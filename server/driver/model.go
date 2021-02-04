package driver

import (
	"fmt"
	"sync"
)

type GameDriver interface {
	// Name of the current state the driver is in
	StateName() string

	// Whether the driver is in a terminal state and no more actions can be performed.
	// If this returns true, calling Transition is a no-op.
	HasTerminated() bool

	// Get the actions that are available for executing in this state
	AvailableActions() map[Seat][]Action

	// Perform the transition to the next state
	//
	// Might return one of several errors:
	// IncorrectActionError in case the given action map is inconsistent with the currently available actions as returned by AvailableActions()
	// TransitionLimitReachedError in case the chain of state transitions that did not require an action became too long
	// GameLogicError in case executing the game logic returned an error.
	Transition(selectedActions map[Seat]int) error
}

func NewGameDriver(initialState *State, transitionLimit int) GameDriver {
	return &productionGameDriver{
		lock:            sync.Mutex{},
		transitionLimit: transitionLimit,
		state:           initialState,
	}
}

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

type IncorrectActionError struct {
	seat Seat
	upperActionIndex int
}

func (e IncorrectActionError) Error() string {
	return fmt.Sprintf("an action is required for seat [%d] within range [0 to %d]", e.seat, e.upperActionIndex)
}

type TransitionLimitReachedError struct {
	transitionLimit int
	StateHistory []string
}

func (e TransitionLimitReachedError) Error() string {
	return fmt.Sprintf("transitioning to next action state took more than [%d] steps", e.transitionLimit)
}

type GameLogicError struct {
	Err error
}

func (e GameLogicError) Error() string {
	return fmt.Sprintf("game logic error: %s", e.Err.Error())
}