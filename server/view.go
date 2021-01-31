package main

import (
	"errors"
	"fmt"
	"github.com/roelofruis/mahjong-learn/game"
	"github.com/roelofruis/mahjong-learn/game/model"
)

var TileOrder = []model.Tile{
	model.Bamboo1,
	model.Bamboo2,
	model.Bamboo3,
	model.Bamboo4,
	model.Bamboo5,
	model.Bamboo6,
	model.Bamboo7,
	model.Bamboo8,
	model.Bamboo9,
	model.Circles1,
	model.Circles2,
	model.Circles3,
	model.Circles4,
	model.Circles5,
	model.Circles6,
	model.Circles7,
	model.Circles8,
	model.Circles9,
	model.Characters1,
	model.Characters2,
	model.Characters3,
	model.Characters4,
	model.Characters5,
	model.Characters6,
	model.Characters7,
	model.Characters8,
	model.Characters9,
	model.RedDragon,
	model.GreenDragon,
	model.WhiteDragon,
	model.EastWind,
	model.SouthWind,
	model.WestWind,
	model.NorthWind,
	model.FlowerPlumb,
	model.FlowerOrchid,
	model.FlowerChrysanthemum,
	model.FlowerBamboo,
	model.SeasonSpring,
	model.SeasonSummer,
	model.SeasonAutumn,
	model.SeasonWinter,
}

var TileNames = map[model.Tile]string{
	model.Bamboo1:             "Bamboo 1",
	model.Bamboo2:             "Bamboo 2",
	model.Bamboo3:             "Bamboo 3",
	model.Bamboo4:             "Bamboo 4",
	model.Bamboo5:             "Bamboo 5",
	model.Bamboo6:             "Bamboo 6",
	model.Bamboo7:             "Bamboo 7",
	model.Bamboo8:             "Bamboo 8",
	model.Bamboo9:             "Bamboo 9",
	model.Circles1:            "Circles 1",
	model.Circles2:            "Circles 2",
	model.Circles3:            "Circles 3",
	model.Circles4:            "Circles 4",
	model.Circles5:            "Circles 5",
	model.Circles6:            "Circles 6",
	model.Circles7:            "Circles 7",
	model.Circles8:            "Circles 8",
	model.Circles9:            "Circles 9",
	model.Characters1:         "Characters 1",
	model.Characters2:         "Characters 2",
	model.Characters3:         "Characters 3",
	model.Characters4:         "Characters 4",
	model.Characters5:         "Characters 5",
	model.Characters6:         "Characters 6",
	model.Characters7:         "Characters 7",
	model.Characters8:         "Characters 8",
	model.Characters9:         "Characters 9",
	model.RedDragon:           "Red Dragon",
	model.GreenDragon:         "Green Dragon",
	model.WhiteDragon:         "White Dragon",
	model.EastWind:            "East Wind",
	model.SouthWind:           "South Wind",
	model.WestWind:            "West Wind",
	model.NorthWind:           "North Wind",
	model.FlowerPlumb:         "Plumb (flower)",
	model.FlowerOrchid:        "Orchid (flower)",
	model.FlowerChrysanthemum: "Chrysanthemum (flower)",
	model.FlowerBamboo:        "Bamboo (flower)",
	model.SeasonSpring:        "Spring (season)",
	model.SeasonSummer:        "Summer (season)",
	model.SeasonAutumn:        "Autumn (season)",
	model.SeasonWinter:        "Winter (season)",
}

var WindNames = map[model.Wind]string{
	model.East:  "East",
	model.South: "South",
	model.West:  "West",
	model.North: "North",
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

func View(stateMachine *game.StateMachine) *HumanView {
	g, s, a := stateMachine.View()

	var activePlayers []int
	playerViews := make(map[int]PlayerView, 4)
	for _, seat := range []int{0, 1, 2, 3} {
		actions, has := a[model.Seat(seat)]
		if !has {
			actions = make([]model.Action, 0)
		} else {
			activePlayers = append(activePlayers, seat+1)
		}
		playerViews[seat+1] = DescribePlayer(g, actions, model.Seat(seat))
	}

	return &HumanView{
		Id:            stateMachine.Id(),
		HasEnded:      s.Transition == nil,
		StateName:     s.Name,
		PrevalentWind: WindNames[g.GetPrevalentWind()],
		ActivePlayers: activePlayers,
		ActiveDiscard: DescribeTilePointer(g.GetActiveDiscard()),
		Players:       playerViews,
		Wall:          Describe(g.GetWall()),
	}
}

func DescribeTilePointer(t *model.Tile) string {
	if t == nil {
		return "none"
	}
	return TileNames[*t]
}

func DescribeCombinations(combinations []model.Combination) []string {
	descriptions := make([]string, len(combinations))
	for i, combi := range combinations {
		switch c := combi.(type) {
		case model.BonusTile:
			descriptions[i] = fmt.Sprintf("Bonus tile %s", TileNames[c.Tile])

		case model.Chow:
			descriptions[i] = fmt.Sprintf("Chow %s", TileNames[c.FirstTile])

		case model.Pung:
			descriptions[i] = fmt.Sprintf("Pung %s", TileNames[c.Tile])

		case model.Kong:
			descriptions[i] = fmt.Sprintf("Kong %s", TileNames[c.Tile])

		default:
			// This should not happen..!
			descriptions[i] = "unknown combination"
		}
	}
	return descriptions
}

func Describe(t *model.TileCollection) []string {
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

func DescribePlayer(g model.Game, actions []model.Action, seat model.Seat) PlayerView {
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

func DescribeAction(action model.Action) string {
	switch a := action.(type) {
	case model.Discard:
		return fmt.Sprintf("Discard a %s", TileNames[a.Tile])
	case model.DeclareConcealedKong:
		return fmt.Sprintf("Declare a concealed Kong of %s", TileNames[a.Tile])
	case model.ExposedPungToKong:
		return fmt.Sprintf("Add to exposed pung of %s", TileNames[a.Tile])
	case model.DoNothing:
		return "Do nothing"
	case model.DeclareChow:
		return fmt.Sprintf("Declare chow up from %s", TileNames[a.Tile])
	case model.DeclarePung:
		return "Declare a pung"
	case model.DeclareKong:
		return "Declare a kong"
	case model.DeclareMahjong:
		return "Declare mahjong"

	default:
		panic(errors.New("unknown action")) // This should never happen..!
	}
}
