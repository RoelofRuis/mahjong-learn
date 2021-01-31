package model

import "sort"

type Action interface {
	// ActionIndex, has to be unique among all defined actions (to guarantee a stable sorting)
	ActionIndex() int
}

type ByIndex []Action

func (a ByIndex) Len() int           { return len(a) }
func (a ByIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIndex) Less(i, j int) bool { return a[i].ActionIndex() < a[j].ActionIndex() }

type Discard struct{ Tile Tile }

func (d Discard) ActionIndex() int { return int(d.Tile) }

type DeclareConcealedKong struct{ Tile Tile }

func (d DeclareConcealedKong) ActionIndex() int { return int(d.Tile) + 100 }

type ExposedPungToKong struct{ Tile Tile }

func (d ExposedPungToKong) ActionIndex() int { return int(d.Tile) + 200 }

type DeclareMahjong struct{}

func (d DeclareMahjong) ActionIndex() int { return -1 }

type DoNothing struct{}

func (d DoNothing) ActionIndex() int { return 0 }

type DeclareChow struct{ Tile Tile }

func (d DeclareChow) ActionIndex() int { return int(d.Tile) }

type DeclarePung struct{}

func (d DeclarePung) ActionIndex() int { return 100 }

type DeclareKong struct{}

func (d DeclareKong) ActionIndex() int { return 101 }

func (p *Player) GetTileReceivedActions() []Action {
	availableActions := make([]Action, 0)

	receivedTile := *p.received

	availableActions = append(availableActions, Discard{Tile: receivedTile})

	for t, c := range p.concealed.tiles {
		if c > 0 && t != receivedTile {
			availableActions = append(availableActions, Discard{Tile: t})
		}
		if c == 4 || (c == 3 && t == receivedTile) {
			availableActions = append(availableActions, DeclareConcealedKong{Tile: t})
		}
	}

	if p.exposed.Contains(Pung{Tile: receivedTile}) {
		availableActions = append(availableActions, ExposedPungToKong{Tile: receivedTile})
	}

	// TODO: add declare mahjong

	sort.Sort(ByIndex(availableActions))

	return availableActions
}

func (p *Player) GetTileDiscardedActions(discarded Tile, isNextSeat bool) []Action {
	availableActions := make([]Action, 0)

	availableActions = append(availableActions, DoNothing{})

	if p.concealed.NumOf(discarded) == 2 {
		availableActions = append(availableActions, DeclarePung{})
	}

	if p.concealed.NumOf(discarded) == 3 {
		availableActions = append(availableActions, DeclareKong{})
	}

	if isNextSeat {
		for _, c := range possibleChows(p.concealed, discarded) {
			availableActions = append(availableActions, DeclareChow{Tile: c})
		}
	}

	// TODO: add mahjong

	sort.Sort(ByIndex(availableActions))

	return availableActions
}

func possibleChows(hand *TileCollection, tile Tile) []Tile {
	if !IsSuit(tile) {
		return nil
	}

	// TODO: calculate chow possibilities
	return nil
}