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
		Transition:    handleTileReceivedActions,
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

func initialize(g *model.Game, _ map[model.Seat]model.Action) (*State, error) {
	g.DealConcealed(13, 0)
	g.DealConcealed(13, 1)
	g.DealConcealed(13, 2)
	g.DealConcealed(13, 3)

	return stateNextTurn, nil
}

func tryDealTile(g *model.Game, _ map[model.Seat]model.Action) (*State, error) {
	if g.GetWallSize() <= 14 {
		return stateNextRound, nil
	}

	g.DealToActivePlayer()

	return stateMustDiscard, nil
}

func mustDiscardActions(g *model.Game) map[model.Seat][]model.Action {
	actionMap := make(map[model.Seat][]model.Action, 1)

	// TODO: probably move this to game
	var availableActions []model.Action
	if g.GetActivePlayer().GetReceivedTile() == nil {
		availableActions = g.GetActivePlayer().GetDiscardAfterCombinationActions()
	} else {
		availableActions = g.GetActivePlayer().GetTileReceivedActions()
	}

	sort.Sort(model.ByActionOrder(availableActions))
	actionMap[g.GetActiveSeat()] = availableActions

	return actionMap
}

func handleTileReceivedActions(g *model.Game, actions map[model.Seat]model.Action) (*State, error) {
	switch a := actions[g.GetActiveSeat()].(type) {
	case model.Discard:
		g.ActivePlayerDiscards(a.Tile)
		return stateTileDiscarded, nil

	case model.DeclareConcealedKong:
		g.ActivePlayerDeclaresConcealedKong(a.Tile)
		g.DealToActivePlayer()
		return stateMustDiscard, nil

	case model.ExposedPungToKong:
		g.ActivePlayerAddsToExposedPung()
		g.DealToActivePlayer()
		return stateMustDiscard, nil

	case model.DeclareMahjong:
		// TODO: double check no more logic is needed here
		return stateNextRound, nil

	default:
		return nil, fmt.Errorf("illegal action %+v", a)
	}
}

func tileDiscardedActions(g *model.Game) map[model.Seat][]model.Action {
	m := make(map[model.Seat][]model.Action, 3)

	activeDiscard := *g.GetActiveDiscard()

	for s, p := range g.GetReactingPlayers() {
		isNextSeat := (g.GetActiveSeat() + 1)%4 == s
		m[s] = p.GetTileDiscardedActions(activeDiscard, isNextSeat)
	}

	return m
}

func handleTileDiscardedActions(g *model.Game, actions map[model.Seat]model.Action) (*State, error) {
	var bestValue = 0
	var bestSeat model.Seat
	for _, seatIndex := range []model.Seat{ (g.GetActiveSeat() + 1) % 4, (g.GetActiveSeat() + 2) % 4, (g.GetActiveSeat() + 3) % 4 } {
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
		g.ActivePlayerTakesDiscarded()
		g.ActivateSeat(bestSeat)
		return stateNextTurn, nil

	case model.DeclareChow:
		g.ActivateSeat(bestSeat)
		g.ActivePlayerTakesChow(a.Tile)
		return stateMustDiscard, nil

	case model.DeclarePung:
		g.ActivateSeat(bestSeat)
		g.ActivePlayerTakesPung()
		return stateMustDiscard, nil

	case model.DeclareKong:
		g.ActivateSeat(bestSeat)
		g.ActivePlayerTakesKong()
		g.DealToActivePlayer()
		return stateMustDiscard, nil

	case model.DeclareMahjong:
		// TODO: double check no more logic is needed here
		return stateNextRound, nil
	}

	panic(fmt.Sprintf("invalid state encountered after resolving tile discarded.\nall actions %+v\nbest action %+v", actions, bestAction))
}

func tryNextRound(g *model.Game, _ map[model.Seat]model.Action) (*State, error) {
	// TODO: tally scores

	// Game ends if player 3 has been North
	if g.GetPrevalentWind() == model.North && g.GetPlayerAtSeat(model.Seat(3)).GetSeatWind() == model.North {
		return stateGameEnded, nil
	}

	if g.GetPlayerAtSeat(model.Seat(3)).GetSeatWind() == g.GetPrevalentWind() {
		g.SetNextPrevalentWind()
	}

	g.ResetWall()
	g.PrepareNextRound()

	return stateNextTurn, nil
}
