package model

type Seat int

type Wind int

const (
	East  Wind = 0
	South Wind = 1
	West  Wind = 2
	North Wind = 3
)

type Game struct {
	prevalentWind Wind
	wall          *TileCollection
	activeDiscard *Tile
	players       map[Seat]*Player
	activeSeat    Seat
}

func NewGame() *Game {
	players := make(map[Seat]*Player, 4)

	wall := NewMahjongSet()
	players[0] = NewPlayer(East)
	players[1] = NewPlayer(South)
	players[2] = NewPlayer(West)
	players[3] = NewPlayer(North)

	return &Game{
		prevalentWind: East,
		wall:          wall,
		activeDiscard: nil,
		players:       players,
		activeSeat:    0,
	}
}

// Getters

func (g *Game) GetWallSize() int {
	return g.wall.Size()
}

func (g *Game) GetActiveSeat() Seat {
	return g.activeSeat
}

func (g *Game) GetReactingPlayers() map[Seat]*Player {
	reactingPlayers := make(map[Seat]*Player, 3)
	for s, p := range g.players {
		if s != g.activeSeat {
			reactingPlayers[s] = p
		}
	}
	return reactingPlayers
}

func (g *Game) GetActivePlayer() *Player {
	return g.GetPlayerAtSeat(g.activeSeat)
}

func (g *Game) GetPlayerAtSeat(seat Seat) *Player {
	return g.players[seat]
}

func (g *Game) GetPrevalentWind() Wind {
	return g.prevalentWind
}

func (g *Game) GetActiveDiscard() *Tile {
	return g.activeDiscard
}

func (g *Game) GetWall() *TileCollection {
	return g.wall
}

// State Updates

func (g *Game) DealToActivePlayer() {
	activePlayer := g.GetActivePlayer()

	for {
		wallTile := g.wall.RemoveRandom()

		if !IsBonusTile(wallTile) {
			activePlayer.received = &wallTile
			break
		}

		activePlayer.exposed.Add(BonusTile{wallTile})
	}
}

func (g *Game) DealConcealed(n int, s Seat) {
	activePlayer := g.players[s]

	for i := n; i > 0; i-- {
		for {
			wallTile := g.wall.RemoveRandom()

			if !IsBonusTile(wallTile) {
				activePlayer.concealed.Add(wallTile)
				break
			}

			activePlayer.exposed.Add(BonusTile{wallTile})
		}
	}
}

func (g *Game) SetNextPrevalentWind() {
	g.prevalentWind++
}

func (g *Game) ResetWall() {
	g.wall = NewMahjongSet()
}

func (g *Game) PrepareNextRound() {
	for s, p := range g.players {
		p.received = nil
		p.discarded.Empty()
		p.concealed.Empty()
		p.exposed.Empty()
		p.seatWind = (p.seatWind + 5) % 4
		g.DealConcealed(13, s)
	}
}

func (g *Game) ActivateNextSeat() {
	g.activeSeat = Seat((int(g.activeSeat) + 1) % 4)
}

func (g *Game) ActivePlayerDiscards(tile Tile) {
	activePlayer := g.GetActivePlayer()

	// transfer received tile to hand
	activePlayer.concealed.Add(*activePlayer.received)
	activePlayer.received = nil

	// remove selected tile from hand
	activePlayer.concealed.Remove(tile)

	// set active discard to selected tile
	g.activeDiscard = &tile
}

func (g *Game) ActivePlayerTakesDiscarded() {
	if g.activeDiscard != nil {
		g.GetActivePlayer().discarded.Add(*g.activeDiscard)
		g.activeDiscard = nil
	}
}
