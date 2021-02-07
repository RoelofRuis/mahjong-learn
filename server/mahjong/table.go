package mahjong

type Wind int

const (
	East  Wind = 0
	South Wind = 1
	West  Wind = 2
	North Wind = 3
)

type Table struct {
	prevalentWind Wind
	wall          *TileCollection
	activeDiscard *Tile
	players       map[int]*Player
	activePlayer  int
}

func newTable() *Table {
	players := make(map[int]*Player, 4)

	wall := newMahjongSet()
	players[0] = newPlayer(East)
	players[1] = newPlayer(South)
	players[2] = newPlayer(West)
	players[3] = newPlayer(North)

	return &Table{
		prevalentWind: East,
		wall:          wall,
		activeDiscard: nil,
		players:       players,
		activePlayer:  0,
	}
}

// Getters

func (t *Table) GetWallSize() int {
	return t.wall.Size()
}

func (t *Table) GetActivePlayerIndex() int {
	return t.activePlayer
}

func (t *Table) GetReactingPlayers() map[int]*Player {
	reactingPlayers := make(map[int]*Player, 3)
	for s, p := range t.players {
		if s != t.activePlayer {
			reactingPlayers[s] = p
		}
	}
	return reactingPlayers
}

func (t *Table) GetActivePlayer() *Player {
	return t.GetPlayerByIndex(t.activePlayer)
}

func (t *Table) GetPlayerByIndex(player int) *Player {
	return t.players[player]
}

func (t *Table) GetPrevalentWind() Wind {
	return t.prevalentWind
}

func (t *Table) GetActiveDiscard() *Tile {
	return t.activeDiscard
}

func (t *Table) GetWall() *TileCollection {
	return t.wall
}

// State Updates

func (t *Table) dealToActivePlayer() {
	activePlayer := t.GetActivePlayer()

	for {
		wallTile := t.wall.removeRandom()

		if !wallTile.IsBonusTile() {
			activePlayer.received = &wallTile
			break
		}

		activePlayer.exposed.add(BonusTile{wallTile})
	}
}

func (t *Table) dealConcealed(n int, player int) {
	activePlayer := t.players[player]

	for i := n; i > 0; i-- {
		for {
			wallTile := t.wall.removeRandom()

			if !wallTile.IsBonusTile() {
				activePlayer.concealed.add(wallTile)
				break
			}

			activePlayer.exposed.add(BonusTile{wallTile})
		}
	}
}

func (t *Table) setNextPrevalentWind() {
	t.prevalentWind++
}

func (t *Table) resetWall() {
	t.wall = newMahjongSet()
}

func (t *Table) prepareNextRound() {
	for s, p := range t.players {
		p.received = nil
		p.discarded.empty()
		p.concealed.empty()
		p.exposed.empty()
		p.wind = (p.wind + 5) % 4
		t.dealConcealed(13, s)
	}
}

func (t *Table) makePlayerActive(player int) {
	t.activePlayer = player
}

func (t *Table) activePlayerDeclaresConcealedKong(tile Tile) {
	activePlayer := t.GetActivePlayer()

	if activePlayer.received != nil {
		// always add to hand first, the player may be declaring a different concealed kong than with the tile just dealt.
		activePlayer.concealed.add(*activePlayer.received)
		activePlayer.received = nil
	}

	activePlayer.concealed.removeAll(tile)
	activePlayer.exposed.add(Kong{
		Tile:      tile,
		Concealed: false,
	})
}

func (t *Table) activePlayerAddsToExposedPung() {
	activePlayer := t.GetActivePlayer()

	activePlayer.exposed.replace(
		Pung{Tile: *activePlayer.received},
		Kong{Tile: *activePlayer.received, Concealed: false},
	)

	activePlayer.received = nil
}

func (t *Table) activePlayerDiscards(tile Tile) {
	activePlayer := t.GetActivePlayer()

	if activePlayer.received != nil {
		activePlayer.concealed.add(*activePlayer.received)
		activePlayer.received = nil
	}

	activePlayer.concealed.remove(tile)

	t.activeDiscard = &tile
}

func (t *Table) activePlayerTakesDiscarded() {
	if t.activeDiscard != nil {
		t.GetActivePlayer().discarded.add(*t.activeDiscard)
		t.activeDiscard = nil
	}
}

func (t *Table) activePlayerTakesChow(tile Tile) {
	if t.activeDiscard != nil {
		activePlayer := t.GetActivePlayer()
		activePlayer.concealed.add(*t.activeDiscard)
		activePlayer.concealed.remove(tile)
		activePlayer.concealed.remove(tile + 1)
		activePlayer.concealed.remove(tile + 2)
		activePlayer.exposed.add(Chow{FirstTile: tile})
		t.activeDiscard = nil
	}
}

func (t *Table) activePlayerTakesPung() {
	if t.activeDiscard != nil {
		activePlayer := t.GetActivePlayer()
		activePlayer.concealed.remove(*t.activeDiscard)
		activePlayer.concealed.remove(*t.activeDiscard)
		activePlayer.exposed.add(Pung{Tile: *t.activeDiscard})
		t.activeDiscard = nil
	}
}

func (t *Table) activePlayerTakesKong() {
	if t.activeDiscard != nil {
		activePlayer := t.GetActivePlayer()
		activePlayer.concealed.removeAll(*t.activeDiscard)
		activePlayer.exposed.add(Kong{Tile: *t.activeDiscard, Concealed: false})
		t.activeDiscard = nil
	}
}
