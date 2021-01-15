package game

func (g *Game) DealStartingHands() {
	g.Wall.Transfer(13, g.Players[0].Concealed)
	g.Wall.Transfer(13, g.Players[1].Concealed)
	g.Wall.Transfer(13, g.Players[2].Concealed)
	g.Wall.Transfer(13, g.Players[3].Concealed)
}
