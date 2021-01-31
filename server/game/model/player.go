package model

type Player struct {
	score int

	received  *Tile
	seatWind  Wind
	concealed *TileCollection
	exposed   *CombinationCollection
	discarded *TileCollection
}

func NewPlayer(seatWind Wind) *Player {
	return &Player{
		score:     0,
		received:  nil,
		seatWind:  seatWind,
		concealed: NewEmptyTileCollection(),
		exposed:   NewCombinationCollection(),
		discarded: NewEmptyTileCollection(),
	}
}

// Getters

func (p *Player) GetConcealedTiles() *TileCollection {
	return p.concealed
}

func (p *Player) GetSeatWind() Wind {
	return p.seatWind
}

func (p *Player) GetExposedCombinations() []Combination {
	return p.exposed.combinations
}

func (p *Player) GetDiscardedTiles() *TileCollection {
	return p.discarded
}

func (p *Player) GetReceivedTile() *Tile {
	return p.received
}
