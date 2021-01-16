package game

import "math/rand"

func (g *Game) DealTiles(n int, seat Seat) {
	player := g.Players[seat]
	g.Wall.TransferRandom(n, player.Concealed)

	for {
		numExtra := player.ForceExposeTiles()
		if numExtra == 0 {
			break
		}

		g.Wall.TransferRandom(numExtra, player.Concealed)
	}
}

func (p *Player) ForceExposeTiles() int {
	var transferred = 0
	for t, c := range p.Concealed.Tiles {
		if IsBonusTile(t) && c > 0 {
			exposed := NewEmptyTileCollection()
			p.Concealed.Transfer(t, exposed)
			p.Exposed = append(p.Exposed, exposed)
			transferred++
		}
	}

	return transferred
}

func (t *TileCollection) Size() int {
	var count = 0
	for _, c := range t.Tiles {
		count += c
	}

	return count
}

func (t *TileCollection) Transfer(tile Tile, target *TileCollection) {
	n, has := t.Tiles[tile]
	if !has || n == 0{
		return
	}

	t.Tiles[tile]--
	target.Tiles[tile]++
}

// Transfers n randomly picked tiles from this tile collection to the target tile collection.
func (t *TileCollection) TransferRandom(n int, target *TileCollection) {
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