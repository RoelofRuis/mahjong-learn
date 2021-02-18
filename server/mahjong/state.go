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
	table := newTable()
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

type stateGenerator func(table *Table) *state.State

var (
	stateNewGame       stateGenerator
	stateNextRound     stateGenerator
	stateNextTurn      stateGenerator
	stateMustDiscard   stateGenerator
	stateTileDiscarded stateGenerator
	stateGameEnded     stateGenerator
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
	t.dealConcealed(13, 0)
	t.dealConcealed(13, 1)
	t.dealConcealed(13, 2)
	t.dealConcealed(13, 3)

	return stateNextTurn(t)
}

func (t *Table) tryDealTile() *state.State {
	if t.GetWallSize() <= 14 {
		return stateNextRound(t)
	}

	t.dealToActivePlayer()

	return stateMustDiscard(t)
}

func (t *Table) mustDiscardActions() map[int][]state.Action {
	actionMap := make(map[int][]state.Action, 1)

	if t.GetActivePlayer().GetReceivedTile() == nil {
		// player must discard after declaring combination
		actionMap[t.GetActivePlayerIndex()] = t.GetActivePlayer().getDiscardAfterCombinationActions()
	} else {
		// player received tile
		actionMap[t.GetActivePlayerIndex()] = t.GetActivePlayer().getTileReceivedActions()
	}

	return actionMap
}

func (t *Table) handleMustDiscardActions(actions map[int]state.Action) (*state.State, error) {
	switch a := actions[t.GetActivePlayerIndex()].(type) {
	case Discard:
		t.activePlayerDiscards(a.Tile)
		return stateTileDiscarded(t), nil

	case DeclareConcealedKong:
		t.activePlayerDeclaresConcealedKong(a.Tile)
		t.dealToActivePlayer()
		return stateMustDiscard(t), nil

	case ExposedPungToKong:
		t.activePlayerAddsToExposedPung()
		t.dealToActivePlayer()
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
		m[s] = p.getTileDiscardedActions(activeDiscard, isNextPlayer)
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
		t.activePlayerTakesDiscarded()
		t.makePlayerActive(bestPlayer)
		return stateNextTurn(t), nil

	case DeclareChow:
		t.makePlayerActive(bestPlayer)
		t.activePlayerTakesChow(a.Tile)
		return stateMustDiscard(t), nil

	case DeclarePung:
		t.makePlayerActive(bestPlayer)
		t.activePlayerTakesPung()
		return stateMustDiscard(t), nil

	case DeclareKong:
		t.makePlayerActive(bestPlayer)
		t.activePlayerTakesKong()
		t.dealToActivePlayer()
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
		t.setNextPrevalentWind()
	}

	t.resetWall()
	t.prepareNextRound()

	return stateNextTurn(t)
}
