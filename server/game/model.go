package game

type StateMachine struct {
	State *State
	Game  *Game
}

type Action func(*Game) *State

type State struct {
	// Name just to display human readable information.
	Name string

	// Transfer to next state via action, or nil if player input is required.
	TransferAction Action
	// Show required player actions. This requires TransferAction to be nil.
	RequiredActions func(*Game) map[Seat][]Action
}

type Game struct {
	// FIXME: add lock to Game so we can modify data freely in a request and block simultaneous requests
	Id uint64

	PrevalentWind Wind
	Wall          *TileCollection
	Players       map[Seat]Player
	ActiveSeat    Seat
}

type Player struct {
	Score int

	SeatWind  Wind
	Concealed *TileCollection
	Exposed   []*TileCollection
	Discarded *TileCollection
}

type TileCollection struct {
	Tiles map[Tile]int
}

type Seat int

type Wind int

const (
	East  Wind = 0
	South Wind = 1
	West  Wind = 2
	North Wind = 3
)

type Tile int

const (
	Bamboo1             Tile = 1
	Bamboo2             Tile = 2
	Bamboo3             Tile = 3
	Bamboo4             Tile = 4
	Bamboo5             Tile = 5
	Bamboo6             Tile = 6
	Bamboo7             Tile = 7
	Bamboo8             Tile = 8
	Bamboo9             Tile = 9
	Circles1            Tile = 11
	Circles2            Tile = 12
	Circles3            Tile = 13
	Circles4            Tile = 14
	Circles5            Tile = 15
	Circles6            Tile = 16
	Circles7            Tile = 17
	Circles8            Tile = 18
	Circles9            Tile = 19
	Characters1         Tile = 21
	Characters2         Tile = 22
	Characters3         Tile = 23
	Characters4         Tile = 24
	Characters5         Tile = 25
	Characters6         Tile = 26
	Characters7         Tile = 27
	Characters8         Tile = 28
	Characters9         Tile = 29
	RedDragon           Tile = 30
	GreenDragon         Tile = 31
	WhiteDragon         Tile = 32
	EastWind            Tile = 40
	SouthWind           Tile = 41
	WestWind            Tile = 42
	NorthWind           Tile = 43
	FlowerPlumb         Tile = 50
	FlowerOrchid        Tile = 51
	FlowerChrysanthemum Tile = 52
	FlowerBamboo        Tile = 53
	SeasonSpring        Tile = 60
	SeasonSummer        Tile = 61
	SeasonAutumn        Tile = 62
	SeasonWinter        Tile = 63
)
