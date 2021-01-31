package game

import (
	"fmt"
	"github.com/roelofruis/mahjong-learn/game/model"
)

var (
	stateNewGame       *State
	stateNextRound     *State
	stateNextTurn      *State
	stateTileReceived  *State
	stateTileDiscarded *State
	stateGameEnded     *State
)

func init() {
	// initialize states here to prevent loop in references
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

	stateTileReceived = &State{
		Name:          "Tile Received",
		PlayerActions: tileReceivedActions,
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

	return stateTileReceived, nil
}

func tileReceivedActions(g *model.Game) map[model.Seat][]model.Action {
	actionMap := make(map[model.Seat][]model.Action, 1)
	actionMap[g.GetActiveSeat()] = g.GetActivePlayer().GetTileReceivedActions()

	return actionMap
}

func handleTileReceivedActions(g *model.Game, actions map[model.Seat]model.Action) (*State, error) {
	switch a := actions[g.GetActiveSeat()].(type) {
	case model.Discard:
		g.ActivePlayerDiscards(a.Tile)
		return stateTileDiscarded, nil

		// TODO: handle other possible cases

	default:
		return nil, fmt.Errorf("illegal action %+v", a)
	}
}

func tileDiscardedActions(g *model.Game) map[model.Seat][]model.Action {
	m := make(map[model.Seat][]model.Action, 3)

	activeDiscard := *g.GetActiveDiscard()

	for s, p := range g.GetReactingPlayers() {
		m[s] = p.GetTileDiscardedActions(activeDiscard)
	}

	return m
}

func handleTileDiscardedActions(g *model.Game, _ map[model.Seat]model.Action) (*State, error) {
	// TODO: handle actions

	g.ActivePlayerTakesDiscarded()
	g.ActivateNextSeat()

	return stateNextTurn, nil
}

func tryNextRound(g *model.Game, _ map[model.Seat]model.Action) (*State, error) {
	if g.GetPrevalentWind() == model.North && g.GetPlayerAtSeat(model.Seat(3)).GetSeatWind() == model.North {
		// Done if player 3 has been North
		return stateGameEnded, nil
	}

	// TODO: tally scores

	if g.GetPlayerAtSeat(model.Seat(3)).GetSeatWind() == g.GetPrevalentWind() {
		g.SetNextPrevalentWind()
	}

	g.ResetWall()
	g.PrepareNextRound()

	return stateNextTurn, nil
}
