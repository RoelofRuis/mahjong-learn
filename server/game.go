package main

var TILESET = []*Tile{
	{name: "bamboo 1"},
}

func InitGame(id uint64) Game {
	players := make(map[int]*Player, 4)

	return Game{
		Id: id,
		HasEnded: false,
		Wall: TILESET,
		Players: players,
	}
}
