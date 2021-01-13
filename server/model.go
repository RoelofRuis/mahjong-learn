package main

type Tile struct {
	Name string
}

type Combination struct {
	Tiles []Tile
}

type Player struct {
	Score int

	ConcealedTiles []Tile
	ExposedCombinations []Combination
	Discards []Tile
}

type Game struct {
	Id uint64

	HasEnded bool

	Wall []Tile
	Players map[int]Player
}
