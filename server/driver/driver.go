package driver

import (
	"sync"
)

type GameDriver struct {
	lock sync.RWMutex

	transitionLimit int

	state State
}

func NewGameDriver(initialState State, transitionLimit int) *GameDriver {
	return &GameDriver{
		lock:            sync.RWMutex{},
		transitionLimit: transitionLimit,
		state:           initialState,
	}
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
			// initialize empty so we can return IncorrectActionError for erroneous seat
			selectedActions = make(map[Seat]int, 0)
		}

		for seat, actions := range possibleActions {
			selected, has := selectedActions[seat]
			if !has || selected < 0 || selected >= len(actions) {
				return IncorrectActionError{seat: seat, upperActionIndex: len(actions) -1}
			}
			seatActions[seat] = actions[selected]
		}
	}

	var stateHistory []string
	for {
		state, err := m.state.Transition(seatActions)
		if err != nil {
			return GameLogicError{Err: err}
		}
		m.state = state
		seatActions = nil // only use player actions in first transition

		// transition until we are in a terminal state, or another player action is required
		if m.state.IsTerminal() || m.state.Actions() != nil {
			return nil
		}

		stateHistory = append(stateHistory, m.state.Name())
		if len(stateHistory) > m.transitionLimit {
			return TransitionLimitReachedError{
				transitionLimit: m.transitionLimit,
				StateHistory:    stateHistory,
			}
		}
	}
}
