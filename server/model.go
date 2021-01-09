package main

type Tile struct {
	name string
}

type Combination struct {
	tiles []*Tile
}

type Wall struct {
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

	Wall Wall
	Players map[int]*Player
}
