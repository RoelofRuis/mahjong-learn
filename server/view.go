package main

import (
	"fmt"
	"github.com/roelofruis/mahjong-learn/game"
)

type TileDescriptor struct {
	Tile game.Tile
	Name string
}

var TileDescriptors = []TileDescriptor{
	{Tile: game.Bamboo1, Name: "Bamboo 1"},
	{Tile: game.Bamboo2, Name: "Bamboo 2"},
	{Tile: game.Bamboo3, Name: "Bamboo 3"},
	{Tile: game.Bamboo4, Name: "Bamboo 4"},
	{Tile: game.Bamboo5, Name: "Bamboo 5"},
	{Tile: game.Bamboo6, Name: "Bamboo 6"},
	{Tile: game.Bamboo7, Name: "Bamboo 7"},
	{Tile: game.Bamboo8, Name: "Bamboo 8"},
	{Tile: game.Bamboo9, Name: "Bamboo 1"},
	{Tile: game.Circles1, Name: "Circles 1"},
	{Tile: game.Circles2, Name: "Circles 2"},
	{Tile: game.Circles3, Name: "Circles 3"},
	{Tile: game.Circles4, Name: "Circles 4"},
	{Tile: game.Circles5, Name: "Circles 5"},
	{Tile: game.Circles6, Name: "Circles 6"},
	{Tile: game.Circles7, Name: "Circles 7"},
	{Tile: game.Circles8, Name: "Circles 8"},
	{Tile: game.Circles9, Name: "Circles 9"},
	{Tile: game.Characters1, Name: "Characters 1"},
	{Tile: game.Characters2, Name: "Characters 2"},
	{Tile: game.Characters3, Name: "Characters 3"},
	{Tile: game.Characters4, Name: "Characters 4"},
	{Tile: game.Characters5, Name: "Characters 5"},
	{Tile: game.Characters6, Name: "Characters 6"},
	{Tile: game.Characters7, Name: "Characters 7"},
	{Tile: game.Characters8, Name: "Characters 8"},
	{Tile: game.Characters9, Name: "Characters 9"},
	{Tile: game.RedDragon, Name: "Red Dragon"},
	{Tile: game.GreenDragon, Name: "Green Dragon"},
	{Tile: game.WhiteDragon, Name: "White Dragon"},
	{Tile: game.EastWind, Name: "East Wind"},
	{Tile: game.SouthWind, Name: "South Wind"},
	{Tile: game.WestWind, Name: "West Wind"},
	{Tile: game.NorthWind, Name: "North Wind"},
	{Tile: game.FlowerPlumb, Name: "Plumb (flower)"},
	{Tile: game.FlowerOrchid, Name: "Orchid (flower)"},
	{Tile: game.FlowerChrysanthemum, Name: "Chrysanthemum (flower)"},
	{Tile: game.FlowerBamboo, Name: "Bamboo (flower)"},
	{Tile: game.SeasonSpring, Name: "Spring (season)"},
	{Tile: game.SeasonSummer, Name: "Summer (season)"},
	{Tile: game.SeasonAutumn, Name: "Autumn (season)"},
	{Tile: game.SeasonWinter, Name: "Winter (season)"},
}

var WindNames = map[game.Wind]string{
	game.East:  "East",
	game.South: "South",
	game.West:  "West",
	game.North: "North",
}

type PlayerView struct {
	Wind      string     `json:"wind"`
	Concealed []string   `json:"hand"`
	Exposed   [][]string `json:"exposed"`
	Discarded []string   `json:"discarded"`
}

type HumanView struct {
	Id            uint64     `json:"id"`
	State         int        `json:"state"` // TODO: map to human readable string
	PrevalentWind string     `json:"prevalent_wind"`
	ActivePlayer  int        `json:"active_player"`
	Wall          []string   `json:"wall"`
	Player1       PlayerView `json:"player_1"`
	Player2       PlayerView `json:"player_2"`
	Player3       PlayerView `json:"player_3"`
	Player4       PlayerView `json:"player_4"`
}

func View(g *game.Game) *HumanView {
	return &HumanView{
		Id:            g.Id,
		State:         int(g.State),
		Wall:          Describe(g.Wall),
		PrevalentWind: WindNames[g.PrevalentWind],
		ActivePlayer:  int(g.ActiveSeat) + 1,
		Player1: PlayerView{
			Wind:      WindNames[g.Players[0].SeatWind],
			Concealed: Describe(g.Players[0].Concealed),
			Exposed:   DescribeAll(g.Players[0].Exposed),
			Discarded: Describe(g.Players[0].Discarded),
		},
		Player2: PlayerView{
			Wind:      WindNames[g.Players[1].SeatWind],
			Concealed: Describe(g.Players[1].Concealed),
			Exposed:   DescribeAll(g.Players[1].Exposed),
			Discarded: Describe(g.Players[1].Discarded),
		},
		Player3: PlayerView{
			Wind:      WindNames[g.Players[2].SeatWind],
			Concealed: Describe(g.Players[2].Concealed),
			Exposed:   DescribeAll(g.Players[2].Exposed),
			Discarded: Describe(g.Players[2].Discarded),
		},
		Player4: PlayerView{
			Wind:      WindNames[g.Players[3].SeatWind],
			Concealed: Describe(g.Players[3].Concealed),
			Exposed:   DescribeAll(g.Players[3].Exposed),
			Discarded: Describe(g.Players[3].Discarded),
		},
	}
}

func DescribeAll(t []*game.TileCollection) [][]string {
	descriptions := make([][]string, len(t))
	for i, col := range t {
		descriptions[i] = Describe(col)
	}
	return descriptions
}

func Describe(t *game.TileCollection) []string {
	descriptions := make([]string, 0)
	for _, d := range TileDescriptors {
		count, has := t.Tiles[d.Tile]
		if !has || count == 0 {
			continue
		}
		text := fmt.Sprintf("%d× %s", count, d.Name)
		descriptions = append(descriptions, text)
	}
	return descriptions
}

func ToVector(t *game.TileCollection) []int {
	tileVector := make([]int, len(TileDescriptors))
	for i, d := range TileDescriptors {
		count, has := t.Tiles[d.Tile]
		if !has {
			count = 0
		}
		tileVector[i] = count
	}
	return tileVector
}
