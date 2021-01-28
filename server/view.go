package main

import (
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

var WindNames = map[model.Wind]string{
	model.East:  "East",
	model.South: "South",
	model.West:  "West",
	model.North: "North",
}

type PlayerView struct {
	Actions   map[int]string `json:"actions"`
	Wind      string         `json:"wind"`
	Concealed []string       `json:"hand"`
	Exposed   []string       `json:"exposed"`
	Discarded []string       `json:"discarded"`
}

type HumanView struct {
	Id            uint64     `json:"id"`
	HasEnded      bool       `json:"has_ended"`
	StateName     string     `json:"state_name"`
	PrevalentWind string     `json:"prevalent_wind"`
	ActivePlayer  int        `json:"active_player"`
	ActiveDiscard string     `json:"active_discard"`
	Player1       PlayerView `json:"player_1"`
	Player2       PlayerView `json:"player_2"`
	Player3       PlayerView `json:"player_3"`
	Player4       PlayerView `json:"player_4"`
	Wall          []string   `json:"wall"`
}

func View(stateMachine *game.StateMachine) *HumanView {
	g, s, a := stateMachine.View()
	return &HumanView{
		Id:            stateMachine.Id(),
		HasEnded:      s.Transition == nil,
		StateName:     s.Name,
		PrevalentWind: WindNames[g.GetPrevalentWind()],
		ActivePlayer:  int(g.GetActiveSeat()) + 1,
		ActiveDiscard: DescribeActiveDiscard(g.GetActiveDiscard()),
		Player1:       DescribePlayer(g, a, 0),
		Player2:       DescribePlayer(g, a, 1),
		Player3:       DescribePlayer(g, a, 2),
		Player4:       DescribePlayer(g, a, 3),
		Wall:          Describe(g.GetWall()),
	}
}

func DescribeActiveDiscard(t *model.Tile) string {
	if t == nil {
		return "none"
	}
	return model.TileNames[*t]
}

func DescribeCombinations(combinations []model.Combination) []string {
	descriptions := make([]string, len(combinations))
	for i, combi := range combinations {
		switch c := combi.(type) {
		case model.BonusTile:
			descriptions[i] = fmt.Sprintf("Bonus tile %s", model.TileNames[c.Tile])

		case model.Chow:
			descriptions[i] = fmt.Sprintf("Chow %s", model.TileNames[c.FirstTile])

		case model.Pung:
			descriptions[i] = fmt.Sprintf("Pung %s", model.TileNames[c.Tile])

		case model.Kong:
			descriptions[i] = fmt.Sprintf("Kong %s", model.TileNames[c.Tile])

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
		text := fmt.Sprintf("%d× %s", count, model.TileNames[tile])
		descriptions = append(descriptions, text)
	}
	return descriptions
}

func DescribePlayer(g model.Game, a map[model.Seat][]model.Action, player int) PlayerView {
	seat := model.Seat(player)
	actions, has := a[seat]
	if !has {
		actions = make([]model.Action, 0)
	}

	actionMap := make(map[int]string)
	for i, a := range actions {
		actionMap[i] = DescribeAction(a)
	}

	return PlayerView{
		Actions:   actionMap,
		Wind:      WindNames[g.GetPlayerAtSeat(seat).GetSeatWind()],
		Concealed: Describe(g.GetPlayerAtSeat(seat).GetConcealedTiles()),
		Exposed:   DescribeCombinations(g.GetPlayerAtSeat(seat).GetExposedCombinations()),
		Discarded: Describe(g.GetPlayerAtSeat(seat).GetDiscardedTiles()),
	}
}

func DescribeAction(action model.Action) string {
	switch a := action.(type) {
	case model.DoNothing:
		return "Do nothing"
	case model.Discard:
		return fmt.Sprintf("Discard a %s", model.TileNames[a.Tile])

	default:
		// This should not happen..!
		return "unknown action"
	}
}

func ToVector(t *model.TileCollection) []int {
	tileVector := make([]int, len(TileOrder))
	for i, tile := range TileOrder {
		tileVector[i] = t.NumOf(tile)
	}
	return tileVector
}
