package game

import (
	"github.com/roelofruis/mahjong-learn/driver"
	"math/rand"
	"testing"
)

func Benchmark100GameRuns(b *testing.B) {
	rand.Seed(0)

	for i := 0; i < 100; i++ {
		game := NewMahjongGame(uint64(i))

		for {
			state := game.Driver.GetState()
			actions := state.Actions

			if state.Transition == nil {
				break
			}

			selectedActions := make(map[driver.Seat]int)
			for seat, a := range actions() {
				selectedActions[seat] = rand.Intn(len(a))
			}

			_ = game.Driver.Transition(selectedActions)
		}
	}
}