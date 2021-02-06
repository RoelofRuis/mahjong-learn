package state

import (
	"fmt"
)

type Transitioner interface {
	// Perform the transition to the next state based on the selected actions.
	//
	// This might return one of several errors:
	// IncorrectActionError in case the given action map is inconsistent with the currently available actions as returned by AvailableActions().
	// TooManyIntermediateStatesError in case the chain of intermediate states became too long.
	// TransitionLogicError in case executing the transition logic returned an error.
	Transition(machine *StateMachine, selectedActions map[int]int) error
}

type ProductionTransitioner struct {
	IntermediateTransitionLimit int
}

func (t *ProductionTransitioner) Transition(m *StateMachine, selectedActions map[int]int) error {
	playerActions := make(map[int]Action)

	if m.state.actions != nil {
		if selectedActions == nil {
			// initialize empty so we can return IncorrectActionError for erroneous player
			selectedActions = make(map[int]int, 0)
		}

		for player, actions := range m.state.actions {
			selected, has := selectedActions[player]
			if !has || selected < 0 || selected >= len(actions) {
				return IncorrectActionError{player: player, upperActionIndex: len(actions) - 1}
			}
			playerActions[player] = actions[selected]
		}
	}

	statesVisited := 0
	for {
		state, err := m.state.transition(playerActions)
		if err != nil {
			return TransitionLogicError{Err: err}
		}
		m.state = state
		playerActions = nil // only use player actions in first Transition

		if m.HasTerminated() || m.state.actions != nil {
			// Transition until we are in a terminal state, or another player action is required
			return nil
		}

		statesVisited++

		if statesVisited > t.IntermediateTransitionLimit {
			return TooManyIntermediateStatesError{
				transitionLimit: t.IntermediateTransitionLimit,
			}
		}
	}
}

type DebugTransitioner struct {
	IntermediateTransitionLimit int
	ActionLimit                 int
	actionsPerformed            int

	LastActions   map[int][]Action
	LastSelection map[int]int
}

func (t *DebugTransitioner) Transition(m *StateMachine, selectedActions map[int]int) error {
	t.actionsPerformed++
	if t.actionsPerformed > t.ActionLimit {
		return TransitionLogicError{Err: fmt.Errorf("more actions than the configured maximum of [%d] were performed", t.ActionLimit)}
	}

	t.LastActions = m.state.actions
	t.LastSelection = selectedActions

	playerActions := make(map[int]Action)

	if m.state.actions != nil {
		if selectedActions == nil {
			// initialize empty so we can return IncorrectActionError for erroneous player
			selectedActions = make(map[int]int, 0)
		}

		for player, actions := range m.state.actions {
			selected, has := selectedActions[player]
			if !has || selected < 0 || selected >= len(actions) {
				return IncorrectActionError{player: player, upperActionIndex: len(actions) - 1}
			}
			playerActions[player] = actions[selected]
		}
	}

	statesVisited := 0
	for {
		state, err := m.state.transition(playerActions)
		if err != nil {
			return TransitionLogicError{Err: err}
		}
		m.state = state
		playerActions = nil // only use player actions in first Transition

		if m.HasTerminated() || m.state.actions != nil {
			// Transition until we are in a terminal state, or another player action is required
			return nil
		}

		statesVisited++

		if statesVisited > t.IntermediateTransitionLimit {
			return TooManyIntermediateStatesError{
				transitionLimit: t.IntermediateTransitionLimit,
			}
		}
	}
}
