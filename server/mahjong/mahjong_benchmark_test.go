package mahjong

import (
	"github.com/roelofruis/mahjong-learn/state"
	"math/rand"
	"testing"
)

func Benchmark1kRuns(b *testing.B) {
	rand.Seed(0)

	for i := 0; i < b.N; i++ {
		game, _ := NewGame(&state.ProductionTransitioner{TransitionLimit: 10})

		for {
			if game.StateMachine.HasTerminated() {
				break
			}

			selectedActions := make(map[int]int)
			for player, a := range game.StateMachine.AvailableActions() {
				selectedActions[player] = rand.Intn(len(a))
			}

			_ = game.StateMachine.Transition(selectedActions)
		}
	}
}
