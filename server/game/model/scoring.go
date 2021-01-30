package model

func (p *Player) CanPung(t Tile) bool {
	return p.concealed.NumOf(t) == 2
}

func (p *Player) CanKong(t Tile) bool {
	return p.concealed.NumOf(t) == 3
}

func (p *Player) PossibleDiscards() []Tile {
	var l []Tile
	for t, c := range p.concealed.tiles {
		if c > 0 {
			l = append(l, t)
		}
	}
	return l
}

func (p *Player) PossibleHiddenKongs() []Tile {
	var l []Tile
	for t, c := range p.concealed.tiles {
		if c == 4 {
			l = append(l, t)
		}
	}
	return l
}

func (p *Player) CanDeclareMahjong() bool {
	// TODO: implement
	return false
}
