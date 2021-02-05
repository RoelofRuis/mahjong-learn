package state

import (
	"sort"
	"sync"
)

type productionStateMachine struct {
	lock sync.Mutex

	transitionLimit int

	state *State
}

func (m *productionStateMachine) StateName() string {
	return m.state.name
}

func (m *productionStateMachine) AvailableActions() map[Seat][]Action {
	if m.state.actions == nil {
		return make(map[Seat][]Action)
	}

	for s, a := range m.state.actions {
		sort.Sort(byActionOrder(a))
		m.state.actions[s] = a
	}

	return m.state.actions
}

func (m *productionStateMachine) HasTerminated() bool {
	return m.state.transition == nil
}

func (m *productionStateMachine) Transition(selectedActions map[Seat]int) error {
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
				return IncorrectActionError{seat: seat, upperActionIndex: len(actions) - 1}
			}
			seatActions[seat] = actions[selected]
		}
	}

	var stateHistory []string
	for {
		state, err := m.state.transition(seatActions)
		if err != nil {
			return TransitionLogicError{Err: err}
		}
		m.state = state
		seatActions = nil // only use player actions in first transition

		if m.HasTerminated() || m.state.actions != nil {
			// transition until we are in a terminal state, or another player action is required
			return nil
		}

		stateHistory = append(stateHistory, m.StateName())
		if len(stateHistory) > m.transitionLimit {
			return TooManyIntermediateStatesError{
				transitionLimit: m.transitionLimit,
				StateHistory:    stateHistory,
			}
		}
	}
}
