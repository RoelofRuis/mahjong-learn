package main

import (
	"fmt"
	"github.com/roelofruis/mahjong-learn/driver"
	"github.com/roelofruis/mahjong-learn/mahjong"
)

var TileOrder = []mahjong.Tile{
	mahjong.Bamboo1,
	mahjong.Bamboo2,
	mahjong.Bamboo3,
	mahjong.Bamboo4,
	mahjong.Bamboo5,
	mahjong.Bamboo6,
	mahjong.Bamboo7,
	mahjong.Bamboo8,
	mahjong.Bamboo9,
	mahjong.Circles1,
	mahjong.Circles2,
	mahjong.Circles3,
	mahjong.Circles4,
	mahjong.Circles5,
	mahjong.Circles6,
	mahjong.Circles7,
	mahjong.Circles8,
	mahjong.Circles9,
	mahjong.Characters1,
	mahjong.Characters2,
	mahjong.Characters3,
	mahjong.Characters4,
	mahjong.Characters5,
	mahjong.Characters6,
	mahjong.Characters7,
	mahjong.Characters8,
	mahjong.Characters9,
	mahjong.RedDragon,
	mahjong.GreenDragon,
	mahjong.WhiteDragon,
	mahjong.EastWind,
	mahjong.SouthWind,
	mahjong.WestWind,
	mahjong.NorthWind,
	mahjong.FlowerPlumb,
	mahjong.FlowerOrchid,
	mahjong.FlowerChrysanthemum,
	mahjong.FlowerBamboo,
	mahjong.SeasonSpring,
	mahjong.SeasonSummer,
	mahjong.SeasonAutumn,
	mahjong.SeasonWinter,
}

var TileNames = map[mahjong.Tile]string{
	mahjong.Bamboo1:             "Bamboo 1",
	mahjong.Bamboo2:             "Bamboo 2",
	mahjong.Bamboo3:             "Bamboo 3",
	mahjong.Bamboo4:             "Bamboo 4",
	mahjong.Bamboo5:             "Bamboo 5",
	mahjong.Bamboo6:             "Bamboo 6",
	mahjong.Bamboo7:             "Bamboo 7",
	mahjong.Bamboo8:             "Bamboo 8",
	mahjong.Bamboo9:             "Bamboo 9",
	mahjong.Circles1:            "Circles 1",
	mahjong.Circles2:            "Circles 2",
	mahjong.Circles3:            "Circles 3",
	mahjong.Circles4:            "Circles 4",
	mahjong.Circles5:            "Circles 5",
	mahjong.Circles6:            "Circles 6",
	mahjong.Circles7:            "Circles 7",
	mahjong.Circles8:            "Circles 8",
	mahjong.Circles9:            "Circles 9",
	mahjong.Characters1:         "Characters 1",
	mahjong.Characters2:         "Characters 2",
	mahjong.Characters3:         "Characters 3",
	mahjong.Characters4:         "Characters 4",
	mahjong.Characters5:         "Characters 5",
	mahjong.Characters6:         "Characters 6",
	mahjong.Characters7:         "Characters 7",
	mahjong.Characters8:         "Characters 8",
	mahjong.Characters9:         "Characters 9",
	mahjong.RedDragon:           "Red Dragon",
	mahjong.GreenDragon:         "Green Dragon",
	mahjong.WhiteDragon:         "White Dragon",
	mahjong.EastWind:            "East Wind",
	mahjong.SouthWind:           "South Wind",
	mahjong.WestWind:            "West Wind",
	mahjong.NorthWind:           "North Wind",
	mahjong.FlowerPlumb:         "Plumb (flower)",
	mahjong.FlowerOrchid:        "Orchid (flower)",
	mahjong.FlowerChrysanthemum: "Chrysanthemum (flower)",
	mahjong.FlowerBamboo:        "Bamboo (flower)",
	mahjong.SeasonSpring:        "Spring (season)",
	mahjong.SeasonSummer:        "Summer (season)",
	mahjong.SeasonAutumn:        "Autumn (season)",
	mahjong.SeasonWinter:        "Winter (season)",
}

var WindNames = map[mahjong.Wind]string{
	mahjong.East:  "East",
	mahjong.South: "South",
	mahjong.West:  "West",
	mahjong.North: "North",
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

func View(game *mahjong.Game) *HumanView {
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

func DescribeTilePointer(t *mahjong.Tile) string {
	if t == nil {
		return "none"
	}
	return TileNames[*t]
}

func DescribeCombinations(combinations []mahjong.Combination) []string {
	descriptions := make([]string, len(combinations))
	for i, combi := range combinations {
		switch c := combi.(type) {
		case mahjong.BonusTile:
			descriptions[i] = fmt.Sprintf("Bonus tile %s", TileNames[c.Tile])

		case mahjong.Chow:
			descriptions[i] = fmt.Sprintf("Chow %s", TileNames[c.FirstTile])

		case mahjong.Pung:
			descriptions[i] = fmt.Sprintf("Pung %s", TileNames[c.Tile])

		case mahjong.Kong:
			descriptions[i] = fmt.Sprintf("Kong %s", TileNames[c.Tile])

		default:
			// This should not happen..!
			descriptions[i] = "unknown combination"
		}
	}
	return descriptions
}

func Describe(t *mahjong.TileCollection) []string {
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

func DescribePlayer(g mahjong.Table, actions []driver.Action, seat driver.Seat) PlayerView {
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
	case mahjong.Discard:
		return fmt.Sprintf("Discard a %s", TileNames[a.Tile])
	case mahjong.DeclareConcealedKong:
		return fmt.Sprintf("Declare a concealed Kong of %s", TileNames[a.Tile])
	case mahjong.ExposedPungToKong:
		return fmt.Sprintf("Add to exposed pung")
	case mahjong.DoNothing:
		return "Do nothing"
	case mahjong.DeclareChow:
		return fmt.Sprintf("Declare chow up from %s", TileNames[a.Tile])
	case mahjong.DeclarePung:
		return "Declare a pung"
	case mahjong.DeclareKong:
		return "Declare a kong"
	case mahjong.DeclareMahjong:
		return "Declare mahjong"

	default:
		panic(fmt.Errorf("unknown action %+v", a))
	}
}
