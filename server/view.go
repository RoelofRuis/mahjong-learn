package main

import (
	"fmt"
	"github.com/roelofruis/mahjong-learn/driver"
	"github.com/roelofruis/mahjong-learn/game"
)

var TileOrder = []game.Tile{
	game.Bamboo1,
	game.Bamboo2,
	game.Bamboo3,
	game.Bamboo4,
	game.Bamboo5,
	game.Bamboo6,
	game.Bamboo7,
	game.Bamboo8,
	game.Bamboo9,
	game.Circles1,
	game.Circles2,
	game.Circles3,
	game.Circles4,
	game.Circles5,
	game.Circles6,
	game.Circles7,
	game.Circles8,
	game.Circles9,
	game.Characters1,
	game.Characters2,
	game.Characters3,
	game.Characters4,
	game.Characters5,
	game.Characters6,
	game.Characters7,
	game.Characters8,
	game.Characters9,
	game.RedDragon,
	game.GreenDragon,
	game.WhiteDragon,
	game.EastWind,
	game.SouthWind,
	game.WestWind,
	game.NorthWind,
	game.FlowerPlumb,
	game.FlowerOrchid,
	game.FlowerChrysanthemum,
	game.FlowerBamboo,
	game.SeasonSpring,
	game.SeasonSummer,
	game.SeasonAutumn,
	game.SeasonWinter,
}

var TileNames = map[game.Tile]string{
	game.Bamboo1:             "Bamboo 1",
	game.Bamboo2:             "Bamboo 2",
	game.Bamboo3:             "Bamboo 3",
	game.Bamboo4:             "Bamboo 4",
	game.Bamboo5:             "Bamboo 5",
	game.Bamboo6:             "Bamboo 6",
	game.Bamboo7:             "Bamboo 7",
	game.Bamboo8:             "Bamboo 8",
	game.Bamboo9:             "Bamboo 9",
	game.Circles1:            "Circles 1",
	game.Circles2:            "Circles 2",
	game.Circles3:            "Circles 3",
	game.Circles4:            "Circles 4",
	game.Circles5:            "Circles 5",
	game.Circles6:            "Circles 6",
	game.Circles7:            "Circles 7",
	game.Circles8:            "Circles 8",
	game.Circles9:            "Circles 9",
	game.Characters1:         "Characters 1",
	game.Characters2:         "Characters 2",
	game.Characters3:         "Characters 3",
	game.Characters4:         "Characters 4",
	game.Characters5:         "Characters 5",
	game.Characters6:         "Characters 6",
	game.Characters7:         "Characters 7",
	game.Characters8:         "Characters 8",
	game.Characters9:         "Characters 9",
	game.RedDragon:           "Red Dragon",
	game.GreenDragon:         "Green Dragon",
	game.WhiteDragon:         "White Dragon",
	game.EastWind:            "East Wind",
	game.SouthWind:           "South Wind",
	game.WestWind:            "West Wind",
	game.NorthWind:           "North Wind",
	game.FlowerPlumb:         "Plumb (flower)",
	game.FlowerOrchid:        "Orchid (flower)",
	game.FlowerChrysanthemum: "Chrysanthemum (flower)",
	game.FlowerBamboo:        "Bamboo (flower)",
	game.SeasonSpring:        "Spring (season)",
	game.SeasonSummer:        "Summer (season)",
	game.SeasonAutumn:        "Autumn (season)",
	game.SeasonWinter:        "Winter (season)",
}

var WindNames = map[game.Wind]string{
	game.East:  "East",
	game.South: "South",
	game.West:  "West",
	game.North: "North",
}

type PlayerView struct {
	Actions   map[int]string `json:"actions"`
	Wind      string         `json:"wind"`
	Received  string         `json:"received"`
	Concealed []string       `json:"concealed"`
	Exposed   []string       `json:"exposed"`
	Discarded []string       `json:"discarded"`
}

type HumanView struct {
	Id            uint64             `json:"id"`
	HasEnded      bool               `json:"has_ended"`
	StateName     string             `json:"state_name"`
	PrevalentWind string             `json:"prevalent_wind"`
	ActivePlayers []int              `json:"active_players"`
	ActiveDiscard string             `json:"active_discard"`
	Players       map[int]PlayerView `json:"players"`
	Wall          []string           `json:"wall"`
}

func View(game *game.MahjongGame) *HumanView {
	s := game.Driver.GetState()
	a := s.Actions()
	g := *game.Table

	var activePlayers []int
	playerViews := make(map[int]PlayerView, 4)
	for _, seat := range []int{0, 1, 2, 3} {
		actions, has := a[driver.Seat(seat)]
		if !has {
			actions = make([]driver.Action, 0)
		} else {
			activePlayers = append(activePlayers, seat+1)
		}
		playerViews[seat+1] = DescribePlayer(g, actions, driver.Seat(seat))
	}

	return &HumanView{
		Id:            game.Id,
		HasEnded:      s.Transition == nil,
		StateName:     s.Name,
		PrevalentWind: WindNames[g.GetPrevalentWind()],
		ActivePlayers: activePlayers,
		ActiveDiscard: DescribeTilePointer(g.GetActiveDiscard()),
		Players:       playerViews,
		Wall:          Describe(g.GetWall()),
	}
}

func DescribeTilePointer(t *game.Tile) string {
	if t == nil {
		return "none"
	}
	return TileNames[*t]
}

func DescribeCombinations(combinations []game.Combination) []string {
	descriptions := make([]string, len(combinations))
	for i, combi := range combinations {
		switch c := combi.(type) {
		case game.BonusTile:
			descriptions[i] = fmt.Sprintf("Bonus tile %s", TileNames[c.Tile])

		case game.Chow:
			descriptions[i] = fmt.Sprintf("Chow %s", TileNames[c.FirstTile])

		case game.Pung:
			descriptions[i] = fmt.Sprintf("Pung %s", TileNames[c.Tile])

		case game.Kong:
			descriptions[i] = fmt.Sprintf("Kong %s", TileNames[c.Tile])

		default:
			// This should not happen..!
			descriptions[i] = "unknown combination"
		}
	}
	return descriptions
}

func Describe(t *game.TileCollection) []string {
	descriptions := make([]string, 0)
	for _, tile := range TileOrder {
		count := t.NumOf(tile)
		if count == 0 {
			continue
		}
		text := fmt.Sprintf("%d× %s", count, TileNames[tile])
		descriptions = append(descriptions, text)
	}
	return descriptions
}

func DescribePlayer(g game.Table, actions []driver.Action, seat driver.Seat) PlayerView {
	actionMap := make(map[int]string)
	for i, a := range actions {
		actionMap[i] = DescribeAction(a)
	}

	p := g.GetPlayerAtSeat(seat)

	return PlayerView{
		Actions:   actionMap,
		Wind:      WindNames[p.GetSeatWind()],
		Received:  DescribeTilePointer(p.GetReceivedTile()),
		Concealed: Describe(p.GetConcealedTiles()),
		Exposed:   DescribeCombinations(p.GetExposedCombinations()),
		Discarded: Describe(p.GetDiscardedTiles()),
	}
}

func DescribeAction(action driver.Action) string {
	switch a := action.(type) {
	case game.Discard:
		return fmt.Sprintf("Discard a %s", TileNames[a.Tile])
	case game.DeclareConcealedKong:
		return fmt.Sprintf("Declare a concealed Kong of %s", TileNames[a.Tile])
	case game.ExposedPungToKong:
		return fmt.Sprintf("Add to exposed pung")
	case game.DoNothing:
		return "Do nothing"
	case game.DeclareChow:
		return fmt.Sprintf("Declare chow up from %s", TileNames[a.Tile])
	case game.DeclarePung:
		return "Declare a pung"
	case game.DeclareKong:
		return "Declare a kong"
	case game.DeclareMahjong:
		return "Declare mahjong"

	default:
		panic(fmt.Errorf("unknown action %+v", a))
	}
}
