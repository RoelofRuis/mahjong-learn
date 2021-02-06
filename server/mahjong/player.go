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

func NewPlayer(wind Wind) *Player {
	return &Player{
		score:     0,
		received:  nil,
		wind:      wind,
		concealed: NewEmptyTileCollection(),
		exposed:   NewCombinationCollection(),
		discarded: NewEmptyTileCollection(),
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

func (p *Player) GetDiscardedTiles() *TileCollection {
	return p.discarded
}

func (p *Player) GetReceivedTile() *Tile {
	return p.received
}

func (p *Player) GetScore() int {
	return p.score
}
