package mahjong

import (
	"github.com/roelofruis/mahjong-learn/state_machine"
	"math/rand"
	"testing"
)

func Benchmark100GameRuns(b *testing.B) {
	rand.Seed(0)

	for i := 0; i < 100; i++ {
		game := NewGame(uint64(i))

		for {
			if game.StateMachine.HasTerminated() {
				break
			}

			selectedActions := make(map[state_machine.Seat]int)
			for seat, a := range game.StateMachine.AvailableActions() {
				selectedActions[seat] = rand.Intn(len(a))
			}

			_ = game.StateMachine.Transition(selectedActions)
		}
	}
}