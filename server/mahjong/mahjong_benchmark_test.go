package mahjong

import (
	"github.com/roelofruis/mahjong-learn/state"
	"math/rand"
	"testing"
)

func Benchmark100GameRuns(b *testing.B) {
	rand.Seed(0)

	for i := 0; i < 100; i++ {
		game, _ := NewGame(&state.ProductionTransitioner{TransitionLimit: 10})

		for {
			if game.StateMachine.HasTerminated() {
				break
			}

			selectedActions := make(map[state.Seat]int)
			for seat, a := range game.StateMachine.AvailableActions() {
				selectedActions[seat] = rand.Intn(len(a))
			}

			_ = game.StateMachine.Transition(selectedActions)
		}
	}
}
