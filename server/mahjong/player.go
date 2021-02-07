package mahjong

import "sort"

type Player struct {
	score int

	received  *Tile
	wind      Wind
	concealed *TileCollection
	exposed   *CombinationCollection
	discarded *TileCollection
}

func newPlayer(wind Wind) *Player {
	return &Player{
		score:     0,
		received:  nil,
		wind:      wind,
		concealed: newEmptyTileCollection(),
		exposed:   newCombinationCollection(),
		discarded: newEmptyTileCollection(),
	}
}

// Getters

func (p *Player) GetConcealedTiles() *TileCollection {
	return p.concealed
}

func (p *Player) GetWind() Wind {
	return p.wind
}

func (p *Player) GetExposedCombinations() []Combination {
	sort.Sort(ByCombinationOrder(p.exposed.combinations))

	return p.exposed.combinations
}

func (p *Player) GetExposedCombinationCollection() *CombinationCollection {
	return p.exposed
}

func (p *Player) GetDiscardedTiles() *TileCollection {
	return p.discarded
}

func (p *Player) GetReceivedTile() *Tile {
	return p.received
}

func (p *Player) GetScore() int {
	return p.score
}
