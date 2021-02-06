package state

type Transitioner interface {
	// Perform the transition to the next state
	//
	// Might return one of several errors:
	// IncorrectActionError in case the given action map is inconsistent with the currently available actions as returned by AvailableActions()
	// TooManyIntermediateStatesError in case the chain of state transitions that did not require an action became too long
	// TransitionLogicError in case executing the transition logic returned an error.
	Transition(machine *StateMachine, selectedActions map[Seat]int) error
}

type ProductionTransitioner struct {
	TransitionLimit int
}

func (t *ProductionTransitioner) Transition(m *StateMachine, selectedActions map[Seat]int) error {
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
		seatActions = nil // only use player actions in first Transition

		if m.HasTerminated() || m.state.actions != nil {
			// Transition until we are in a terminal state, or another player action is required
			return nil
		}

		statesVisited++

		if statesVisited > t.TransitionLimit {
			return TooManyIntermediateStatesError{
				transitionLimit: t.TransitionLimit,
			}
		}
	}
}

type DebugTransitioner struct {
	TransitionLimit int

	LastActions   map[Seat][]Action
	LastSelection map[Seat]int
}

func (t *DebugTransitioner) Transition(m *StateMachine, selectedActions map[Seat]int) error {
	t.LastActions = m.state.actions
	t.LastSelection = selectedActions

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
		seatActions = nil // only use player actions in first Transition

		if m.HasTerminated() || m.state.actions != nil {
			// Transition until we are in a terminal state, or another player action is required
			return nil
		}

		statesVisited++

		if statesVisited > t.TransitionLimit {
			return TooManyIntermediateStatesError{
				transitionLimit: t.TransitionLimit,
			}
		}
	}
}
