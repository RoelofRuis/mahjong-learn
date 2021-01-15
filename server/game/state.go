package game

func (g *Game) Transition() {
	state := g.State
	for {
		switch state {
		case StateNextRound:
			// select next player
			// deal starting hands
			// -> StateNextTurn
		case StateNextTurn:
			// IF more tiles available
			//   deal single tile
			//   -> TileReceived
			// ELSE
			//   -> StateRoundEnded
		case StateKongDeclared:
			// update player exposed tiles
			// deal single tile
			// -> TileReceived
		case StateRoundEnded:
			// update scores
			// -> StateNextRound
		}
		if g.State == state || g.State.IsObservable() {
			break
		}
		state = g.State
	}

}

func (g *Game) PerformAction(a Action) error {
	// TODO: implement
	// Return error if an action is not allowed in the current state
	return nil
}

func (g *Game) DealStartingHands() {
	g.Wall.Transfer(13, g.Players[0].Concealed)
	g.Wall.Transfer(13, g.Players[1].Concealed)
	g.Wall.Transfer(13, g.Players[2].Concealed)
	g.Wall.Transfer(13, g.Players[3].Concealed)
}

func (g *Game) DealTile(player int) {
	g.Wall.Transfer(1, g.Players[player].Concealed)
}