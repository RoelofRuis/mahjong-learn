package model

import "math/rand"

type TileCollection struct {
	tiles map[Tile]int
}

func NewEmptyTileCollection() *TileCollection {
	return &TileCollection{tiles: make(map[Tile]int)}
}

func NewMahjongSet() *TileCollection {
	return &TileCollection{tiles: map[Tile]int{
		Bamboo1: 4, Bamboo2: 4, Bamboo3: 4, Bamboo4: 4, Bamboo5: 4, Bamboo6: 4, Bamboo7: 4, Bamboo8: 4, Bamboo9: 4,
		Circles1: 4, Circles2: 4, Circles3: 4, Circles4: 4, Circles5: 4, Circles6: 4, Circles7: 4, Circles8: 4, Circles9: 4,
		Characters1: 4, Characters2: 4, Characters3: 4, Characters4: 4, Characters5: 4, Characters6: 4, Characters7: 4, Characters8: 4, Characters9: 4,

		RedDragon: 4, GreenDragon: 4, WhiteDragon: 4,
		EastWind: 4, SouthWind: 4, WestWind: 4, NorthWind: 4,

		FlowerPlumb: 1, FlowerOrchid: 1, FlowerChrysanthemum: 1, FlowerBamboo: 1,
		SeasonSpring: 1, SeasonSummer: 1, SeasonAutumn: 1, SeasonWinter: 1,
	}}
}

// Getters

func (t *TileCollection) NumOf(tile Tile) int {
	count, has := t.tiles[tile]
	if !has {
		return 0
	}
	return count
}

func (t *TileCollection) Size() int {
	var count = 0
	for _, c := range t.tiles {
		count += c
	}

	return count
}

// State Modifiers

func (t *TileCollection) Empty() {
	t.tiles = make(map[Tile]int)
}

func (t *TileCollection) Remove(tile Tile) {
	n, has := t.tiles[tile]
	if !has || n < 1 {
		return
	}
	t.tiles[tile]--
}

func (t *TileCollection) Add(tile Tile) {
	t.tiles[tile]++
}

func (t *TileCollection) Transfer(tile Tile, target *TileCollection) {
	n, has := t.tiles[tile]
	if !has || n == 0 {
		return
	}

	t.tiles[tile]--
	target.Add(tile)
}

// Transfers n randomly picked tiles from this tile collection to the target tile collection.
func (t *TileCollection) TransferRandom(n int, target *TileCollection) {
	var tileList = make([]Tile, 0)
	for k, v := range t.tiles {
		for i := v; i > 0; i-- {
			tileList = append(tileList, k)
		}
	}
	for i := n; i > 0; i-- {
		numTiles := len(tileList)
		pos := rand.Intn(numTiles)
		picked := tileList[pos]

		tileList[pos] = tileList[numTiles-1]
		tileList = tileList[:numTiles-1]

		t.tiles[picked]--
		target.tiles[picked]++
	}
}
