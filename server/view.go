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
	Actions []string `json:"actions"`
	Wind      string     `json:"wind"`
	Concealed []string   `json:"hand"`
	Exposed   [][]string `json:"exposed"`
	Discarded []string   `json:"discarded"`
}

type HumanView struct {
	Id              uint64        `json:"id"`
	StateName       string        `json:"state_name"`
	PrevalentWind   string        `json:"prevalent_wind"`
	ActivePlayer    int           `json:"active_player"`
	Player1         PlayerView    `json:"player_1"`
	Player2         PlayerView    `json:"player_2"`
	Player3         PlayerView    `json:"player_3"`
	Player4         PlayerView    `json:"player_4"`
	Wall            []string      `json:"wall"`
}

func View(stateMachine *game.StateMachine) *HumanView {
	g, s, a := stateMachine.View()
	return &HumanView{
		Id:            g.Id,
		StateName:     s.Name,
		PrevalentWind: WindNames[g.PrevalentWind],
		ActivePlayer:  int(g.ActiveSeat) + 1,
		Player1: DescribePlayer(g, a, 0),
		Player2: DescribePlayer(g, a, 1),
		Player3: DescribePlayer(g, a, 2),
		Player4: DescribePlayer(g, a, 3),
		Wall:          Describe(g.Wall),
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

func DescribePlayer(g game.Game, a map[game.Seat][]game.PlayerAction, player int) PlayerView {
	seat := game.Seat(player)
	actions, has := a[seat]
	if ! has {
		actions = make([]game.PlayerAction, 0)
	}

	actionList := make([]string, 0)
	for _, a := range actions {
		actionList = append(actionList, a.Name)
	}

	return PlayerView{
		Actions: actionList,
		Wind:      WindNames[g.Players[seat].SeatWind],
		Concealed: Describe(g.Players[seat].Concealed),
		Exposed:   DescribeAll(g.Players[seat].Exposed),
		Discarded: Describe(g.Players[seat].Discarded),
	}
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
