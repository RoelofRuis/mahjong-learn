package driver

import (
	"fmt"
	"sync"
)

type Seat int

type GameDriver struct {
	lock sync.RWMutex

	id uint64

	transitionLimit int

	state State
}

type Action interface {
	ActionOrder() uint64
}

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

func (m *GameDriver) Transition(selectedActions map[Seat]int) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.state.IsTerminal() {
		return nil
	}

	seatActions := make(map[Seat]Action)
	possibleActions := m.state.Actions()

	if possibleActions != nil {
		if selectedActions == nil {
			return fmt.Errorf("state requires actions but a nil actions map was provided")
		}

		for seat, actions := range possibleActions {
			selected, has := selectedActions[seat]
			if !has {
				return fmt.Errorf("state requires action for seat [%d] but no action was given", seat)
			}
			if selected < 0 || selected >= len(actions) {
				return fmt.Errorf("selected action for seat [%d] is out of range (%d not in 0 to %d)", seat, selected, len(actions)-1)
			}
			seatActions[seat] = actions[selected]
		}
	}

	var stateHistory []string
	for {
		state, err := m.state.Transition(seatActions)
		if err != nil {
			return err
		}
		m.state = state
		seatActions = nil // only use player actions in first transition

		// transition until we are in a terminal state, or another player action is required
		if m.state.IsTerminal() || m.state.Actions() != nil {
			return nil
		}

		stateHistory = append(stateHistory, m.state.Name())
		if len(stateHistory) > m.transitionLimit {
			// there is probably some infinite loop in the transition logic...
			stateDebug := ""
			for _, s := range stateHistory {
				stateDebug += fmt.Sprintf("%s\n", s)
			}
			return fmt.Errorf("transitioning took more than %d steps. "+
				"There is probably an infinite loop in the state transitions.\nVisited stateds were:\n%s", m.transitionLimit, stateDebug)
		}
	}
}
