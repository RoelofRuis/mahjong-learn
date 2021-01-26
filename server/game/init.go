package game

func NewGameStateMachine(id uint64) *StateMachine {
	players := make(map[Seat]*Player, 4)

	players[0] = NewPlayer(East)
	players[1] = NewPlayer(South)
	players[2] = NewPlayer(West)
	players[3] = NewPlayer(North)
	tileSet := NewMahjongSet()

	g := &Game{
		PrevalentWind: East,
		Wall:          tileSet,
		Players:       players,
		ActiveSeat:    0,
	}

	return &StateMachine{
		id:    id,
		state: stateNewGame,
		game:  g,
	}
}

func NewPlayer(seatWind Wind) *Player {
	return &Player{
		Score:     0,
		SeatWind:  seatWind,
		Concealed: NewEmptyTileCollection(),
		Exposed:   []Combination{},
		Discarded: NewEmptyTileCollection(),
	}
}

func NewEmptyTileCollection() *TileCollection {
	return &TileCollection{Tiles: make(map[Tile]int)}
}

func NewMahjongSet() *TileCollection {
	return &TileCollection{Tiles: map[Tile]int{
		Bamboo1: 4, Bamboo2: 4, Bamboo3: 4, Bamboo4: 4, Bamboo5: 4, Bamboo6: 4, Bamboo7: 4, Bamboo8: 4, Bamboo9: 4,
		Circles1: 4, Circles2: 4, Circles3: 4, Circles4: 4, Circles5: 4, Circles6: 4, Circles7: 4, Circles8: 4, Circles9: 4,
		Characters1: 4, Characters2: 4, Characters3: 4, Characters4: 4, Characters5: 4, Characters6: 4, Characters7: 4, Characters8: 4, Characters9: 4,

		RedDragon: 4, GreenDragon: 4, WhiteDragon: 4,
		EastWind: 4, SouthWind: 4, WestWind: 4, NorthWind: 4,

		FlowerPlumb: 1, FlowerOrchid: 1, FlowerChrysanthemum: 1, FlowerBamboo: 1,
		SeasonSpring: 1, SeasonSummer: 1, SeasonAutumn: 1, SeasonWinter: 1,
	}}
}
