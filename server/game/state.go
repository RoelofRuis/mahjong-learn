package game

import (
	"fmt"
	"sort"
)

var (
	stateNewGame       *State
	stateNextRound     *State
	stateNextTurn      *State
	stateTileReceived  *State
	stateTileDiscarded *State
)

func init() {
	stateNewGame = &State{
		Name:          "New Game",
		PlayerActions: nil,
		Transition:    initialize,
	}

	stateNextRound = &State{
		Name:          "Next Round",
		PlayerActions: nil,
		Transition:    nil,
	}

	stateNextTurn = &State{
		Name:          "Next Turn",
		PlayerActions: nil,
		Transition:    tryDealTile,
	}

	stateTileReceived = &State{
		Name:          "Tile Received",
		PlayerActions: tileReceivedActions,
		Transition:    handleTileReceivedActions,
	}

	stateTileDiscarded = &State{
		Name:          "Tile Discarded",
		PlayerActions: tileDiscardedActions,
		Transition:    handleTileDiscardedActions,
	}
}

func (m *StateMachine) Transition(selectedActions map[Seat]int) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.state.Transition == nil {
		return nil // we cannot transition from this state
	}

	playerActions := make(map[Seat]Action)

	if m.state.PlayerActions != nil {
		if selectedActions == nil {
			return fmt.Errorf("a nil actions map was provided")
		}

		for seat, actions := range m.state.PlayerActions(m.game) {
			selected, has := selectedActions[seat]
			if !has {
				return fmt.Errorf("state requires action for seat [%d] but no action was given", seat)
			}
			if selected < 0 || selected >= len(actions) {
				return fmt.Errorf("selected action for seat [%d] is out of range (%d not in 0 to %d)", seat, selected, len(actions)-1)
			}
			playerActions[seat] = actions[selected].Action
		}
	}

	for {
		state, err := m.state.Transition(m.game, playerActions)
		if err != nil {
			return err
		}
		m.state = state

		// transition until we are in a state where another player action is required
		if m.state.PlayerActions != nil {
			return nil
		}
	}
}

// If the StateMachine is viewed, internals should be exposed in a consistent manner, so one function returns everything.
func (m *StateMachine) View() (Game, State, map[Seat][]PlayerAction) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var playerActions = make(map[Seat][]PlayerAction)
	if m.state.PlayerActions != nil {
		playerActions = m.state.PlayerActions(m.game)
	}

	return *m.game, *m.state, playerActions
}

func initialize(g *Game, _ map[Seat]Action) (*State, error) {
	g.DealTiles(13, 0)
	g.DealTiles(13, 1)
	g.DealTiles(13, 2)
	g.DealTiles(13, 3)

	return stateNextTurn, nil
}

func tryDealTile(g *Game, _ map[Seat]Action) (*State, error) {
	if g.Wall.Size() <= 14 {
		// tally scores?
		return stateNextRound, nil
	}

	g.DealTiles(1, g.ActiveSeat)

	return stateTileReceived, nil
}

func tileReceivedActions(g *Game) map[Seat][]PlayerAction {
	m := make(map[Seat][]PlayerAction, 1)

	a := make([]PlayerAction, 0)
	for t, c := range g.Players[g.ActiveSeat].Concealed.Tiles {
		if c < 1 {
			continue
		}

		a = append(a, PlayerAction{
			Name:   fmt.Sprintf("Discard a %s", TileNames[t]),
			Action: Discard{Tile: t},
		})
	}

	// TODO: check whether player can declare kong or mahjong and add to available actions

	sort.Sort(ByIndex(a))

	m[g.ActiveSeat] = a

	return m
}

func handleTileReceivedActions(g *Game, actions map[Seat]Action) (*State, error) {
	switch a := actions[g.ActiveSeat].(type) {
	case Discard:
		g.Players[g.ActiveSeat].Concealed.Transfer(a.Tile, g.Players[g.ActiveSeat].Discarded)
		return stateTileDiscarded, nil

		// TODO: handle other possible cases

	default:
		return nil, fmt.Errorf("illegal action %+v", a)
	}
}

func tileDiscardedActions(g *Game) map[Seat][]PlayerAction {
	m := make(map[Seat][]PlayerAction, 4)

	for i := 0; i < 4; i++ {
		a := make([]PlayerAction, 0)

		a = append(a, PlayerAction{
			Name:   fmt.Sprintf("Do nothing"),
			Action: DoNothing{},
		})

		// TODO: check whether player can declare pung, kong, chow or mahjong and add to available actions

		m[Seat(i)] = a
	}

	return m
}

func handleTileDiscardedActions(g *Game, actions map[Seat]Action) (*State, error) {
	// TODO: handle actions

	g.ActiveSeat = Seat(int(g.ActiveSeat) + 1%4)

	return stateNextTurn, nil
}
