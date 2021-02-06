package view

import (
	"github.com/roelofruis/mahjong-learn/mahjong"
	"github.com/roelofruis/mahjong-learn/state"
)

type GamePlayerView struct {
	Actions   map[int]string `json:"actions"`
	Score     int            `json:"score"`
	Wind      string         `json:"wind"`
	Received  string         `json:"received"`
	Concealed []string       `json:"concealed"`
	Exposed   []string       `json:"exposed"`
	Discarded []string       `json:"discarded"`
}

type GameView struct {
	HasEnded      bool                   `json:"has_ended"`
	StateName     string                 `json:"state_name"`
	PrevalentWind string                 `json:"prevalent_wind"`
	ActivePlayers []int                  `json:"active_players"`
	ActiveDiscard string                 `json:"active_discard"`
	Players       map[int]GamePlayerView `json:"players"`
	Wall          []string               `json:"wall"`
}

func ViewGame(game *mahjong.Game) *GameView {
	table := *game.Table

	var activePlayers []int
	playerViews := make(map[int]GamePlayerView, 4)
	for _, player := range []int{0, 1, 2, 3} {
		actions, has := game.StateMachine.AvailableActions()[player]
		if !has {
			actions = make([]state.Action, 0)
		} else {
			activePlayers = append(activePlayers, player+1)
		}
		playerViews[player+1] = describeGamePlayer(table, actions, player)
	}

	return &GameView{
		HasEnded:      game.StateMachine.HasTerminated(),
		StateName:     game.StateMachine.StateName(),
		PrevalentWind: WindNames[table.GetPrevalentWind()],
		ActivePlayers: activePlayers,
		ActiveDiscard: describeTilePointer(table.GetActiveDiscard()),
		Players:       playerViews,
		Wall:          describeTileCollection(table.GetWall()),
	}
}

func describeGamePlayer(g mahjong.Table, actions []state.Action, player int) GamePlayerView {
	actionMap := make(map[int]string)
	for i, a := range actions {
		actionMap[i] = describeAction(a)
	}

	p := g.GetPlayerByIndex(player)

	return GamePlayerView{
		Actions:   actionMap,
		Score:     p.GetScore(),
		Wind:      WindNames[p.GetWind()],
		Received:  describeTilePointer(p.GetReceivedTile()),
		Concealed: describeTileCollection(p.GetConcealedTiles()),
		Exposed:   describeCombinations(p.GetExposedCombinations()),
		Discarded: describeTileCollection(p.GetDiscardedTiles()),
	}
}
