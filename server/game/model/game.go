package model

type Seat int

var SEATS = []Seat{Seat(0), Seat(1), Seat(2), Seat(3)}

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

func (g *Game) GetWallSize() int {
	return g.wall.Size()
}

func (g *Game) GetActiveSeat() Seat {
	return g.activeSeat
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

func (g *Game) DealTiles(n int, seat Seat) {
	player := g.players[seat]
	g.wall.TransferRandom(n, player.concealed)

	for {
		numExtra := player.ForceExposeTiles()
		if numExtra == 0 {
			break
		}

		g.wall.TransferRandom(numExtra, player.concealed)
	}
}

func (g *Game) SetNextPrevalentWind() {
	g.prevalentWind++
}

func (g *Game) ResetWall() {
	g.wall = NewMahjongSet()
}

func (g *Game) ActivateNextSeat() {
	g.activeSeat = Seat((int(g.activeSeat) + 1) % 4)
}

func (g *Game) ActivePlayerDiscards(tile Tile) {
	g.players[g.activeSeat].concealed.Remove(tile)
	g.activeDiscard = &tile
}

func (g *Game) ActivePlayerTakesDiscarded() {
	if g.activeDiscard != nil {
		g.players[g.activeSeat].discarded.Add(*g.activeDiscard)
		g.activeDiscard = nil
	}
}
