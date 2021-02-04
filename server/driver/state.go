package driver

type State struct {
	// name just to display human readable information.
	name string

	// Required actions per seat. May be set to nil if this state requires no actions.
	actions map[Seat][]Action

	// transition to next state. Selected actions are passed if applicable.
	// Set to nil to make this a terminal state.
	transition func(map[Seat]Action) (*State, error)
}

func NewState(name string, actions map[Seat][]Action, transition func(map[Seat]Action) (*State, error)) *State {
	return &State{
		name:       name,
		actions:    actions,
		transition: transition,
	}
}

func NewIntermediateState(name string, transition func(map[Seat]Action) (*State, error)) *State {
	return &State {
		name: name,
		actions: nil,
		transition: transition,
	}
}

func NewTerminalState(name string) *State {
	return &State {
		name: name,
		actions: nil,
		transition: nil,
	}
}

type Action interface {
	ActionOrder() int
}

type ByActionOrder []Action

func (a ByActionOrder) Len() int           { return len(a) }
func (a ByActionOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByActionOrder) Less(i, j int) bool { return a[i].ActionOrder() < a[j].ActionOrder() }

type Seat int