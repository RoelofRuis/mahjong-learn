package view

import (
	"github.com/roelofruis/mahjong-learn/mahjong"
	"github.com/roelofruis/mahjong-learn/state"
)

type OtherPlayer struct {
	Score int `json:"score"`
	Wind string `json:"wind"`
	Exposed []string `json:"exposed"`
	Discarded []string `json:"discarded"`
	ActiveDiscard string `json:"active_discarded"`
}

type PlayerView struct {
	PrevalentWind string             `json:"prevalent_wind"`
	OtherPlayers map[int]OtherPlayer `json:"other_players"`

	Score int `json:"score"`
	Wind string `json:"wind"`
	Received string `json:"received"`
	Concealed []string `json:"concealed"`
	Exposed []string `json:"exposed"`
	Discarded []string `json:"discarded"`
}

func ViewPlayer(game *mahjong.Game, seat state.Seat) *PlayerView {
	// TODO: implement

	return &PlayerView{}
}
