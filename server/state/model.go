package state

import (
	"fmt"
	"sort"
	"sync"
)

type StateMachine struct {
	lock sync.Mutex

	state *State

	transitioner Transitioner
}

// Name of the current state the state machine is in
func (s *StateMachine) StateName() string {
	return s.state.name
}

// Whether the state machine is in a terminal state and no more actions can be performed.
// If this returns true, calling Transition is a no-op.
func (s *StateMachine) HasTerminated() bool {
	return s.state.transition == nil
}

// Get the actions that are available for executing in this state per player
func (s *StateMachine) AvailableActions() map[int][]Action {
	if s.state.actions == nil {
		return make(map[int][]Action)
	}

	for player, a := range s.state.actions {
		sort.Sort(byActionOrder(a))
		s.state.actions[player] = a
	}

	return s.state.actions
}

func (s *StateMachine) Transition(selectedActions map[int]int) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.HasTerminated() {
		return nil
	}

	return s.transitioner.Transition(s, selectedActions)
}

func NewStateMachine(initialState *State, transitioner Transitioner) *StateMachine {
	return &StateMachine{
		lock:         sync.Mutex{},
		state:        initialState,
		transitioner: transitioner,
	}
}

type State struct {
	// name just to display human readable information.
	name string

	// Required actions per player. May be set to nil if this state requires no actions.
	actions map[int][]Action

	// transition to next state. Selected actions are passed if applicable.
	// Set to nil to make this a terminal state.
	transition func(map[int]Action) (*State, error)
}

type Action interface {
	// Defines an order for the actions returned.
	// Needs to be unique among simultaneous action options to guarantee a stable sorting.
	ActionOrder() int
}

func NewState(name string, actions map[int][]Action, transition func(map[int]Action) (*State, error)) *State {
	return &State{
		name:       name,
		actions:    actions,
		transition: transition,
	}
}

func NewIntermediateState(name string, transition func() *State) *State {
	return &State{
		name:       name,
		actions:    nil,
		transition: func(_ map[int]Action) (*State, error) { return transition(), nil },
	}
}

func NewTerminalState(name string) *State {
	return &State{
		name:       name,
		actions:    nil,
		transition: nil,
	}
}

type byActionOrder []Action

func (a byActionOrder) Len() int           { return len(a) }
func (a byActionOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byActionOrder) Less(i, j int) bool { return a[i].ActionOrder() < a[j].ActionOrder() }

type IncorrectActionError struct {
	player           int
	upperActionIndex int
}

func (e IncorrectActionError) Error() string {
	return fmt.Sprintf("an action is required for player [%d] within range [0 to %d]", e.player, e.upperActionIndex)
}

type TooManyIntermediateStatesError struct {
	transitionLimit int
}

func (e TooManyIntermediateStatesError) Error() string {
	return fmt.Sprintf("transitioning to next actionable state took more than [%d] steps", e.transitionLimit)
}

type TransitionLogicError struct {
	Err error
}

func (e TransitionLogicError) Error() string {
	return fmt.Sprintf("transition logic error: %s", e.Err.Error())
}
