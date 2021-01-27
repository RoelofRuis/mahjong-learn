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
	g.DealTiles(13, 0)
	g.DealTiles(13, 1)
	g.DealTiles(13, 2)
	g.DealTiles(13, 3)

	return stateNextTurn, nil
}

func tryDealTile(g *model.Game, _ map[model.Seat]model.Action) (*State, error) {
	if g.GetWallSize() <= 14 {
		return stateNextRound, nil
	}

	g.DealTiles(1, g.GetActiveSeat())

	return stateTileReceived, nil
}

func tileReceivedActions(g *model.Game) map[model.Seat][]PlayerAction {
	actionMap := make(map[model.Seat][]PlayerAction, 1)

	availableActions := make([]PlayerAction, 0)
	for t, c := range g.GetActivePlayer().GetConcealedTiles().AsMap() {
		if c < 1 {
			continue
		}

		availableActions = append(availableActions, PlayerAction{
			Name:   fmt.Sprintf("Discard a %s", model.TileNames[t]),
			Action: model.Discard{Tile: t},
		})
	}

	// TODO: check whether player can declare kong or mahjong and add to available actions

	sort.Sort(ByIndex(availableActions))

	actionMap[g.GetActiveSeat()] = availableActions

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

func tileDiscardedActions(g *model.Game) map[model.Seat][]PlayerAction {
	m := make(map[model.Seat][]PlayerAction, 4)

	for _, s := range model.SEATS {
		if s == g.GetActiveSeat() {
			continue
		}

		a := make([]PlayerAction, 0)

		a = append(a, PlayerAction{
			Name:   fmt.Sprintf("Do nothing"),
			Action: model.DoNothing{},
		})

		// TODO: check whether player can declare pung, kong, chow or mahjong and add to available actions

		m[s] = a
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
	if g.GetPrevalentWind() == model.North && g.GetPlayerAtSeat(model.Seat(1)).GetSeatWind() == model.North {
		// Done if player 3 has been North
		return stateGameEnded, nil
	}

	// TODO: tally scores

	// FIXME: not sure if it is player 1 or player 3 that gets to be prevalent wind last...
	if g.GetPlayerAtSeat(model.Seat(1)).GetSeatWind() == g.GetPrevalentWind() {
		g.SetNextPrevalentWind()
	}

	g.ResetWall()

	for _, s := range model.SEATS {
		g.GetPlayerAtSeat(s).PrepareNextRound()
		g.DealTiles(13, s)
	}


	return stateNextTurn, nil
}
