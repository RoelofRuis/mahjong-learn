package mahjong

import (
	"math/rand"
	"sort"
)

type TileCollection struct {
	tiles map[Tile]uint8
}

func newEmptyTileCollection() *TileCollection {
	return &TileCollection{tiles: make(map[Tile]uint8)}
}

func newMahjongSet() *TileCollection {
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

func (t *TileCollection) OrderedCounts() []TileCount {
	i := 0
	ordered := make([]TileCount, len(t.tiles))
	for tile, count := range t.tiles {
		ordered[i] = TileCount{tile, count}
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].Count > ordered[j].Count
	})
	return ordered
}

type TileCount struct {
	Tile Tile
	Count uint8
}

// State Modifiers

func (t *TileCollection) empty() {
	t.tiles = make(map[Tile]uint8)
}

func (t *TileCollection) removeAll(tile Tile) {
	delete(t.tiles, tile)
}

func (t *TileCollection) remove(tile Tile) {
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

func (t *TileCollection) add(tile Tile) {
	t.tiles[tile]++
}

func (t *TileCollection) removeRandom() Tile {
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

func newCombinationCollection() *CombinationCollection {
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

func (c *CombinationCollection) empty() {
	c.combinations = []Combination{}
}

func (c *CombinationCollection) replace(old Combination, new Combination) {
	for i, comb := range c.combinations {
		if comb == old {
			c.combinations[i] = new
		}
	}
}

func (c *CombinationCollection) add(combination Combination) {
	c.combinations = append(c.combinations, combination)
}

type Combination interface {
	// has to be unique among all defined combinations (to guarantee a stable sorting)
	CombinationOrder() int
}

type Chow struct{ FirstTile Tile }

func (c Chow) CombinationOrder() int { return int(c.FirstTile) }

type Pung struct{ Tile Tile }

func (c Pung) CombinationOrder() int { return int(c.Tile) + 100 }

type Kong struct {
	Tile      Tile
	Concealed bool
}

func (c Kong) CombinationOrder() int { return int(c.Tile) + 200 }

type BonusTile struct {
	Tile Tile
}

func (c BonusTile) CombinationOrder() int { return int(c.Tile) + 300 }

type ByCombinationOrder []Combination

func (a ByCombinationOrder) Len() int      { return len(a) }
func (a ByCombinationOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCombinationOrder) Less(i, j int) bool {
	return a[i].CombinationOrder() < a[j].CombinationOrder()
}
