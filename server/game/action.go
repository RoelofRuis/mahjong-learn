package game

import (
	"github.com/roelofruis/mahjong-learn/driver"
	"sort"
)

// Tile received actions
type Discard struct{ Tile Tile }

func (d Discard) ActionOrder() int { return int(d.Tile) }

type DeclareConcealedKong struct{ Tile Tile }

func (d DeclareConcealedKong) ActionOrder() int { return int(d.Tile) + 100 }

type ExposedPungToKong struct{}

func (d ExposedPungToKong) ActionOrder() int { return 200 }

// Tile discarded actions
type DoNothing struct{}

func (d DoNothing) ActionOrder() int { return 0 }

type DeclareChow struct{ Tile Tile }

func (d DeclareChow) ActionOrder() int { return int(d.Tile) }

type DeclarePung struct{}

func (d DeclarePung) ActionOrder() int { return 100 }

type DeclareKong struct{}

func (d DeclareKong) ActionOrder() int { return 101 }

// Both received and discarded actions
type DeclareMahjong struct{}

func (d DeclareMahjong) ActionOrder() int { return -1 }

// Player actions

func (p *Player) GetDiscardAfterCombinationActions() []driver.Action {
	availableActions := make([]driver.Action, 0)

	for t, c := range p.concealed.tiles {
		availableActions = append(availableActions, Discard{Tile: t})
		if c == 4 {
			availableActions = append(availableActions, DeclareConcealedKong{Tile: t})
		}
	}

	// TODO: add declare mahjong

	sort.Sort(driver.ByActionOrder(availableActions))

	return availableActions
}

func (p *Player) GetTileReceivedActions() []driver.Action {
	availableActions := make([]driver.Action, 0)

	receivedTile := *p.received

	availableActions = append(availableActions, Discard{Tile: receivedTile})

	for t, c := range p.concealed.tiles {
		if t != receivedTile {
			availableActions = append(availableActions, Discard{Tile: t})
		}
		if c == 4 || (c == 3 && t == receivedTile) {
			availableActions = append(availableActions, DeclareConcealedKong{Tile: t})
		}
	}

	if p.exposed.Contains(Pung{Tile: receivedTile}) {
		availableActions = append(availableActions, ExposedPungToKong{})
	}

	// TODO: add declare mahjong

	sort.Sort(driver.ByActionOrder(availableActions))

	return availableActions
}

func (p *Player) GetTileDiscardedActions(discarded Tile, isNextSeat bool) []driver.Action {
	availableActions := make([]driver.Action, 0)

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

	// TODO: add declare mahjong

	sort.Sort(driver.ByActionOrder(availableActions))

	return availableActions
}

func possibleChows(hand *TileCollection, tile Tile) []Tile {
	tileList := make([]Tile, 0)

	if !IsSuit(tile) {
		return nil
	}

	suitType := int(tile) / 10
	suitNumber := int(tile) % 10
	for _, i := range []int{1, 2, 3, 4, 5, 6, 7} {
		diff := suitNumber - i
		if diff >= 0 && diff <= 2 {
			if i != suitNumber && hand.NumOf(Tile((suitType * 10) + i)) == 0 {
				continue
			}
			if i + 1 != suitNumber && hand.NumOf(Tile((suitType * 10) + i + 1)) == 0 {
				continue
			}
			if i + 2 != suitNumber && hand.NumOf(Tile((suitType * 10) + i + 2)) == 0 {
				continue
			}
			tileList = append(tileList, Tile((suitType * 10) + i))
		}
	}

	return tileList
}