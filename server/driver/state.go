package driver

type State struct {
	// Name just to display human readable information.
	Name string

	// Required actions per seat. May be nil if this state requires no actions.
	Actions func() map[Seat][]Action

	// Transition to next state. Selected actions are passed if applicable.
	// Set to nil to make this a terminal state.
	Transition func(map[Seat]Action) (*State, error)
}

type Action interface {
	ActionOrder() uint64
}

type Seat int