package game

type StateMachine struct {
	State State

	Game *Game
}

type Action struct {
	Update func(*Game) error
}

type State struct {
	Name string
	StateActions []Action
	IsTransition bool
	TransitionTo *State
	PlayerActions func(*Game) map[Seat][]Action
}
