package mahjong

import (
	"errors"
	"fmt"
	"github.com/roelofruis/mahjong-learn/state_machine"
)

type Game struct {
	Id uint64

	Table *Table

	StateMachine state_machine.StateMachine
}

func NewGame(id uint64) (*Game, error) {
	table := NewTable()
	state := stateNewGame(table)

	sm := state_machine.NewStateMachine(state, 10)

	err := sm.Transition(nil)
	if err != nil {
		return nil, err
	}

	return &Game{
		Id:           id,
		Table:        table,
		StateMachine: sm,
	}, nil
}

type State func(table *Table) *state_machine.State

var (
	stateNewGame       State
	stateNextRound     State
	stateNextTurn      State
	stateMustDiscard   State
	stateTileDiscarded State
	stateGameEnded     State
)

func init() {
	// initialize states in `init` to prevent loops in references
	stateNewGame = func(table *Table) *state_machine.State {
		return state_machine.NewIntermediateState("New Game", table.initialize)
	}

	stateNextRound = func(table *Table) *state_machine.State {
		return state_machine.NewIntermediateState("Next Round", table.tryNextRound)
	}

	stateNextTurn = func(table *Table) *state_machine.State {
		return state_machine.NewIntermediateState("Next turn", table.tryDealTile)
	}

	stateMustDiscard = func(table *Table) *state_machine.State {
		return state_machine.NewState("Must Discard", table.mustDiscardActions(), table.handleMustDiscardActions)
	}

	stateTileDiscarded = func(table *Table) *state_machine.State {
		return state_machine.NewState("Tile Discarded", table.tileDiscardedActions(), table.handleTileDiscardedActions)
	}

	stateGameEnded = func(table *Table) *state_machine.State {
		return state_machine.NewTerminalState("Game Ended")
	}
}

func (t *Table) initialize(_ map[state_machine.Seat]state_machine.Action) (*state_machine.State, error) {
	t.DealConcealed(13, 0)
	t.DealConcealed(13, 1)
	t.DealConcealed(13, 2)
	t.DealConcealed(13, 3)

	return stateNextTurn(t), nil
}

func (t *Table) tryDealTile(_ map[state_machine.Seat]state_machine.Action) (*state_machine.State, error) {
	if t.GetWallSize() <= 14 {
		return stateNextRound(t), nil
	}

	t.DealToActivePlayer()

	return stateMustDiscard(t), nil
}

func (t *Table) mustDiscardActions() map[state_machine.Seat][]state_machine.Action {
	actionMap := make(map[state_machine.Seat][]state_machine.Action, 1)

	if t.GetActivePlayer().GetReceivedTile() == nil {
		actionMap[t.GetActiveSeat()] = t.GetActivePlayer().GetDiscardAfterCombinationActions()
	} else {
		actionMap[t.GetActiveSeat()] = t.GetActivePlayer().GetTileReceivedActions()
	}

	return actionMap
}

func (t *Table) handleMustDiscardActions(actions map[state_machine.Seat]state_machine.Action) (*state_machine.State, error) {
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

func (t *Table) tileDiscardedActions() map[state_machine.Seat][]state_machine.Action {
	m := make(map[state_machine.Seat][]state_machine.Action, 3)

	activeDiscard := *t.GetActiveDiscard()

	for s, p := range t.GetReactingPlayers() {
		isNextSeat := (t.GetActiveSeat() + 1)%4 == s
		m[s] = p.GetTileDiscardedActions(activeDiscard, isNextSeat)
	}

	return m
}

func (t *Table) handleTileDiscardedActions(actions map[state_machine.Seat]state_machine.Action) (*state_machine.State, error) {
	var bestValue = 0
	var bestSeat state_machine.Seat
	for _, seatIndex := range []state_machine.Seat{(t.GetActiveSeat() + 1) % 4, (t.GetActiveSeat() + 2) % 4, (t.GetActiveSeat() + 3) % 4 } {
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
			return nil, errors.New("invalid action given in response to `handleTileDiscarded`")
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

	return nil, fmt.Errorf("invalid state encountered after resolving tile discarded.\nall actions %+v\nbest action %+v", actions, bestAction)
}

func (t *Table) tryNextRound(_ map[state_machine.Seat]state_machine.Action) (*state_machine.State, error) {
	// TODO: tally scores

	// Game ends if player 3 has been North
	if t.GetPrevalentWind() == North && t.GetPlayerAtSeat(state_machine.Seat(3)).GetSeatWind() == North {
		return stateGameEnded(t), nil
	}

	if t.GetPlayerAtSeat(state_machine.Seat(3)).GetSeatWind() == t.GetPrevalentWind() {
		t.SetNextPrevalentWind()
	}

	t.ResetWall()
	t.PrepareNextRound()

	return stateNextTurn(t), nil
}
