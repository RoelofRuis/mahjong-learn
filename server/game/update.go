package game

import "math/rand"

func (g *Game) DealStartingHands() {
	g.Wall.Transfer(13, g.Players[0].Concealed)
	g.Wall.Transfer(13, g.Players[1].Concealed)
	g.Wall.Transfer(13, g.Players[2].Concealed)
	g.Wall.Transfer(13, g.Players[3].Concealed)
}

func (g *Game) DealTile(seat Seat) {
	g.Wall.Transfer(1, g.Players[seat].Concealed)
}

// Transfers n randomly picked tiles from this tile collection to the target tile collection.
func (t *TileCollection) Transfer(n int, target *TileCollection) {
	var tileList = make([]Tile, 0)
	for k, v := range t.Tiles {
		for i := v; i > 0; i-- {
			tileList = append(tileList, k)
		}
	}
	for i := n; i > 0; i-- {
		numTiles := len(tileList)
		pos := rand.Intn(numTiles)
		picked := tileList[pos]

		tileList[pos] = tileList[numTiles-1]
		tileList = tileList[:numTiles-1]

		t.Tiles[picked]--
		target.Tiles[picked]++
	}
}
