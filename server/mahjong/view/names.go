package view

import (
	"fmt"
	"github.com/roelofruis/mahjong-learn/mahjong"
	"github.com/roelofruis/mahjong-learn/state"
)

var tileNames = map[mahjong.Tile]string{
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

var windNames = map[mahjong.Wind]string{
	mahjong.East:  "East",
	mahjong.South: "South",
	mahjong.West:  "West",
	mahjong.North: "North",
}

func tileName(t *mahjong.Tile) string {
	if t == nil {
		return "none"
	}
	return tileNames[*t]
}

func combinationNames(combinations []mahjong.Combination) []string {
	descriptions := make([]string, len(combinations))
	for i, combi := range combinations {
		switch c := combi.(type) {
		case mahjong.BonusTile:
			descriptions[i] = fmt.Sprintf("Bonus tile %s", tileNames[c.Tile])

		case mahjong.Chow:
			descriptions[i] = fmt.Sprintf("Chow %s", tileNames[c.FirstTile])

		case mahjong.Pung:
			descriptions[i] = fmt.Sprintf("Pung %s", tileNames[c.Tile])

		case mahjong.Kong:
			descriptions[i] = fmt.Sprintf("Kong %s", tileNames[c.Tile])

		default:
			// This should not happen..!
			descriptions[i] = "unknown combination"
		}
	}
	return descriptions
}

func tileCollectionNames(t *mahjong.TileCollection) []string {
	orderedCounts := t.OrderedCounts()
	descriptions := make([]string, len(orderedCounts))
	for i, count := range orderedCounts {
		descriptions[i] = fmt.Sprintf("%d× %s", count.Count, tileNames[count.Tile])
	}
	return descriptions
}

func actionNames(action state.Action) string {
	switch a := action.(type) {
	case mahjong.Discard:
		return fmt.Sprintf("Discard a %s", tileNames[a.Tile])
	case mahjong.DeclareConcealedKong:
		return fmt.Sprintf("Declare a concealed Kong of %s", tileNames[a.Tile])
	case mahjong.ExposedPungToKong:
		return fmt.Sprintf("Add to exposed pung")
	case mahjong.DoNothing:
		return "Do nothing"
	case mahjong.DeclareChow:
		return fmt.Sprintf("Declare chow up from %s", tileNames[a.Tile])
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
