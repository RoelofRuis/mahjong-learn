package game

import "math/rand"

type Tile int

const (
	Bamboo1             Tile = 1
	Bamboo2             Tile = 2
	Bamboo3             Tile = 3
	Bamboo4             Tile = 4
	Bamboo5             Tile = 5
	Bamboo6             Tile = 6
	Bamboo7             Tile = 7
	Bamboo8             Tile = 8
	Bamboo9             Tile = 9
	Circles1            Tile = 11
	Circles2            Tile = 12
	Circles3            Tile = 13
	Circles4            Tile = 14
	Circles5            Tile = 15
	Circles6            Tile = 16
	Circles7            Tile = 17
	Circles8            Tile = 18
	Circles9            Tile = 19
	Characters1         Tile = 21
	Characters2         Tile = 22
	Characters3         Tile = 23
	Characters4         Tile = 24
	Characters5         Tile = 25
	Characters6         Tile = 26
	Characters7         Tile = 27
	Characters8         Tile = 28
	Characters9         Tile = 29
	RedDragon           Tile = 30
	GreenDragon         Tile = 31
	WhiteDragon         Tile = 32
	EastWind            Tile = 40
	SouthWind           Tile = 41
	WestWind            Tile = 42
	NorthWind           Tile = 43
	FlowerPlumb         Tile = 50
	FlowerOrchid        Tile = 51
	FlowerChrysanthemum Tile = 52
	FlowerBamboo        Tile = 53
	SeasonSpring        Tile = 60
	SeasonSummer        Tile = 61
	SeasonAutumn        Tile = 62
	SeasonWinter        Tile = 63
)

type Game struct {
	// FIXME: add lock to Game so we can modify data freely in a request and block simultaneous requests
	Id uint64

	HasEnded bool

	Wall    *TileCollection
	Players map[int]Player
}

type Player struct {
	Score int

	Concealed *TileCollection
	Exposed   []*TileCollection
	Discarded *TileCollection
}

type TileCollection struct {
	Tiles map[Tile]int
}

// Transfers n randomly picked tiles from this tile collection to the target tile collection.
func (t *TileCollection) Transfer(n int, target *TileCollection) {
	var tileList = make([]Tile, 0)
	for k, v := range t.Tiles {
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

		t.Tiles[picked]--
		target.Tiles[picked]++
	}
}
