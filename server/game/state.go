package game

import (
	"fmt"
	"sort"
)

func (m *StateMachine) Transition(selectedActions map[Seat]int) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.state.PlayerActions == nil {
		// no player actions required
		if m.state.Transfer == nil {
			return nil // end state
		}

		for {
			// move forward without selecting actions
			state, err := m.state.Transfer(m.game, nil)
			if err != nil {
				return err
			}
			m.state = state

			if m.state.PlayerActions != nil {
				// move forward until we are in a state where a player action is required
				return nil
			}
		}
	}

	if selectedActions == nil {
		// actions are required but none provided
		return fmt.Errorf("a nil actions map was provided")
	}

	// collect selected actions for all players
	pickedActions := make(map[Seat]Action)
	for seat, actions := range m.state.PlayerActions(m.game) {
		selected, has := selectedActions[seat]
		if !has {
			return fmt.Errorf("state requires action for seat [%d] but no action was given", seat)
		}
		if selected < 0 || selected >= len(actions) {
			return fmt.Errorf("selected action for seat [%d] is out of range (%d not in 0 to %d)", seat, selected, len(actions)-1)
		}
		pickedActions[seat] = actions[selected].Action
	}

	state, err := m.state.Transfer(m.game, pickedActions)
	if err != nil {
		return err
	}
	m.state = state

	return nil
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

var StateNewGame = &State{
	Name:          "New Game",
	PlayerActions: nil,
	Transfer:      Initialize,
}

var StateNextRound = &State{
	Name:          "Next Round",
	PlayerActions: nil,
	Transfer:      nil,
}

var StateNextTurn = &State{
	Name:          "Next Turn",
	PlayerActions: nil,
	Transfer:      TryDealTile,
}

var StateTileReceived = &State{
	Name:          "Tile Received",
	PlayerActions: TileReceivedReactions,
	Transfer:      HandleTileReceived,
}

var StateTileDiscarded = &State{
	Name:          "Tile Discarded",
	PlayerActions: TileDiscardedReactions,
	Transfer:      nil,
}

func Initialize(g *Game, _ map[Seat]Action) (*State, error) {
	g.DealTiles(13, 0)
	g.DealTiles(13, 1)
	g.DealTiles(13, 2)
	g.DealTiles(13, 3)

	return StateNextTurn, nil
}

func TryDealTile(g *Game, _ map[Seat]Action) (*State, error) {
	if g.Wall.Size() <= 14 {
		// tally scores?
		return StateNextRound, nil
	}

	g.DealTiles(1, g.ActiveSeat)

	return StateTileReceived, nil
}

func TileReceivedReactions(g *Game) map[Seat][]PlayerAction {
	m := make(map[Seat][]PlayerAction, 1)

	a := make([]PlayerAction, 0)
	for t, c := range g.Players[g.ActiveSeat].Concealed.Tiles {
		if c < 1 {
			continue
		}

		a = append(a, PlayerAction{
			Name:           fmt.Sprintf("Discard a %s", TileNames[t]),
			Action: Discard{Tile: t},
		})
	}

	// TODO: check whether player can declare kong or mahjong and add to available actions

	sort.Sort(ByIndex(a))

	m[g.ActiveSeat] = a

	return m
}

func HandleTileReceived(g *Game, actions map[Seat]Action) (*State, error) {
	switch a := actions[g.ActiveSeat].(type) {
	case Discard:
		g.Players[g.ActiveSeat].Concealed.Transfer(a.Tile, g.Players[g.ActiveSeat].Discarded)
		return StateTileDiscarded, nil
	default:
		return nil, fmt.Errorf("illegal action %+v", a)
	}
}

func TileDiscardedReactions(g *Game) map[Seat][]PlayerAction {
	m := make(map[Seat][]PlayerAction, 4)

	for i := 0; i < 4; i++ {
		a := make([]PlayerAction, 0)

		a = append(a, PlayerAction{
			Name: fmt.Sprintf("Do nothing"),
			Action: DoNothing{},
		})

		// TODO: check whether player can declare pung, kong, chow or mahjong and add to available actions

		m[Seat(i)] = a
	}

	return m
}
