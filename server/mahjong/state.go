package mahjong

import (
	"fmt"
	"github.com/roelofruis/mahjong-learn/driver"
)

// TODO: refactor to remove panic calls!

type MahjongGame struct {
	Id uint64

	Table *Table

	Driver *driver.GameDriver
}

func NewMahjongGame(id uint64) *MahjongGame {
	table := NewTable()
	state := stateNewGame(table)

	gameDriver := driver.NewGameDriver(state, 10)

	err := gameDriver.Transition(nil)
	if err != nil {
		panic(err) // TODO: clean return instead of panic!
	}

	return &MahjongGame{
		Id:     id,
		Table:  table,
		Driver: gameDriver,
	}
}

type MahjongState func(table *Table) *driver.State

var (
	stateNewGame       MahjongState
	stateNextRound     MahjongState
	stateNextTurn      MahjongState
	stateMustDiscard   MahjongState
	stateTileDiscarded MahjongState
	stateGameEnded     MahjongState
)

func init() {
	// initialize states in `init` to prevent loops in references
	stateNewGame = func(table *Table) *driver.State {
		return &driver.State{
			Name:          "New Game",
			Actions: nil,
			Transition:    table.initialize,
		}
	}

	stateNextRound = func(table *Table) *driver.State {
		return &driver.State{
			Name:          "Next Round",
			Actions: nil,
			Transition:    table.tryNextRound,
		}
	}

	stateNextTurn = func(table *Table) *driver.State {
		return &driver.State{
			Name:    "Next Turn",
			Actions: nil,
			Transition: table.tryDealTile,
		}
	}

	stateMustDiscard = func(table *Table) *driver.State {
		return &driver.State{
			Name:          "Must Discard",
			Actions: table.mustDiscardActions,
			Transition:    table.handleMustDiscardActions,
		}
	}

	stateTileDiscarded = func(table *Table) *driver.State {
		return &driver.State{
			Name:          "Tile Discarded",
			Actions: table.tileDiscardedActions,
			Transition:    table.handleTileDiscardedActions,
		}
	}

	stateGameEnded = func(table *Table) *driver.State {
		return &driver.State{
			Name:          "Game Ended",
			Actions: nil,
			Transition:    nil,
		}
	}
}

func (t *Table) initialize(_ map[driver.Seat]driver.Action) (*driver.State, error) {
	t.DealConcealed(13, 0)
	t.DealConcealed(13, 1)
	t.DealConcealed(13, 2)
	t.DealConcealed(13, 3)

	return stateNextTurn(t), nil
}

func (t *Table) tryDealTile(_ map[driver.Seat]driver.Action) (*driver.State, error) {
	if t.GetWallSize() <= 14 {
		return stateNextRound(t), nil
	}

	t.DealToActivePlayer()

	return stateMustDiscard(t), nil
}

func (t *Table) mustDiscardActions() map[driver.Seat][]driver.Action {
	actionMap := make(map[driver.Seat][]driver.Action, 1)

	// TODO: probably move this to mahjong..?
	if t.GetActivePlayer().GetReceivedTile() == nil {
		actionMap[t.GetActiveSeat()] = t.GetActivePlayer().GetDiscardAfterCombinationActions()
	} else {
		actionMap[t.GetActiveSeat()] = t.GetActivePlayer().GetTileReceivedActions()
	}

	return actionMap
}

func (t *Table) handleMustDiscardActions(actions map[driver.Seat]driver.Action) (*driver.State, error) {
	switch a := actions[t.GetActiveSeat()].(type) {
	case Discard:
		t.ActivePlayerDiscards(a.Tile)
		return stateTileDiscarded(t), nil

	case DeclareConcealedKong:
		t.ActivePlayerDeclaresConcealedKong(a.Tile)
		t.DealToActivePlayer()
		return stateMustDiscard(t), nil

	case ExposedPungToKong:
		t.ActivePlayerAddsToExposedPung()
		t.DealToActivePlayer()
		return stateMustDiscard(t), nil

	case DeclareMahjong:
		// TODO: double check no more logic is needed here
		return stateNextRound(t), nil

	default:
		return nil, fmt.Errorf("illegal action %+v", a)
	}
}

func (t *Table) tileDiscardedActions() map[driver.Seat][]driver.Action {
	m := make(map[driver.Seat][]driver.Action, 3)

	activeDiscard := *t.GetActiveDiscard()

	for s, p := range t.GetReactingPlayers() {
		isNextSeat := (t.GetActiveSeat() + 1)%4 == s
		m[s] = p.GetTileDiscardedActions(activeDiscard, isNextSeat)
	}

	return m
}

func (t *Table) handleTileDiscardedActions(actions map[driver.Seat]driver.Action) (*driver.State, error) {
	var bestValue = 0
	var bestSeat driver.Seat
	for _, seatIndex := range []driver.Seat{ (t.GetActiveSeat() + 1) % 4, (t.GetActiveSeat() + 2) % 4, (t.GetActiveSeat() + 3) % 4 } {
		var value int
		switch actions[seatIndex].(type) {
		case DoNothing:
			value = 1
		case DeclareChow:
			value = 2
		case DeclarePung:
			value = 3
		case DeclareKong:
			value = 4
		case DeclareMahjong:
			value = 5
		default:
			panic("invalid action given in response to 'handleTileDiscarded'")
		}
		if value > bestValue {
			bestValue = value
			bestSeat = seatIndex
		}
	}
	bestAction := actions[bestSeat]

	switch a := bestAction.(type) {
	case DoNothing:
		t.ActivePlayerTakesDiscarded()
		t.ActivateSeat(bestSeat)
		return stateNextTurn(t), nil

	case DeclareChow:
		t.ActivateSeat(bestSeat)
		t.ActivePlayerTakesChow(a.Tile)
		return stateMustDiscard(t), nil

	case DeclarePung:
		t.ActivateSeat(bestSeat)
		t.ActivePlayerTakesPung()
		return stateMustDiscard(t), nil

	case DeclareKong:
		t.ActivateSeat(bestSeat)
		t.ActivePlayerTakesKong()
		t.DealToActivePlayer()
		return stateMustDiscard(t), nil

	case DeclareMahjong:
		// TODO: double check no more logic is needed here
		return stateNextRound(t), nil
	}

	panic(fmt.Sprintf("invalid state encountered after resolving tile discarded.\nall actions %+v\nbest action %+v", actions, bestAction))
}

func (t *Table) tryNextRound(_ map[driver.Seat]driver.Action) (*driver.State, error) {
	// TODO: tally scores

	// Game ends if player 3 has been North
	if t.GetPrevalentWind() == North && t.GetPlayerAtSeat(driver.Seat(3)).GetSeatWind() == North {
		return stateGameEnded(t), nil
	}

	if t.GetPlayerAtSeat(driver.Seat(3)).GetSeatWind() == t.GetPrevalentWind() {
		t.SetNextPrevalentWind()
	}

	t.ResetWall()
	t.PrepareNextRound()

	return stateNextTurn(t), nil
}
