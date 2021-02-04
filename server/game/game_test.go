package game

import (
	"fmt"
	"github.com/roelofruis/mahjong-learn/driver"
	"math/rand"
	"testing"
)

func TestGameLogic(t *testing.T) {
	game := NewMahjongGame(1)

	numTransitions := 0
	var stateHistory []string
	for {
		state := game.Driver.GetState()
		actions := game.Driver.GetState().Actions

		stateHistory = append(stateHistory, state.Name)

		if state.Transition == nil {
			break
		}

		if actions == nil {
			t.Logf("state after transition should define some actions")
			t.FailNow()
		}

		selectedActions := make(map[driver.Seat]int)
		for seat, a := range actions() {
			selectedActions[seat] = rand.Intn(len(a))
		}

		err := game.Driver.Transition(selectedActions)
		if err != nil {
			t.Logf("game transition raised an error: %s", err.Error())
			t.FailNow()
		}

		numTransitions++

		err = checkInvariants(*game.Table)
		if err != nil {
			t.Logf("invariant failed after [%d] transitions: %s", numTransitions, err.Error())
			t.Logf("state was [%s] selected actions were [%+v]", state.Name, selectedActions)
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
	if player.GetReceivedTile() != nil{
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