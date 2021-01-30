package model

type Player struct {
	score int

	seatWind  Wind
	concealed *TileCollection
	exposed   []Combination
	discarded *TileCollection
}

func NewPlayer(seatWind Wind) *Player {
	return &Player{
		score:     0,
		seatWind:  seatWind,
		concealed: NewEmptyTileCollection(),
		exposed:   NewEmptyCombinationList(),
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
	return p.exposed
}

func (p *Player) GetDiscardedTiles() *TileCollection {
	return p.discarded
}

// State modifiers

func (p *Player) ForceExposeTiles() int {
	var transferred = 0
	for t, c := range p.concealed.tiles {
		if IsBonusTile(t) && c > 0 {
			p.concealed.Remove(t)

			p.exposed = append(p.exposed, BonusTile{t})
			transferred++
		}
	}

	return transferred
}

func (p *Player) PrepareNextRound() {
	p.discarded.Empty()
	p.concealed.Empty()
	p.exposed = NewEmptyCombinationList()
	p.seatWind = (p.seatWind + 5) % 4
}
