package main

import "math/rand"

func InitGame(id uint64) *Game {
	players := make(map[int]Player, 4)

	players[0] = NewPlayer()
	players[1] = NewPlayer()
	players[2] = NewPlayer()
	players[3] = NewPlayer()
	tileSet := NewMahjongSet()

	tileSet.Transfer(14, players[0].Concealed)
	tileSet.Transfer(13, players[1].Concealed)
	tileSet.Transfer(13, players[2].Concealed)
	tileSet.Transfer(13, players[3].Concealed)

	return &Game{
		Id:       id,
		HasEnded: false,
		Wall:     tileSet,
		Players:  players,
	}
}

func NewPlayer() Player {
	return Player{
		Score:     0,
		Concealed: NewEmptyTileCollection(),
		Exposed:   []*TileCollection{},
		Discarded: NewEmptyTileCollection(),
	}
}

func NewEmptyTileCollection() *TileCollection {
	return &TileCollection{Tiles: make(map[Tile]int)}
}

func NewMahjongSet() *TileCollection {
	return &TileCollection{Tiles: map[Tile]int{
		Bamboo1: 4, Bamboo2: 4, Bamboo3: 4, Bamboo4: 4, Bamboo5: 4, Bamboo6: 4, Bamboo7: 4, Bamboo8: 4, Bamboo9: 4,
		Circles1: 4, Circles2: 4, Circles3: 4, Circles4: 4, Circles5: 4, Circles6: 4, Circles7: 4, Circles8: 4, Circles9: 4,
		Characters1: 4, Characters2: 4, Characters3: 4, Characters4: 4, Characters5: 4, Characters6: 4, Characters7: 4, Characters8: 4, Characters9: 4,

		RedDragon: 4, GreenDragon: 4, WhiteDragon: 4,
		EastWind: 4, SouthWind: 4, WestWind: 4, NorthWind: 4,
		FlowerPlumb: 1, FlowerOrchid: 1, FlowerChrysanthemum: 1, FlowerBamboo: 1,
		SeasonSpring: 1, SeasonSummer: 1, SeasonAutumn: 1, SeasonWinter: 1,
	}}
}

// Transfers n tiles from this tile collection to the target tile collection.
func (t *TileCollection) Transfer(n int, target *TileCollection) {
	var tileList = make([]Tile, 0)
	for k, v := range t.Tiles {
		for i:=v; i>0; i-- {
			tileList = append(tileList, k)
		}
	}
	for i:=n; i>0; i-- {
		numTiles := len(tileList)
		pos := rand.Intn(numTiles)
		picked := tileList[pos]

		tileList[pos] = tileList[numTiles-1]
		tileList = tileList[:numTiles-1]

		t.Tiles[picked]--
		target.Tiles[picked]++
	}
}