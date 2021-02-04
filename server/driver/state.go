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
	ActionOrder() int
}

type ByActionOrder []Action

func (a ByActionOrder) Len() int           { return len(a) }
func (a ByActionOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByActionOrder) Less(i, j int) bool { return a[i].ActionOrder() < a[j].ActionOrder() }

type Seat int