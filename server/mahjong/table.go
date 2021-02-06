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

func NewTable() *Table {
	players := make(map[int]*Player, 4)

	wall := NewMahjongSet()
	players[0] = NewPlayer(East)
	players[1] = NewPlayer(South)
	players[2] = NewPlayer(West)
	players[3] = NewPlayer(North)

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

func (t *Table) DealToActivePlayer() {
	activePlayer := t.GetActivePlayer()

	for {
		wallTile := t.wall.RemoveRandom()

		if !IsBonusTile(wallTile) {
			activePlayer.received = &wallTile
			break
		}

		activePlayer.exposed.Add(BonusTile{wallTile})
	}
}

func (t *Table) DealConcealed(n int, player int) {
	activePlayer := t.players[player]

	for i := n; i > 0; i-- {
		for {
			wallTile := t.wall.RemoveRandom()

			if !IsBonusTile(wallTile) {
				activePlayer.concealed.Add(wallTile)
				break
			}

			activePlayer.exposed.Add(BonusTile{wallTile})
		}
	}
}

func (t *Table) SetNextPrevalentWind() {
	t.prevalentWind++
}

func (t *Table) ResetWall() {
	t.wall = NewMahjongSet()
}

func (t *Table) PrepareNextRound() {
	for s, p := range t.players {
		p.received = nil
		p.discarded.Empty()
		p.concealed.Empty()
		p.exposed.Empty()
		p.wind = (p.wind + 5) % 4
		t.DealConcealed(13, s)
	}
}

func (t *Table) MakePlayerActive(player int) {
	t.activePlayer = player
}

func (t *Table) ActivePlayerDeclaresConcealedKong(tile Tile) {
	activePlayer := t.GetActivePlayer()

	if activePlayer.received != nil {
		// always add to hand first, the player may be declaring a different concealed kong than with the tile just dealt.
		activePlayer.concealed.Add(*activePlayer.received)
		activePlayer.received = nil
	}

	activePlayer.concealed.RemoveAll(tile)
	activePlayer.exposed.Add(Kong{
		Tile:      tile,
		Concealed: false,
	})
}

func (t *Table) ActivePlayerAddsToExposedPung() {
	activePlayer := t.GetActivePlayer()

	activePlayer.exposed.Replace(
		Pung{Tile: *activePlayer.received},
		Kong{Tile: *activePlayer.received, Concealed: false},
	)

	activePlayer.received = nil
}

func (t *Table) ActivePlayerDiscards(tile Tile) {
	activePlayer := t.GetActivePlayer()

	if activePlayer.received != nil {
		activePlayer.concealed.Add(*activePlayer.received)
		activePlayer.received = nil
	}

	activePlayer.concealed.Remove(tile)

	t.activeDiscard = &tile
}

func (t *Table) ActivePlayerTakesDiscarded() {
	if t.activeDiscard != nil {
		t.GetActivePlayer().discarded.Add(*t.activeDiscard)
		t.activeDiscard = nil
	}
}

func (t *Table) ActivePlayerTakesChow(tile Tile) {
	if t.activeDiscard != nil {
		activePlayer := t.GetActivePlayer()
		activePlayer.concealed.Add(*t.activeDiscard)
		activePlayer.concealed.Remove(tile)
		activePlayer.concealed.Remove(tile + 1)
		activePlayer.concealed.Remove(tile + 2)
		activePlayer.exposed.Add(Chow{FirstTile: tile})
		t.activeDiscard = nil
	}
}

func (t *Table) ActivePlayerTakesPung() {
	if t.activeDiscard != nil {
		activePlayer := t.GetActivePlayer()
		activePlayer.concealed.Remove(*t.activeDiscard)
		activePlayer.concealed.Remove(*t.activeDiscard)
		activePlayer.exposed.Add(Pung{Tile: *t.activeDiscard})
		t.activeDiscard = nil
	}
}

func (t *Table) ActivePlayerTakesKong() {
	if t.activeDiscard != nil {
		activePlayer := t.GetActivePlayer()
		activePlayer.concealed.RemoveAll(*t.activeDiscard)
		activePlayer.exposed.Add(Kong{Tile: *t.activeDiscard, Concealed: false})
		t.activeDiscard = nil
	}
}
