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

type GameView struct {
	Id uint64 `json:"id"`
	Wall []string `json:"wall"`
	WallVector []int `json:"wall_vector"`
}

func (g *Game) View() *GameView {
	return &GameView{
		Id: g.Id,
		Wall: g.Wall.Describe(),
		WallVector: g.Wall.ToVector(),
	}
}

func (t *TileCollection) Describe() []string {
	var descriptions []string
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