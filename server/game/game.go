package game

import (
	"fmt"
	"github.com/roelofruis/mahjong-learn/game/model"
	"sync"
)

type Game struct {
	lock sync.RWMutex

	id uint64

	transitionLimit int

	state *State
	table *model.Table
}

type State struct {
	// Name just to display human readable information.
	Name string

	// Required player actions. May be nil if this state requires no player actions.
	PlayerActions func(*model.Table) map[model.Seat][]model.Action

	// Transition to next state. Selected actions are passed if applicable.
	Transition func(*model.Table, map[model.Seat]model.Action) (*State, error)
}

func NewGame(id uint64) *Game {
	return &Game{
		id: id,

		transitionLimit: 10,

		state: stateNewGame,
		table: model.NewTable(),
	}
}

func (m *Game) Id() uint64 {
	return m.id
}

func (m *Game) Transition(selectedActions map[model.Seat]int) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.state.Transition == nil {
		return nil
	}

	playerActions := make(map[model.Seat]model.Action)

	if m.state.PlayerActions != nil {
		if selectedActions == nil {
			return fmt.Errorf("a nil actions map was provided")
		}

		for seat, actions := range m.state.PlayerActions(m.table) {
			selected, has := selectedActions[seat]
			if !has {
				return fmt.Errorf("state requires action for seat [%d] but no action was given", seat)
			}
			if selected < 0 || selected >= len(actions) {
				return fmt.Errorf("selected action for seat [%d] is out of range (%d not in 0 to %d)", seat, selected, len(actions)-1)
			}
			playerActions[seat] = actions[selected]
		}
	}

	var stateHistory []string
	for {
		state, err := m.state.Transition(m.table, playerActions)
		if err != nil {
			return err
		}
		m.state = state
		playerActions = nil // only use player actions in first transition

		// transition until we are in a terminal state, or another player action is required
		if m.state.Transition == nil || m.state.PlayerActions != nil {
			return nil
		}

		stateHistory = append(stateHistory, m.state.Name)
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

// If the Game is viewed, internals should be exposed in a consistent manner, so one function returns everything.
func (m *Game) View() (model.Table, State, map[model.Seat][]model.Action) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var playerActions = make(map[model.Seat][]model.Action)
	if m.state.PlayerActions != nil {
		playerActions = m.state.PlayerActions(m.table)
	}

	return *m.table, *m.state, playerActions
}