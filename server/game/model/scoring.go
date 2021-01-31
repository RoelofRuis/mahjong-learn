package model

import "sort"

// TODO: maybe move these functions to `action.go`

func (p *Player) GetTileReceivedActions() []Action {
	availableActions := make([]Action, 0)

	receivedTile := *p.received

	availableActions = append(availableActions, Discard{Tile: receivedTile})

	for t, c := range p.concealed.tiles {
		if c > 0 && t != receivedTile {
			// a player can discard any tile he has at least one of
			availableActions = append(availableActions, Discard{Tile: t})
		}
		if c == 4 || (c == 3 && t == receivedTile) {
			// a player can declare a concealed kong of four equal tiles
			availableActions = append(availableActions, DeclareConcealedKong{Tile: t})
		}
	}

	if p.exposed.Contains(Pung{Tile: receivedTile}) {
		// add to an exposed pung to make a kong
		availableActions = append(availableActions, ExposedPungToKong{Tile: receivedTile})
	}

	// TODO: add declare mahjong

	sort.Sort(ByIndex(availableActions))

	return availableActions
}

func (p *Player) GetTileDiscardedActions(discarded Tile) []Action {
	availableActions := make([]Action, 0)

	availableActions = append(availableActions, DoNothing{})

	if p.concealed.NumOf(discarded) == 2 {
		availableActions = append(availableActions, DeclarePung{})
	}
	if p.concealed.NumOf(discarded) == 3 {
		availableActions = append(availableActions, DeclareKong{})
	}

	// TODO: add chow
	// TODO: add mahjong

	sort.Sort(ByIndex(availableActions))

	return availableActions
}
