package mahjong

import "github.com/roelofruis/mahjong-learn/state"

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
	players       map[state.Seat]*Player
	activeSeat    state.Seat
}

func NewTable() *Table {
	players := make(map[state.Seat]*Player, 4)

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
		activeSeat:    0,
	}
}

// Getters

func (t *Table) GetWallSize() int {
	return t.wall.Size()
}

func (t *Table) GetActiveSeat() state.Seat {
	return t.activeSeat
}

func (t *Table) GetReactingPlayers() map[state.Seat]*Player {
	reactingPlayers := make(map[state.Seat]*Player, 3)
	for s, p := range t.players {
		if s != t.activeSeat {
			reactingPlayers[s] = p
		}
	}
	return reactingPlayers
}

func (t *Table) GetActivePlayer() *Player {
	return t.GetPlayerAtSeat(t.activeSeat)
}

func (t *Table) GetPlayerAtSeat(seat state.Seat) *Player {
	return t.players[seat]
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

func (t *Table) DealConcealed(n int, s state.Seat) {
	activePlayer := t.players[s]

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
		p.seatWind = (p.seatWind + 5) % 4
		t.DealConcealed(13, s)
	}
}

func (t *Table) ActivateSeat(seat state.Seat) {
	t.activeSeat = seat
}

func (t *Table) ActivePlayerDeclaresConcealedKong(tile Tile) {
	activePlayer := t.GetActivePlayer()

	activePlayer.concealed.RemoveAll(tile)
	activePlayer.exposed.Add(Kong{
		Tile:      tile,
		Concealed: false,
	})

	activePlayer.received = nil
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
