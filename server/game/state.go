package game

type StateMachine struct {
	State *State

	Game *Game
}

type Action struct {
	Update func(*Game) *State
}

type State struct {
	// Name just to display human readable information.
	Name string

	// Transfer to next state via action, or nil if player input is required.
	TransferAction *Action
	// Show required player actions. This requires TransferAction to be nil.
	RequiredActions func(*Game) map[Seat][]Action
}

func (m *StateMachine) Transition() {
	for {
		if m.State.TransferAction == nil {
			break
		}
		m.State = m.State.TransferAction.Update(m.Game)
	}
	// TODO: player actions
}

var StateNewGame = &State{
	Name:            "New Game",
	TransferAction:  nil,
	RequiredActions: nil,
}

var StateNextTurn = &State{
	Name:            "Next Turn",
	TransferAction:  nil,
	RequiredActions: nil,
}

var InitializeGame = &Action{
	Update: func(*Game) *State {
		return nil
	},
}
