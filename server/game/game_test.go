package game

import (
	"github.com/roelofruis/mahjong-learn/game/model"
	"math/rand"
	"testing"
)

// TODO: expand with invariant checking

func TestGameLogic(t *testing.T) {
	game := NewGame(1)
	err := game.Transition(nil)
	if err != nil {
		t.Errorf("game transition raised an error: %s", err.Error())
	}

	numTransitions := 0
	for {
		_, state, actions := game.View()

		if state.Transition == nil {
			break
		}

		if actions == nil {
			t.Errorf("state after transition should define some actions")
		}

		selectedActions := make(map[model.Seat]int)
		for seat, actions := range actions {
			selectedActions[seat] = rand.Intn(len(actions))
		}

		err := game.Transition(selectedActions)
		if err != nil {
			t.Errorf("game transition raised an error: %s", err.Error())
		}

		numTransitions++
	}

	t.Logf("game ended without errors after [%d] actions", numTransitions)
}
