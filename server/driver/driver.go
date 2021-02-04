package driver

import (
	"sort"
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

type productionGameDriver struct {
	lock sync.Mutex

	transitionLimit int

	state *State
}

func NewGameDriver(initialState *State, transitionLimit int) GameDriver {
	return &productionGameDriver{
		lock:            sync.Mutex{},
		transitionLimit: transitionLimit,
		state:           initialState,
	}
}

func (m *productionGameDriver) StateName() string {
	return m.state.name
}

func (m *productionGameDriver) AvailableActions() map[Seat][]Action {
	if m.state.actions == nil {
		return make(map[Seat][]Action)
	}

	for s, a := range m.state.actions {
		sort.Sort(byActionOrder(a))
		m.state.actions[s] = a
	}

	return m.state.actions
}

func (m *productionGameDriver) HasTerminated() bool {
	return m.state.transition == nil
}

func (m *productionGameDriver) Transition(selectedActions map[Seat]int) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.HasTerminated() {
		return nil
	}

	seatActions := make(map[Seat]Action)

	if m.state.actions != nil {
		if selectedActions == nil {
			// initialize empty so we can return IncorrectActionError for erroneous seat
			selectedActions = make(map[Seat]int, 0)
		}

		for seat, actions := range m.state.actions {
			selected, has := selectedActions[seat]
			if !has || selected < 0 || selected >= len(actions) {
				return IncorrectActionError{seat: seat, upperActionIndex: len(actions) -1}
			}
			seatActions[seat] = actions[selected]
		}
	}

	var stateHistory []string
	for {
		state, err := m.state.transition(seatActions)
		if err != nil {
			return GameLogicError{Err: err}
		}
		m.state = state
		seatActions = nil // only use player actions in first transition

		if m.HasTerminated() || m.state.actions != nil {
			// transition until we are in a terminal state, or another player action is required
			return nil
		}

		stateHistory = append(stateHistory, m.StateName())
		if len(stateHistory) > m.transitionLimit {
			return TransitionLimitReachedError{
				transitionLimit: m.transitionLimit,
				StateHistory:    stateHistory,
			}
		}
	}
}
