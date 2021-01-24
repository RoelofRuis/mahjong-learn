package game

import (
	"fmt"
	"sort"
)

func (m *StateMachine) Transition(selectedActions map[Seat]int) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.state.TransferAction == nil {
		return m.executePlayerActions(selectedActions)
	}

	for {
		if m.state.TransferAction == nil {
			return nil
		}
		m.state = m.state.TransferAction(m.game)
	}
}

func (m *StateMachine) executePlayerActions(selectedActions map[Seat]int) error {
	if selectedActions == nil {
		return fmt.Errorf("a nil actions map was provided")
	}

	available := m.state.PlayerActions(m.game)

	pickedActions := make(map[Seat]PlayerAction)
	for seat, actions := range available {
		selected, has := selectedActions[seat]
		if !has {
			return fmt.Errorf("state requires action for seat [%d] but no action was given", seat)
		}
		if selected < 0 || selected >= len(actions) {
			return fmt.Errorf("selected action for seat [%d] is out of range (%d not in 0 to %d)", seat, selected, len(actions)-1)
		}
		pickedActions[seat] = actions[selected]
	}

	if len(pickedActions) == 1 {
		for _, a := range pickedActions {
			m.state = a.TransferAction(m.game)
			return nil
		}
	}

	// TODO: if multiple seats declared an action, determine which actions to execute (maybe len 1 case can be merged eventually)

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
	Name:           "New Game",
	TransferAction: Initialize,
	PlayerActions:  nil,
}

var StateNextRound = &State{
	Name:           "Next Round",
	TransferAction: nil,
	PlayerActions:  nil,
}

var StateNextTurn = &State{
	Name:           "Next Turn",
	TransferAction: TryDealTile,
	PlayerActions:  nil,
}

var StateTileReceived = &State{
	Name:           "Tile Received",
	TransferAction: nil,
	PlayerActions:  ReactToTile,
}

var StateTileDiscarded = &State{
	Name:           "Tile Discarded",
	TransferAction: nil,
	PlayerActions:  ReactToDiscard,
}

func Initialize(g *Game) *State {
	g.DealTiles(13, 0)
	g.DealTiles(13, 1)
	g.DealTiles(13, 2)
	g.DealTiles(13, 3)
	return StateNextTurn
}

func TryDealTile(g *Game) *State {
	if g.Wall.Size() <= 14 {
		// tally scores?
		return StateNextRound
	}

	g.DealTiles(1, g.ActiveSeat)
	return StateTileReceived
}

func ReactToTile(g *Game) map[Seat][]PlayerAction {
	m := make(map[Seat][]PlayerAction, 1)

	a := make([]PlayerAction, 0)
	for t, c := range g.Players[g.ActiveSeat].Concealed.Tiles {
		if c < 1 {
			continue
		}

		a = append(a, PlayerAction{
			Index:          int(t),
			Name:           fmt.Sprintf("Discard a %s", TileNames[t]),
			TransferAction: DiscardTile(t),
		})
	}

	// TODO: check whether player can declare kong or mahjong and add to available actions

	sort.Sort(ByIndex(a))

	m[g.ActiveSeat] = a

	return m
}

func DiscardTile(tile Tile) func(g *Game) *State {
	return func(g *Game) *State {
		g.Players[g.ActiveSeat].Concealed.Transfer(tile, g.Players[g.ActiveSeat].Discarded)

		return StateTileDiscarded
	}
}

func ReactToDiscard(g *Game) map[Seat][]PlayerAction {
	m := make(map[Seat][]PlayerAction, 4)

	// TODO: implement

	return m
}
