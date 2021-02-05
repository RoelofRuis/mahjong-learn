package mahjong

import (
	"fmt"
	"github.com/roelofruis/mahjong-learn/state"
	"math/rand"
	"strings"
	"testing"
)

func TestGameLogic(t *testing.T) {
	transitioner := &state.DebugTransitioner{TransitionLimit: 10}
	game, _ := NewGame(transitioner)

	numTransitions := 0
	var stateHistory []string
	for {
		actions := game.StateMachine.AvailableActions()

		stateHistory = append(stateHistory, game.StateMachine.StateName())

		if game.StateMachine.HasTerminated() {
			break
		}

		if actions == nil {
			t.Logf("state after transition should define some actions")
			t.FailNow()
		}

		selectedActions := make(map[state.Seat]int)
		for seat, a := range game.StateMachine.AvailableActions() {
			selectedActions[seat] = rand.Intn(len(a))
		}

		err := game.StateMachine.Transition(selectedActions)
		if err != nil {
			t.Logf("game transition raised an error: %s", err.Error())
			t.FailNow()
		}

		numTransitions++

		err = checkInvariants(*game.Table)
		if err != nil {
			t.Logf("invariant failed after [%d] transitions: %s", numTransitions, err.Error())
			t.Logf("%s\n", describeState(transitioner))
			t.FailNow()
		}
	}

	t.Logf("game ended without errors after [%d] actions", numTransitions)
}

func checkInvariants(table Table) error {
	return checkTileCount(table)
}

func checkTileCount(table Table) error {
	wall := table.GetWall().Size()
	discard := 0
	if table.GetActiveDiscard() != nil {
		discard = 1
	}
	player1 := countPlayerTiles(table.GetPlayerAtSeat(0))
	player2 := countPlayerTiles(table.GetPlayerAtSeat(1))
	player3 := countPlayerTiles(table.GetPlayerAtSeat(2))
	player4 := countPlayerTiles(table.GetPlayerAtSeat(3))

	tileCount := wall + discard + player1 + player2 + player3 + player4

	if tileCount != 144 {
		return fmt.Errorf("incorrect tile count [%d]", tileCount)
	}

	return nil
}

func countPlayerTiles(player *Player) int {
	concealed := player.GetConcealedTiles().Size()
	discarded := player.GetDiscardedTiles().Size()
	received := 0
	if player.GetReceivedTile() != nil {
		received = 1
	}
	exposed := 0
	for _, c := range player.GetExposedCombinations() {
		switch c.(type) {
		case Kong:
			exposed += 4
		case BonusTile:
			exposed += 1
		default:
			exposed += 3
		}
	}
	return concealed + discarded + exposed + received
}

func describeState(transitioner *state.DebugTransitioner) string {
	var seatActions []string
	for seat, actions := range transitioner.LastActions {
		var actionNames []string
		for _, a := range actions {
			actionNames = append(actionNames, describeAction(a))
		}
		seatActions = append(seatActions, fmt.Sprintf("seat [%d] : %s", seat, strings.Join(actionNames, ",")))
	}
	return strings.Join(seatActions, "\n")
}

func describeAction(action state.Action) string {
	switch a := action.(type) {
	case Discard:
		return fmt.Sprintf("Discard [%d]", a.Tile)
	case DeclareConcealedKong:
		return fmt.Sprintf("Declare a concealed Kong [%d]", a.Tile)
	case ExposedPungToKong:
		return fmt.Sprintf("Add to exposed pung")
	case DoNothing:
		return "Do nothing"
	case DeclareChow:
		return fmt.Sprintf("Declare chow [%d]", a.Tile)
	case DeclarePung:
		return "Declare a pung"
	case DeclareKong:
		return "Declare a kong"
	case DeclareMahjong:
		return "Declare mahjong"
	}
	return "<UNKNOWN>"
}
