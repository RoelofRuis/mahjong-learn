package view

import (
	"github.com/roelofruis/mahjong-learn/mahjong"
	"github.com/roelofruis/mahjong-learn/state"
)

type OtherPlayer struct {
	Score     int      `json:"score"`
	Wind      string   `json:"wind"`
	Exposed   []string `json:"exposed"`
	Discarded []string `json:"discarded"`
}

type PlayerView struct {
	PrevalentWind    string `json:"prevalent_wind"`
	DiscardingPlayer int    `json:"discarding_player"`
	ActiveDiscard    string `json:"active_discard"`

	OtherPlayers map[int]OtherPlayer `json:"other_players"`

	Score     int      `json:"score"`
	Wind      string   `json:"wind"`
	Received  string   `json:"received"`
	Concealed []string `json:"concealed"`
	Exposed   []string `json:"exposed"`
	Discarded []string `json:"discarded"`
}

func ViewPlayer(game *mahjong.Game, seat state.Seat) *PlayerView {
	table := *game.Table

	discardingPlayer := -1
	activeDiscard := "none"
	if table.GetActiveDiscard() == nil {
		discardingPlayer = int(table.GetActiveSeat())
		activeDiscard = describeTilePointer(table.GetActiveDiscard())
	}

	otherPlayers := make(map[int]OtherPlayer)
	for _, s := range []int{0, 1, 2, 3} {
		if state.Seat(s) == seat {
			continue
		}
		otherPlayers[s] = describeOtherPlayer(table, state.Seat(s))
	}

	p := table.GetPlayerAtSeat(seat)

	return &PlayerView{
		PrevalentWind:    WindNames[table.GetPrevalentWind()],
		DiscardingPlayer: discardingPlayer,
		ActiveDiscard:    activeDiscard,

		OtherPlayers: otherPlayers,

		Score:     p.GetScore(),
		Wind:      WindNames[p.GetSeatWind()],
		Received:  describeTilePointer(p.GetReceivedTile()),
		Concealed: describeTileCollection(p.GetConcealedTiles()),
		Exposed:   describeCombinations(p.GetExposedCombinations()),
		Discarded: describeTileCollection(p.GetDiscardedTiles()),
	}
}

func describeOtherPlayer(table mahjong.Table, seat state.Seat) OtherPlayer {
	p := table.GetPlayerAtSeat(seat)

	return OtherPlayer{
		Score:     p.GetScore(),
		Wind:      WindNames[p.GetSeatWind()],
		Exposed:   describeCombinations(p.GetExposedCombinations()),
		Discarded: describeTileCollection(p.GetDiscardedTiles()),
	}
}
