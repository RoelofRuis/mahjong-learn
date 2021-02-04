package mahjong

import (
	"github.com/roelofruis/mahjong-learn/driver"
	"math/rand"
	"testing"
)

func Benchmark100GameRuns(b *testing.B) {
	rand.Seed(0)

	for i := 0; i < 100; i++ {
		game := NewGame(uint64(i))

		for {
			if game.Driver.HasTerminated() {
				break
			}

			selectedActions := make(map[driver.Seat]int)
			for seat, a := range game.Driver.AvailableActions() {
				selectedActions[seat] = rand.Intn(len(a))
			}

			_ = game.Driver.Transition(selectedActions)
		}
	}
}