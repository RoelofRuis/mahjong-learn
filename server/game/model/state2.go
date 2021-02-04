package model

import (
	"fmt"
	"github.com/roelofruis/mahjong-learn/driver"
	"sort"
)

var (
	stateNewGame       func(table *Table) *driver.State
	stateNextRound     func(table *Table) *driver.State
	stateNextTurn      func(table *Table) *driver.State
	stateMustDiscard   func(table *Table) *driver.State
	stateTileDiscarded func(table *Table) *driver.State
	stateGameEnded     func(table *Table) *driver.State
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

	// TODO: probably move this to game
	var availableActions []driver.Action
	if t.GetActivePlayer().GetReceivedTile() == nil {
		availableActions = t.GetActivePlayer().GetDiscardAfterCombinationActions()
	} else {
		availableActions = t.GetActivePlayer().GetTileReceivedActions()
	}

	sort.Sort(ByActionOrder(availableActions))
	actionMap[t.GetActiveSeat()] = availableActions

	return actionMap
}

func (t *Table) handleMustDiscardActions(actions map[driver.Seat]driver.Action) (*driver.State, error) {
	switch a := actions[t.GetActiveSeat()].(type) {
	case Discard:
		t.ActivePlayerDiscards(a.Tile)
		return stateTileDiscarded, nil

	case DeclareConcealedKong:
		t.ActivePlayerDeclaresConcealedKong(a.Tile)
		t.DealToActivePlayer()
		return stateMustDiscard, nil

	case ExposedPungToKong:
		t.ActivePlayerAddsToExposedPung()
		t.DealToActivePlayer()
		return stateMustDiscard, nil

	case DeclareMahjong:
		// TODO: double check no more logic is needed here
		return stateNextRound, nil

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
		case model.DoNothing:
			value = 1
		case model.DeclareChow:
			value = 2
		case model.DeclarePung:
			value = 3
		case model.DeclareKong:
			value = 4
		case model.DeclareMahjong:
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
	case model.DoNothing:
		t.ActivePlayerTakesDiscarded()
		t.ActivateSeat(bestSeat)
		return stateNextTurn, nil

	case model.DeclareChow:
		t.ActivateSeat(bestSeat)
		t.ActivePlayerTakesChow(a.Tile)
		return stateMustDiscard, nil

	case model.DeclarePung:
		t.ActivateSeat(bestSeat)
		t.ActivePlayerTakesPung()
		return stateMustDiscard, nil

	case model.DeclareKong:
		t.ActivateSeat(bestSeat)
		t.ActivePlayerTakesKong()
		t.DealToActivePlayer()
		return stateMustDiscard, nil

	case model.DeclareMahjong:
		// TODO: double check no more logic is needed here
		return stateNextRound, nil
	}

	panic(fmt.Sprintf("invalid state encountered after resolving tile discarded.\nall actions %+v\nbest action %+v", actions, bestAction))
}

func (t *Table) tryNextRound(_ map[driver.Seat]driver.Action) (*driver.State, error) {
	// TODO: tally scores

	// Game ends if player 3 has been North
	if t.GetPrevalentWind() == North && t.GetPlayerAtSeat(driver.Seat(3)).GetSeatWind() == model.North {
		return stateGameEnded, nil
	}

	if t.GetPlayerAtSeat(model.Seat(3)).GetSeatWind() == t.GetPrevalentWind() {
		t.SetNextPrevalentWind()
	}

	t.ResetWall()
	t.PrepareNextRound()

	return stateNextTurn, nil
}
