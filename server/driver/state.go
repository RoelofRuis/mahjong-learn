package driver

// indicates a player seat number
type Seat int

type State struct {
	// name just to display human readable information.
	name string

	// Required actions per seat. May be set to nil if this state requires no actions.
	actions map[Seat][]Action

	// transition to next state. Selected actions are passed if applicable.
	// Set to nil to make this a terminal state.
	transition func(map[Seat]Action) (*State, error)
}

type Action interface {
	// Defines an order for the actions returned.
	// Needs to be unique among simultaneous action options to guarantee a stable sorting.
	ActionOrder() int
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

type byActionOrder []Action

func (a byActionOrder) Len() int           { return len(a) }
func (a byActionOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byActionOrder) Less(i, j int) bool { return a[i].ActionOrder() < a[j].ActionOrder() }
