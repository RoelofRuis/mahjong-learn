package mahjong

import (
	"errors"
	"fmt"
	"github.com/roelofruis/mahjong-learn/state"
)

type Game struct {
	Table        *Table
	StateMachine *state.StateMachine
}

func NewGame(transitioner state.Transitioner) (*Game, error) {
	table := NewTable()
	generator := stateNewGame(table)

	sm := state.NewStateMachine(generator, transitioner)

	err := sm.Transition(nil)
	if err != nil {
		return nil, err
	}

	return &Game{
		Table:        table,
		StateMachine: sm,
	}, nil
}

type StateGenerator func(table *Table) *state.State

var (
	stateNewGame       StateGenerator
	stateNextRound     StateGenerator
	stateNextTurn      StateGenerator
	stateMustDiscard   StateGenerator
	stateTileDiscarded StateGenerator
	stateGameEnded     StateGenerator
)

func init() {
	// initialize states in `init` to prevent loops in references
	stateNewGame = func(table *Table) *state.State {
		return state.NewIntermediateState("New Game", table.initialize)
	}

	stateNextRound = func(table *Table) *state.State {
		return state.NewIntermediateState("Next Round", table.tryNextRound)
	}

	stateNextTurn = func(table *Table) *state.State {
		return state.NewIntermediateState("Next turn", table.tryDealTile)
	}

	stateMustDiscard = func(table *Table) *state.State {
		return state.NewState("Must Discard", table.mustDiscardActions(), table.handleMustDiscardActions)
	}

	stateTileDiscarded = func(table *Table) *state.State {
		return state.NewState("Tile Discarded", table.tileDiscardedActions(), table.handleTileDiscardedActions)
	}

	stateGameEnded = func(table *Table) *state.State {
		return state.NewTerminalState("Game Ended")
	}
}

func (t *Table) initialize() *state.State {
	t.DealConcealed(13, 0)
	t.DealConcealed(13, 1)
	t.DealConcealed(13, 2)
	t.DealConcealed(13, 3)

	return stateNextTurn(t)
}

func (t *Table) tryDealTile() *state.State {
	if t.GetWallSize() <= 14 {
		return stateNextRound(t)
	}

	t.DealToActivePlayer()

	return stateMustDiscard(t)
}

func (t *Table) mustDiscardActions() map[int][]state.Action {
	actionMap := make(map[int][]state.Action, 1)

	if t.GetActivePlayer().GetReceivedTile() == nil {
		actionMap[t.GetActivePlayerIndex()] = t.GetActivePlayer().GetDiscardAfterCombinationActions()
	} else {
		actionMap[t.GetActivePlayerIndex()] = t.GetActivePlayer().GetTileReceivedActions()
	}

	return actionMap
}

func (t *Table) handleMustDiscardActions(actions map[int]state.Action) (*state.State, error) {
	switch a := actions[t.GetActivePlayerIndex()].(type) {
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

func (t *Table) tileDiscardedActions() map[int][]state.Action {
	m := make(map[int][]state.Action, 3)

	activeDiscard := *t.GetActiveDiscard()

	for s, p := range t.GetReactingPlayers() {
		isNextPlayer := (t.GetActivePlayerIndex()+1)%4 == s
		m[s] = p.GetTileDiscardedActions(activeDiscard, isNextPlayer)
	}

	return m
}

func (t *Table) handleTileDiscardedActions(actions map[int]state.Action) (*state.State, error) {
	var bestValue = 0
	var bestPlayer int
	for _, playerIndex := range []int{(t.GetActivePlayerIndex() + 1) % 4, (t.GetActivePlayerIndex() + 2) % 4, (t.GetActivePlayerIndex() + 3) % 4} {
		var value int
		switch actions[playerIndex].(type) {
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
			bestPlayer = playerIndex
		}
	}
	bestAction := actions[bestPlayer]

	switch a := bestAction.(type) {
	case DoNothing:
		t.ActivePlayerTakesDiscarded()
		t.MakePlayerActive(bestPlayer)
		return stateNextTurn(t), nil

	case DeclareChow:
		t.MakePlayerActive(bestPlayer)
		t.ActivePlayerTakesChow(a.Tile)
		return stateMustDiscard(t), nil

	case DeclarePung:
		t.MakePlayerActive(bestPlayer)
		t.ActivePlayerTakesPung()
		return stateMustDiscard(t), nil

	case DeclareKong:
		t.MakePlayerActive(bestPlayer)
		t.ActivePlayerTakesKong()
		t.DealToActivePlayer()
		return stateMustDiscard(t), nil

	case DeclareMahjong:
		// TODO: double check no more logic is needed here
		return stateNextRound(t), nil
	}

	return nil, fmt.Errorf("invalid state encountered after resolving tile discarded.\nall actions %+v\nbest action %+v", actions, bestAction)
}

func (t *Table) tryNextRound() *state.State {
	// TODO: tally scores

	// Game ends if player 3 has been North
	if t.GetPrevalentWind() == North && t.GetPlayerByIndex(3).GetWind() == North {
		return stateGameEnded(t)
	}

	if t.GetPlayerByIndex(3).GetWind() == t.GetPrevalentWind() {
		t.SetNextPrevalentWind()
	}

	t.ResetWall()
	t.PrepareNextRound()

	return stateNextTurn(t)
}
