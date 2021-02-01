package model

import "math/rand"

type TileCollection struct {
	tiles map[Tile]uint8
}

func NewEmptyTileCollection() *TileCollection {
	return &TileCollection{tiles: make(map[Tile]uint8)}
}

func NewMahjongSet() *TileCollection {
	return &TileCollection{tiles: map[Tile]uint8{
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
		count = 0
	}
	return int(count)
}

func (t *TileCollection) Size() int {
	var count = 0
	for _, c := range t.tiles {
		count += int(c)
	}

	return count
}

// State Modifiers

func (t *TileCollection) Empty() {
	t.tiles = make(map[Tile]uint8)
}

func (t *TileCollection) Remove(tile Tile) {
	n, has := t.tiles[tile]
	if !has {
		return
	}
	if n == 1 {
		delete(t.tiles, tile)
	} else {
		t.tiles[tile]--
	}
}

func (t *TileCollection) Add(tile Tile) {
	t.tiles[tile]++
}

func (t *TileCollection) RemoveRandom() Tile {
	var tileList = make([]Tile, 0)
	for k, v := range t.tiles {
		for i := v; i > 0; i-- {
			tileList = append(tileList, k)
		}
	}
	pos := rand.Intn(len(tileList))
	picked := tileList[pos]

	t.tiles[picked]--

	return picked
}

type CombinationCollection struct {
	combinations []Combination
}

func NewCombinationCollection() *CombinationCollection {
	return &CombinationCollection{combinations: []Combination{}}
}

// Getters

func (c CombinationCollection) Contains(check Combination) bool {
	for _, comb := range c.combinations {
		if comb == check {
			return true
		}
	}
	return false
}

// State modifiers

func (c *CombinationCollection) Empty() {
	c.combinations = []Combination{}
}

func (c *CombinationCollection) Add(combination Combination) {
	c.combinations = append(c.combinations, combination)
}

type Combination interface {
	CombinationIndex() int // TODO: meh, not sure if this is really required...
}

type Chow struct {
	FirstTile Tile
}

func (c Chow) CombinationIndex() int {
	return int(c.FirstTile)
}

type Pung struct {
	Tile Tile
}

func (c Pung) CombinationIndex() int {
	return int(c.Tile) + 100
}

type Kong struct {
	Tile   Tile
	Hidden bool
}

func (c Kong) CombinationIndex() int {
	return int(c.Tile) + 200
}

type BonusTile struct {
	Tile Tile
}

func (c BonusTile) CombinationIndex() int {
	return int(c.Tile) + 300
}
