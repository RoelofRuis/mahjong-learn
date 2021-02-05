package state

type stateTransitioner interface {
	// Perform the transition to the next state
	//
	// Might return one of several errors:
	// IncorrectActionError in case the given action map is inconsistent with the currently available actions as returned by AvailableActions()
	// TooManyIntermediateStatesError in case the chain of state transitions that did not require an action became too long
	// TransitionLogicError in case executing the transition logic returned an error.
	Transition(machine *StateMachine, selectedActions map[Seat]int) error
}

type productionTransitioner struct {
	transitionLimit int
}

func (t *productionTransitioner) Transition(m *StateMachine, selectedActions map[Seat]int) error {
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

	statesVisited := 0
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

		statesVisited++

		if statesVisited > t.transitionLimit {
			return TooManyIntermediateStatesError{
				transitionLimit: t.transitionLimit,
			}
		}
	}
}
