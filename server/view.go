package main

import (
	"fmt"
	"github.com/roelofruis/mahjong-learn/game"
)

var TileOrder = []game.Tile{
	game.Bamboo1,
	game.Bamboo2,
	game.Bamboo3,
	game.Bamboo4,
	game.Bamboo5,
	game.Bamboo6,
	game.Bamboo7,
	game.Bamboo8,
	game.Bamboo9,
	game.Circles1,
	game.Circles2,
	game.Circles3,
	game.Circles4,
	game.Circles5,
	game.Circles6,
	game.Circles7,
	game.Circles8,
	game.Circles9,
	game.Characters1,
	game.Characters2,
	game.Characters3,
	game.Characters4,
	game.Characters5,
	game.Characters6,
	game.Characters7,
	game.Characters8,
	game.Characters9,
	game.RedDragon,
	game.GreenDragon,
	game.WhiteDragon,
	game.EastWind,
	game.SouthWind,
	game.WestWind,
	game.NorthWind,
	game.FlowerPlumb,
	game.FlowerOrchid,
	game.FlowerChrysanthemum,
	game.FlowerBamboo,
	game.SeasonSpring,
	game.SeasonSummer,
	game.SeasonAutumn,
	game.SeasonWinter,
}

var WindNames = map[game.Wind]string{
	game.East:  "East",
	game.South: "South",
	game.West:  "West",
	game.North: "North",
}

type PlayerView struct {
	Actions   map[int]string `json:"actions"`
	Wind      string         `json:"wind"`
	Concealed []string       `json:"hand"`
	Exposed   []string       `json:"exposed"`
	Discarded []string       `json:"discarded"`
}

type HumanView struct {
	Id            uint64     `json:"id"`
	HasEnded      bool       `json:"has_ended"`
	StateName     string     `json:"state_name"`
	PrevalentWind string     `json:"prevalent_wind"`
	ActivePlayer  int        `json:"active_player"`
	ActiveDiscard string     `json:"active_discard"`
	Player1       PlayerView `json:"player_1"`
	Player2       PlayerView `json:"player_2"`
	Player3       PlayerView `json:"player_3"`
	Player4       PlayerView `json:"player_4"`
	Wall          []string   `json:"wall"`
}

func View(stateMachine *game.StateMachine) *HumanView {
	g, s, a := stateMachine.View()
	return &HumanView{
		Id:            stateMachine.Id(),
		HasEnded:      s.IsTerminal(),
		StateName:     s.Name,
		PrevalentWind: WindNames[g.PrevalentWind],
		ActivePlayer:  int(g.ActiveSeat) + 1,
		ActiveDiscard: DescribeActiveDiscard(g.ActiveDiscard),
		Player1:       DescribePlayer(g, a, 0),
		Player2:       DescribePlayer(g, a, 1),
		Player3:       DescribePlayer(g, a, 2),
		Player4:       DescribePlayer(g, a, 3),
		Wall:          Describe(g.Wall),
	}
}

func DescribeActiveDiscard(t *game.Tile) string {
	if t == nil {
		return "none"
	}
	return game.TileNames[*t]
}

func DescribeCombinations(combinations []game.Combination) []string {
	descriptions := make([]string, len(combinations))
	for i, combi := range combinations {
		switch c := combi.(type) {
		case game.BonusTile:
			descriptions[i] = fmt.Sprintf("Bonus tile %s", game.TileNames[c.Tile])

		case game.Chow:
			descriptions[i] = fmt.Sprintf("Chow %s", game.TileNames[c.FirstTile])

		case game.Pung:
			descriptions[i] = fmt.Sprintf("Pung %s", game.TileNames[c.Tile])

		case game.Kong:
			descriptions[i] = fmt.Sprintf("Kong %s", game.TileNames[c.Tile])

		default:
			// This should not happen..!
			descriptions[i] = "unknown combination"
		}
	}
	return descriptions
}

func Describe(t *game.TileCollection) []string {
	descriptions := make([]string, 0)
	for _, tile := range TileOrder {
		count, has := t.Tiles[tile]
		if !has || count == 0 {
			continue
		}
		text := fmt.Sprintf("%d× %s", count, game.TileNames[tile])
		descriptions = append(descriptions, text)
	}
	return descriptions
}

func DescribePlayer(g game.Game, a map[game.Seat][]game.PlayerAction, player int) PlayerView {
	seat := game.Seat(player)
	actions, has := a[seat]
	if !has {
		actions = make([]game.PlayerAction, 0)
	}

	actionMap := make(map[int]string)
	for i, a := range actions {
		actionMap[i] = a.Name
	}

	return PlayerView{
		Actions:   actionMap,
		Wind:      WindNames[g.Players[seat].SeatWind],
		Concealed: Describe(g.Players[seat].Concealed),
		Exposed:   DescribeCombinations(g.Players[seat].Exposed),
		Discarded: Describe(g.Players[seat].Discarded),
	}
}

func ToVector(t *game.TileCollection) []int {
	tileVector := make([]int, len(TileOrder))
	for i, tile := range TileOrder {
		count, has := t.Tiles[tile]
		if !has {
			count = 0
		}
		tileVector[i] = count
	}
	return tileVector
}
