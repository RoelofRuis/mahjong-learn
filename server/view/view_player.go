package view

import (
	"github.com/roelofruis/mahjong-learn/mahjong"
)

type OtherPlayer struct {
	Score     int      `json:"score"`
	Wind      string   `json:"wind"`
	Exposed   []string `json:"exposed"`
	Discarded []string `json:"discarded"`
}

type PlayerView struct {
	Actions map[int]string `json:"actions"`

	PrevalentWind    string `json:"prevalent_wind"`
	DiscardingPlayer int    `json:"discarding_player"`
	ActiveDiscard    string `json:"active_discard"`

	Score     int      `json:"score"`
	Wind      string   `json:"wind"`
	Received  string   `json:"received"`
	Concealed []string `json:"concealed"`
	Exposed   []string `json:"exposed"`
	Discarded []string `json:"discarded"`

	OtherPlayers map[int]OtherPlayer `json:"other_players"`
}

func ViewPlayer(game *mahjong.Game, playerIndex int) *PlayerView {
	table := *game.Table

	discardingPlayer := -1
	activeDiscard := "none"
	if table.GetActiveDiscard() != nil {
		discardingPlayer = table.GetActivePlayerIndex()
		activeDiscard = describeTilePointer(table.GetActiveDiscard())
	}

	otherPlayers := make(map[int]OtherPlayer)
	for _, p := range []int{0, 1, 2, 3} {
		if p == playerIndex {
			continue
		}
		otherPlayers[p] = describeOtherPlayer(table, p)
	}

	player := table.GetPlayerByIndex(playerIndex)

	actionMap := make(map[int]string)
	for i, a := range game.StateMachine.AvailableActions()[playerIndex] {
		actionMap[i] = describeAction(a)
	}

	return &PlayerView{
		PrevalentWind:    WindNames[table.GetPrevalentWind()],
		DiscardingPlayer: discardingPlayer,
		ActiveDiscard:    activeDiscard,

		OtherPlayers: otherPlayers,

		Score:     player.GetScore(),
		Wind:      WindNames[player.GetWind()],
		Received:  describeTilePointer(player.GetReceivedTile()),
		Concealed: describeTileCollection(player.GetConcealedTiles()),
		Exposed:   describeCombinations(player.GetExposedCombinations()),
		Discarded: describeTileCollection(player.GetDiscardedTiles()),

		Actions: actionMap,
	}
}

func describeOtherPlayer(table mahjong.Table, player int) OtherPlayer {
	p := table.GetPlayerByIndex(player)

	return OtherPlayer{
		Score:     p.GetScore(),
		Wind:      WindNames[p.GetWind()],
		Exposed:   describeCombinations(p.GetExposedCombinations()),
		Discarded: describeTileCollection(p.GetDiscardedTiles()),
	}
}
