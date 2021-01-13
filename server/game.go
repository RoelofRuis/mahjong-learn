package main

var TILESET = []Tile{
	{Name: "bamboo 1"},
	{Name: "bamboo 2"},
	{Name: "bamboo 3"},
}

func InitGame(id uint64) Game {
	players := make(map[int]Player, 4)

	tiles := make([]Tile, len(TILESET))
	copy(tiles, TILESET)

	return Game{
		Id: id,
		HasEnded: false,
		Wall: tiles,
		Players: players,
	}
}
