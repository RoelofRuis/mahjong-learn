package game

func (m *StateMachine) Transition() {
	for {
		if m.State.TransferAction == nil {
			break
		}
		m.State = m.State.TransferAction(m.Game)
	}
	// TODO: player actions
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
	g.DealStartingHands()
	return StateNextTurn
}

func TryDealTile(g *Game) *State {
	if g.Wall.Size() <= 14 {
		// tally scores?
		return StateNextRound
	}

	g.DealTile(g.ActiveSeat)
	return StateTileDealt
}
