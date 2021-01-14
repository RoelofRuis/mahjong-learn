package main

import (
	"fmt"
)

type TileDescriptor struct {
	Name string
	Pos int
}

var TileDescriptors = map[Tile]TileDescriptor {
	Bamboo1: {Name: "Bamboo 1", Pos: 0},
	Bamboo2: {Name: "Bamboo 2", Pos: 1},
	Bamboo3: {Name: "Bamboo 3", Pos: 2},
	Bamboo4: {Name: "Bamboo 4", Pos: 3},
	Bamboo5: {Name: "Bamboo 5", Pos: 4},
	Bamboo6: {Name: "Bamboo 6", Pos: 5},
	Bamboo7: {Name: "Bamboo 7", Pos: 6},
	Bamboo8: {Name: "Bamboo 8", Pos: 7},
	Bamboo9: {Name: "Bamboo 1", Pos: 8},
	Circles1: {Name: "Circles 1", Pos: 9},
	Circles2: {Name: "Circles 2", Pos: 10},
	Circles3: {Name: "Circles 3", Pos: 11},
	Circles4: {Name: "Circles 4", Pos: 12},
	Circles5: {Name: "Circles 5", Pos: 13},
	Circles6: {Name: "Circles 6", Pos: 14},
	Circles7: {Name: "Circles 7", Pos: 15},
	Circles8: {Name: "Circles 8", Pos: 16},
	Circles9: {Name: "Circles 9", Pos: 17},
	Characters1: {Name: "Characters 1", Pos: 18},
	Characters2: {Name: "Characters 2", Pos: 19},
	Characters3: {Name: "Characters 3", Pos: 20},
	Characters4: {Name: "Characters 4", Pos: 21},
	Characters5: {Name: "Characters 5", Pos: 22},
	Characters6: {Name: "Characters 6", Pos: 23},
	Characters7: {Name: "Characters 7", Pos: 24},
	Characters8: {Name: "Characters 8", Pos: 25},
	Characters9: {Name: "Characters 9", Pos: 26},
	RedDragon: {Name: "Red Dragon", Pos: 27},
	GreenDragon: {Name: "Green Dragon", Pos: 28},
	WhiteDragon: {Name: "White Dragon", Pos: 29},
	EastWind: {Name: "East Wind", Pos: 30},
	SouthWind: {Name: "South Wind", Pos: 31},
	WestWind: {Name: "West Wind", Pos: 32},
	NorthWind: {Name: "North Wind", Pos: 33},
	FlowerPlumb: {Name: "Plumb (flower)", Pos: 34},
	FlowerOrchid: {Name: "Orchid (flower)", Pos: 35},
	FlowerChrysanthemum: {Name: "Chrysanthemum (flower)", Pos: 36},
	FlowerBamboo: {Name: "Bamboo (flower)", Pos: 37},
	SeasonSpring: {Name: "Spring (season)", Pos: 38},
	SeasonSummer: {Name: "Summer (season)", Pos: 39},
	SeasonAutumn: {Name: "Autumn (season)", Pos: 40},
	SeasonWinter: {Name: "Winter (season)", Pos: 41},
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
	for tile, count := range t.Tiles {
		if count == 0 {
			continue
		}
		d, _ := TileDescriptors[tile]
		text := fmt.Sprintf("%d× %s", count, d.Name)
		descriptions = append(descriptions, text)
	}
	return descriptions
}

func (t *TileCollection) ToVector() []int {
	tileVector := make([]int, 42)
	for tile, count := range t.Tiles {
		d, _ := TileDescriptors[tile]
		tileVector[d.Pos] = count
	}
	return tileVector
}