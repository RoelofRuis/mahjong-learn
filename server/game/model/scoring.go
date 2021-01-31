package model

import "sort"

func (p *Player) GetTileReceivedActions() []Action {
	availableActions := make([]Action, 0)

	for t, c := range p.concealed.tiles {
		if c > 0 {
			availableActions = append(availableActions, Discard{Tile: t})
		}
		if c == 4 {
			availableActions = append(availableActions, DeclareConcealedKong{Tile: t})
		}
		// TODO: add to exposed pung to make kong
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
