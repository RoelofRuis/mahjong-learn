package game

import (
	"github.com/roelofruis/mahjong-learn/game/model"
	"math/rand"
	"testing"
)

func Benchmark100GameRuns(b *testing.B) {
	rand.Seed(0)

	for i := 0; i < 100; i++ {
		game := NewGame(uint64(i))

		var _, state, actions = game.View()

		for {
			if state.Transition == nil {
				break
			}

			selectedActions := make(map[model.Seat]int)
			for seat, actions := range actions {
				selectedActions[seat] = rand.Intn(len(actions))
			}

			_ = game.Transition(selectedActions)
			_, state, actions = game.View()
		}
	}
}