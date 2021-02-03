package game

import (
	"fmt"
	"github.com/roelofruis/mahjong-learn/game/model"
	"sort"
)

var (
	stateNewGame       *State
	stateNextRound     *State
	stateNextTurn      *State
	stateMustDiscard   *State
	stateTileDiscarded *State
	stateGameEnded     *State
)

func init() {
	// initialize states in `init` to prevent loops in references
	stateNewGame = &State{
		Name:          "New Game",
		PlayerActions: nil,
		Transition:    initialize,
	}

	stateNextRound = &State{
		Name:          "Next Round",
		PlayerActions: nil,
		Transition:    tryNextRound,
	}

	stateNextTurn = &State{
		Name:          "Next Turn",
		PlayerActions: nil,
		Transition:    tryDealTile,
	}

	stateMustDiscard = &State{
		Name:          "Must Discard",
		PlayerActions: mustDiscardActions,
		Transition:    handleMustDiscardActions,
	}

	stateTileDiscarded = &State{
		Name:          "Tile Discarded",
		PlayerActions: tileDiscardedActions,
		Transition:    handleTileDiscardedActions,
	}

	stateGameEnded = &State{
		Name:          "Game Ended",
		PlayerActions: nil,
		Transition:    nil,
	}
}

func initialize(t *model.Table, _ map[model.Seat]model.Action) (*State, error) {
	t.DealConcealed(13, 0)
	t.DealConcealed(13, 1)
	t.DealConcealed(13, 2)
	t.DealConcealed(13, 3)

	return stateNextTurn, nil
}

func tryDealTile(t *model.Table, _ map[model.Seat]model.Action) (*State, error) {
	if t.GetWallSize() <= 14 {
		return stateNextRound, nil
	}

	t.DealToActivePlayer()

	return stateMustDiscard, nil
}

func mustDiscardActions(t *model.Table) map[model.Seat][]model.Action {
	actionMap := make(map[model.Seat][]model.Action, 1)

	// TODO: probably move this to game
	var availableActions []model.Action
	if t.GetActivePlayer().GetReceivedTile() == nil {
		availableActions = t.GetActivePlayer().GetDiscardAfterCombinationActions()
	} else {
		availableActions = t.GetActivePlayer().GetTileReceivedActions()
	}

	sort.Sort(model.ByActionOrder(availableActions))
	actionMap[t.GetActiveSeat()] = availableActions

	return actionMap
}

func handleMustDiscardActions(t *model.Table, actions map[model.Seat]model.Action) (*State, error) {
	switch a := actions[t.GetActiveSeat()].(type) {
	case model.Discard:
		t.ActivePlayerDiscards(a.Tile)
		return stateTileDiscarded, nil

	case model.DeclareConcealedKong:
		t.ActivePlayerDeclaresConcealedKong(a.Tile)
		t.DealToActivePlayer()
		return stateMustDiscard, nil

	case model.ExposedPungToKong:
		t.ActivePlayerAddsToExposedPung()
		t.DealToActivePlayer()
		return stateMustDiscard, nil

	case model.DeclareMahjong:
		// TODO: double check no more logic is needed here
		return stateNextRound, nil

	default:
		return nil, fmt.Errorf("illegal action %+v", a)
	}
}

func tileDiscardedActions(t *model.Table) map[model.Seat][]model.Action {
	m := make(map[model.Seat][]model.Action, 3)

	activeDiscard := *t.GetActiveDiscard()

	for s, p := range t.GetReactingPlayers() {
		isNextSeat := (t.GetActiveSeat() + 1)%4 == s
		m[s] = p.GetTileDiscardedActions(activeDiscard, isNextSeat)
	}

	return m
}

func handleTileDiscardedActions(t *model.Table, actions map[model.Seat]model.Action) (*State, error) {
	var bestValue = 0
	var bestSeat model.Seat
	for _, seatIndex := range []model.Seat{ (t.GetActiveSeat() + 1) % 4, (t.GetActiveSeat() + 2) % 4, (t.GetActiveSeat() + 3) % 4 } {
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

func tryNextRound(t *model.Table, _ map[model.Seat]model.Action) (*State, error) {
	// TODO: tally scores

	// Game ends if player 3 has been North
	if t.GetPrevalentWind() == model.North && t.GetPlayerAtSeat(model.Seat(3)).GetSeatWind() == model.North {
		return stateGameEnded, nil
	}

	if t.GetPlayerAtSeat(model.Seat(3)).GetSeatWind() == t.GetPrevalentWind() {
		t.SetNextPrevalentWind()
	}

	t.ResetWall()
	t.PrepareNextRound()

	return stateNextTurn, nil
}
