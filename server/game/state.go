package game

func (m *StateMachine) Transition() {
	m.lock.Lock()
	defer m.lock.Unlock()

	for {
		if m.state.TransferAction == nil {
			break
		}
		m.state = m.state.TransferAction(m.game)
	}
	// TODO: handle player actions
}

// If the StateMachine is viewed, internals should be exposed in a consistent manner, so one function returns everything.
func (m *StateMachine) View() (Game, State, map[Seat][]PlayerAction) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var playerActions map[Seat][]PlayerAction
	if m.state.PlayerActions == nil {
		playerActions = make(map[Seat][]PlayerAction, 0)
	}

	playerActions = m.state.PlayerActions(m.game)

	return *m.game, *m.state, playerActions
}

var StateNewGame = &State{
	Name:           "New Game",
	TransferAction: Initialize,
	PlayerActions:  nil,
}

var StateNextRound = &State{
	Name:           "Next Round",
	TransferAction: nil,
	PlayerActions:  nil,
}

var StateNextTurn = &State{
	Name:           "Next Turn",
	TransferAction: TryDealTile,
	PlayerActions:  nil,
}

var StateTileReceived = &State{
	Name:           "Tile Received",
	TransferAction: nil,
	PlayerActions:  ReactToTile,
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
	return StateTileReceived
}

func ReactToTile(g *Game) map[Seat][]PlayerAction {
	m := make(map[Seat][]PlayerAction, 1)
	var a = []PlayerAction{
		{
			Name: "Discard a tile",
		},
	}

	m[g.ActiveSeat] = a

	return m
}