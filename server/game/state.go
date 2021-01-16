package game

type StateMachine struct {
	State *State

	Game *Game
}

type State struct {
	// Name just to display human readable information.
	Name string

	// Transfer to next state via action, or nil if player input is required.
	TransferAction Action
	// Show required player actions. This requires TransferAction to be nil.
	RequiredActions func(*Game) map[Seat][]Action
}

type Action func(*Game) *State

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

func Initialize(g *Game) *State {
	g.DealStartingHands()
	return StateNextTurn
}

var StateNextTurn = &State{
	Name:            "Next Turn",
	TransferAction:  nil,
	RequiredActions: nil,
}
