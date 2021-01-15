package main

import (
	"fmt"
)

type TileDescriptor struct {
	Tile Tile
	Name string
}

var TileDescriptors = []TileDescriptor {
	{Tile: Bamboo1, Name: "Bamboo 1"},
	{Tile: Bamboo2, Name: "Bamboo 2"},
	{Tile: Bamboo3, Name: "Bamboo 3"},
	{Tile: Bamboo4, Name: "Bamboo 4"},
	{Tile: Bamboo5, Name: "Bamboo 5"},
	{Tile: Bamboo6, Name: "Bamboo 6"},
	{Tile: Bamboo7, Name: "Bamboo 7"},
	{Tile: Bamboo8, Name: "Bamboo 8"},
	{Tile: Bamboo9, Name: "Bamboo 1"},
	{Tile: Circles1, Name: "Circles 1"},
	{Tile: Circles2, Name: "Circles 2"},
	{Tile: Circles3, Name: "Circles 3"},
	{Tile: Circles4, Name: "Circles 4"},
	{Tile: Circles5, Name: "Circles 5"},
	{Tile: Circles6, Name: "Circles 6"},
	{Tile: Circles7, Name: "Circles 7"},
	{Tile: Circles8, Name: "Circles 8"},
	{Tile: Circles9, Name: "Circles 9"},
	{Tile: Characters1, Name: "Characters 1"},
	{Tile: Characters2, Name: "Characters 2"},
	{Tile: Characters3, Name: "Characters 3"},
	{Tile: Characters4, Name: "Characters 4"},
	{Tile: Characters5, Name: "Characters 5"},
	{Tile: Characters6, Name: "Characters 6"},
	{Tile: Characters7, Name: "Characters 7"},
	{Tile: Characters8, Name: "Characters 8"},
	{Tile: Characters9, Name: "Characters 9"},
	{Tile: RedDragon, Name: "Red Dragon"},
	{Tile: GreenDragon, Name: "Green Dragon"},
	{Tile: WhiteDragon, Name: "White Dragon"},
	{Tile: EastWind, Name: "East Wind"},
	{Tile: SouthWind, Name: "South Wind"},
	{Tile: WestWind, Name: "West Wind"},
	{Tile: NorthWind, Name: "North Wind"},
	{Tile: FlowerPlumb, Name: "Plumb (flower)"},
	{Tile: FlowerOrchid, Name: "Orchid (flower)"},
	{Tile: FlowerChrysanthemum, Name: "Chrysanthemum (flower)"},
	{Tile: FlowerBamboo, Name: "Bamboo (flower)"},
	{Tile: SeasonSpring, Name: "Spring (season)"},
	{Tile: SeasonSummer, Name: "Summer (season)"},
	{Tile: SeasonAutumn, Name: "Autumn (season)"},
	{Tile: SeasonWinter, Name: "Winter (season)"},
}

type PlayerView struct {
	Concealed []string `json:"hand"`
	Exposed [][]string `json:"exposed"`
	Discarded []string `json:"discarded"`
}

type HumanView struct {
	Id uint64 `json:"id"`
	Wall []string `json:"wall"`
	Player1 PlayerView `json:"player_1"`
	Player2 PlayerView `json:"player_2"`
	Player3 PlayerView `json:"player_3"`
	Player4 PlayerView `json:"player_4"`
}

func (g *Game) HumanView() *HumanView {
	return &HumanView{
		Id: g.Id,
		Wall: g.Wall.Describe(),
		Player1: PlayerView{
			Concealed: g.Players[0].Concealed.Describe(),
			Exposed: DescribeAll(g.Players[0].Exposed),
			Discarded: g.Players[0].Discarded.Describe(),
		},
		Player2: PlayerView{
			Concealed: g.Players[1].Concealed.Describe(),
			Exposed: DescribeAll(g.Players[1].Exposed),
			Discarded: g.Players[1].Discarded.Describe(),
		},
		Player3: PlayerView{
			Concealed: g.Players[2].Concealed.Describe(),
			Exposed: DescribeAll(g.Players[2].Exposed),
			Discarded: g.Players[2].Discarded.Describe(),
		},
		Player4: PlayerView{
			Concealed: g.Players[3].Concealed.Describe(),
			Exposed: DescribeAll(g.Players[3].Exposed),
			Discarded: g.Players[3].Discarded.Describe(),
		},
	}
}

func DescribeAll(t []*TileCollection) [][]string {
	descriptions := make([][]string, len(t))
	for i, col := range t {
		descriptions[i] = col.Describe()
	}
	return descriptions
}

func (t *TileCollection) Describe() []string {
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

func (t *TileCollection) ToVector() []int {
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