package driver

type State interface {
	// Name just to display human readable information.
	Name() string
	// Required actions per seat. May be nil if this state requires no actions.
	Actions() map[Seat][]Action
	// Transition to next state. Selected actions are passed if applicable.
	Transition(map[Seat]Action) (State, error)
	// Whether this is a terminal state and the state machine cannot progress further.
	IsTerminal() bool
}

type Action interface {
	ActionOrder() uint64
}

type Seat int