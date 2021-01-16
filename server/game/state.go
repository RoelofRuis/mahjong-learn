package game

func (m *StateMachine) Transition() {
	m.Lock()
	for {
		if m.State.TransferAction == nil {
			break
		}
		m.State = m.State.TransferAction(m.Game)
	}
	// TODO: player actions
	m.Unlock()
}

var StateNewGame = &State{
	Name:            "New Game",
	TransferAction:  Initialize,
	RequiredActions: nil,
}

var StateNextRound = &State{
	Name:            "Next Round",
	TransferAction:  nil,
	RequiredActions: nil,
}

var StateNextTurn = &State{
	Name:            "Next Turn",
	TransferAction:  TryDealTile,
	RequiredActions: nil,
}

var StateTileDealt = &State{
	Name:            "Tile Dealt",
	TransferAction:  nil,
	RequiredActions: nil,
}

func Initialize(g *Game) *State {
	g.DealTiles(13, 0)
	g.DealTiles(13, 1)
	g.DealTiles(13, 2)
	g.DealTiles(13, 3)
	return StateNextTurn
}

func TryDealTile(g *Game) *State {
	if g.Wall.Size() <= 14 {
		// tally scores?
		return StateNextRound
	}

	g.DealTiles(1, g.ActiveSeat)
	return StateTileDealt
}
