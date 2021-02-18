package mahjong

import (
	"github.com/roelofruis/mahjong-learn/state"
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

func (p *Player) getDiscardAfterCombinationActions() []state.Action {
	availableActions := make([]state.Action, 0)

	for t, c := range p.concealed.tiles {
		availableActions = append(availableActions, Discard{Tile: t})
		if c == 4 {
			availableActions = append(availableActions, DeclareConcealedKong{Tile: t})
		}
	}

	if mahjongPossible(p) {
		availableActions = append(availableActions, DeclareMahjong{})
	}

	return availableActions
}

func (p *Player) getTileReceivedActions() []state.Action {
	availableActions := make([]state.Action, 0)

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

	if mahjongPossible(p) {
		availableActions = append(availableActions, DeclareMahjong{})
	}

	return availableActions
}

func (p *Player) getTileDiscardedActions(discarded Tile, isNextPlayer bool) []state.Action {
	availableActions := make([]state.Action, 0)

	availableActions = append(availableActions, DoNothing{})

	if p.concealed.NumOf(discarded) == 2 {
		availableActions = append(availableActions, DeclarePung{})
	}

	if p.concealed.NumOf(discarded) == 3 {
		availableActions = append(availableActions, DeclareKong{})
	}

	if isNextPlayer {
		for _, c := range possibleChows(p.concealed, discarded) {
			availableActions = append(availableActions, DeclareChow{Tile: c})
		}
	}

	if mahjongPossible(p) {
		availableActions = append(availableActions, DeclareMahjong{})
	}

	return availableActions
}

func possibleChows(hand *TileCollection, tile Tile) []Tile {
	if !tile.IsSuit() {
		return nil
	}

	tileList := make([]Tile, 0)

	suitType := (int(tile) / 10) * 10
	suitNumber := int(tile) % 10
	for _, i := range []int{1, 2, 3, 4, 5, 6, 7} {
		tileI := suitType + i
		diff := suitNumber - i
		if diff < 0 {
			continue
		}
		if diff > 2 {
			break
		}
		if i != suitNumber && hand.NumOf(Tile(tileI)) == 0 {
			continue
		}
		if i+1 != suitNumber && hand.NumOf(Tile(tileI+1)) == 0 {
			continue
		}
		if i+2 != suitNumber && hand.NumOf(Tile(tileI+2)) == 0 {
			continue
		}
		tileList = append(tileList, Tile(tileI))
	}

	return tileList
}

func mahjongPossible(player *Player) bool {
	numPungs := 0
	numKongs := 0
	numChows := 0
	numDragonCombinations := 0
	numBambooCombination := 0
	numCircleCombinations := 0
	numCharacterCombinations := 0
	numOwnBonus := 0
	hasOwnWind := false

	for _, comb := range player.exposed.combinations {
		switch c := comb.(type) {
		case Chow:
			numChows++
			if c.FirstTile.IsBamboo() {
				numBambooCombination++
			}
			if c.FirstTile.IsCircle() {
				numCircleCombinations++
			}
			if c.FirstTile.IsCharacter() {
				numCharacterCombinations++
			}

		case Pung:
			numPungs++
			if c.Tile.IsDragon() {
				numDragonCombinations++
			}
			if c.Tile.IsSameWindAs(player.wind) {
				hasOwnWind = true
			}
			if c.Tile.IsBamboo() {
				numBambooCombination++
			}
			if c.Tile.IsCircle() {
				numCircleCombinations++
			}
			if c.Tile.IsCharacter() {
				numCharacterCombinations++
			}

		case Kong:
			numKongs++
			if c.Tile.IsDragon() {
				numDragonCombinations++
			}
			if c.Tile.IsSameWindAs(player.wind) {
				hasOwnWind = true
			}
			if c.Tile.IsBamboo() {
				numBambooCombination++
			}
			if c.Tile.IsCircle() {
				numCircleCombinations++
			}
			if c.Tile.IsCharacter() {
				numCharacterCombinations++
			}

		case BonusTile:
			if int(c.Tile) % 10 == int(player.wind) {
				numOwnBonus++
			}
		}
	}



	// 4 groups + pair



	return false
}