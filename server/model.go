package main

type Tile struct {
	name string
}

type Combination struct {
	tiles []*Tile
}

type Player struct {
	score int

	concealedTiles []*Tile
	exposedCombinations []*Combination
	discards []*Tile
}

type Game struct {
	Id uint64

	HasEnded bool

	Wall []*Tile
	Players map[int]*Player
}
