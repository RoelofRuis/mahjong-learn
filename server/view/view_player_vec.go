package view

import "github.com/roelofruis/mahjong-learn/mahjong"

type PlayerVec struct {
	Score            int       `json:"score"`
	BonusTiles       []int     `json:"bonus_tiles"`
	PrevalentWind    []int     `json:"prevalent_wind"`
	PlayerWind       []int     `json:"player_wind"`
	DiscardingPlayer []int     `json:"discarding_player"`
	ActiveDiscard    []int     `json:"active_discard"`
	Received         []int     `json:"received_tile"`
	Concealed        [][]int   `json:"concealed_tiles"`
	ExposedChows     [][][]int `json:"exposed_chows"`
	ExposedPungs     [][]int   `json:"exposed_pungs"`
	ExposedKongs     [][]int   `json:"exposed_kongs"`
	HiddenKongs      [][]int   `json:"hidden_kongs"`
	Discards         [][]int   `json:"discards"`
	// player to the right
	PlayerRScore        int       `json:"right_player_score"`
	PlayerRBonusTiles   []int     `json:"right_player_bonus_tiles"`
	PlayerRWind         []int     `json:"right_player_wind"`
	PlayerRExposedChows [][][]int `json:"right_player_exposed_chows"`
	PlayerRExposedPungs [][]int   `json:"right_player_exposed_pungs"`
	PlayerRExposedKongs [][]int   `json:"right_player_exposed_kongs"`
	PlayerRHiddenKongs  [][]int   `json:"right_player_hidden_kongs"`
	PlayerRDiscards     [][]int   `json:"right_player_discards"`
	// opposite player
	PlayerOScore        int       `json:"opposite_player_score"`
	PlayerOBonusTiles   []int     `json:"opposite_player_bonus_tiles"`
	PlayerOWind         []int     `json:"opposite_player_wind"`
	PlayerOExposedChows [][][]int `json:"opposite_player_exposed_chows"`
	PlayerOExposedPungs [][]int   `json:"opposite_player_exposed_pungs"`
	PlayerOExposedKongs [][]int   `json:"opposite_player_exposed_kongs"`
	PlayerOHiddenKongs  [][]int   `json:"opposite_player_hidden_kongs"`
	PlayerODiscards     [][]int   `json:"opposite_player_discards"`
	// player to the left
	PlayerLScore        int       `json:"left_player_score"`
	PlayerLBonusTiles   []int     `json:"left_player_bonus_tiles"`
	PlayerLWind         []int     `json:"left_player_wind"`
	PlayerLExposedChows [][][]int `json:"left_player_exposed_chows"`
	PlayerLExposedPungs [][]int   `json:"left_player_exposed_pungs"`
	PlayerLExposedKongs [][]int   `json:"left_player_exposed_kongs"`
	PlayerLHiddenKongs  [][]int   `json:"left_player_hidden_kongs"`
	PlayerLDiscards     [][]int   `json:"left_player_discards"`
}

func ViewPlayerVec(game *mahjong.Game, playerIndex int) *PlayerVec {
	table := *game.Table
	player := table.GetPlayerByIndex(playerIndex)

	discardingPlayer := []int{0, 0, 0}
	activeDiscard := tileToVec(nil)
	if table.GetActiveDiscard() != nil {
		discardingPlayer = discardingPlayerVec(table, playerIndex)
		activeDiscard = tileToVec(table.GetActiveDiscard())
	}

	chows, pungs, kongs, hiddenKongs := exposedCombinations(player.GetExposedCombinations())

	playerRIndex := (playerIndex + 1) % 4
	playerR := table.GetPlayerByIndex(playerRIndex)
	rChows, rPungs, rKongs, rHiddenKongs := exposedCombinations(playerR.GetExposedCombinations())
	playerOIndex := (playerIndex + 2) % 4
	playerO := table.GetPlayerByIndex(playerOIndex)
	oChows, oPungs, oKongs, oHiddenKongs := exposedCombinations(playerO.GetExposedCombinations())
	playerLIndex := (playerIndex + 3) % 4
	playerL := table.GetPlayerByIndex(playerLIndex)
	lChows, lPungs, lKongs, lHiddenKongs := exposedCombinations(playerL.GetExposedCombinations())

	return &PlayerVec{
		Score:               player.GetScore(),
		BonusTiles:          bonusTiles(player.GetExposedCombinationCollection()),
		PrevalentWind:       WindVectors[table.GetPrevalentWind()],
		PlayerWind:          WindVectors[player.GetWind()],
		DiscardingPlayer:    discardingPlayer,
		ActiveDiscard:       activeDiscard,
		Received:            tileToVec(player.GetReceivedTile()),
		Concealed:           collectionToVec(player.GetConcealedTiles(), 13),
		ExposedChows:        chows,
		ExposedPungs:        pungs,
		ExposedKongs:        kongs,
		HiddenKongs:         hiddenKongs,
		Discards:            collectionToVec(player.GetDiscardedTiles(), 40), // TODO: maybe this can be lower, determine worst case.
		PlayerRScore:        playerR.GetScore(),
		PlayerRBonusTiles:   bonusTiles(playerR.GetExposedCombinationCollection()),
		PlayerRWind:         WindVectors[playerR.GetWind()],
		PlayerRExposedChows: rChows,
		PlayerRExposedPungs: rPungs,
		PlayerRExposedKongs: rKongs,
		PlayerRHiddenKongs:  rHiddenKongs,
		PlayerRDiscards:     collectionToVec(playerR.GetDiscardedTiles(), 40),
		PlayerOScore:        playerO.GetScore(),
		PlayerOBonusTiles:   bonusTiles(playerO.GetExposedCombinationCollection()),
		PlayerOWind:         WindVectors[playerO.GetWind()],
		PlayerOExposedChows: oChows,
		PlayerOExposedPungs: oPungs,
		PlayerOExposedKongs: oKongs,
		PlayerOHiddenKongs:  oHiddenKongs,
		PlayerODiscards:     collectionToVec(playerO.GetDiscardedTiles(), 40),
		PlayerLScore:        playerL.GetScore(),
		PlayerLBonusTiles:   bonusTiles(playerL.GetExposedCombinationCollection()),
		PlayerLWind:         WindVectors[playerL.GetWind()],
		PlayerLExposedChows: lChows,
		PlayerLExposedPungs: lPungs,
		PlayerLExposedKongs: lKongs,
		PlayerLHiddenKongs:  lHiddenKongs,
		PlayerLDiscards:     collectionToVec(playerL.GetDiscardedTiles(), 40),
	}
}

var WindVectors = map[mahjong.Wind][]int{
	mahjong.East:  {1, 0, 0, 0},
	mahjong.South: {0, 1, 0, 0},
	mahjong.West:  {0, 0, 1, 0},
	mahjong.North: {0, 0, 0, 1},
}

var DiscardingPlayerVectors = map[int][]int{
	0: {0, 0, 0},
	1: {1, 0, 0},
	2: {0, 1, 0},
	3: {0, 0, 1},
}

func collectionToVec(coll *mahjong.TileCollection, maxLen int) [][]int {
	vector := make([][]int, maxLen)
	tileIndex := 0
	for _, t := range TileOrder {
		if t.IsBonusTile() {
			// we know the order, so break early on first bonus tile occurrence
			break
		}
		tileVec := tileToVec(&t)
		for i := coll.NumOf(t); i > 0; i-- {
			vector[tileIndex] = tileVec
			tileIndex++
		}
	}
	for ; tileIndex < maxLen; tileIndex++ {
		vector[tileIndex] = tileToVec(nil)
	}
	return vector
}

func exposedCombinations(combinations []mahjong.Combination) ([][][]int, [][]int, [][]int, [][]int) {
	chowVector := make([][][]int, 4)
	chowIndex := 0
	pungVector := make([][]int, 4)
	pungIndex := 0
	kongVector := make([][]int, 4)
	kongIndex := 0
	hiddenKongVector := make([][]int, 4)
	hiddenKongIndex := 0
	for _, combination := range combinations {
		switch c := combination.(type) {
		case mahjong.Chow:
			chowVector[chowIndex] = [][]int{
				tileToVec(&c.FirstTile),
				tileToVec(c.FirstTile.NextInSuit()),
				tileToVec(c.FirstTile.NextInSuit().NextInSuit()),
			}
			chowIndex++
		case mahjong.Pung:
			pungVector[pungIndex] = tileToVec(&c.Tile)
			pungIndex++
		case mahjong.Kong:
			if c.Concealed == true {
				hiddenKongVector[hiddenKongIndex] = tileToVec(&c.Tile)
				hiddenKongIndex++
			} else {
				kongVector[kongIndex] = tileToVec(&c.Tile)
				kongIndex++
			}
		}
	}
	for i := 0; i < 4; i++ {
		if i >= chowIndex {
			chowVector[i] = [][]int{tileToVec(nil), tileToVec(nil), tileToVec(nil)}
		}
		if i >= pungIndex {
			pungVector[i] = tileToVec(nil)
		}
		if i >= kongIndex {
			kongVector[i] = tileToVec(nil)
		}
		if i >= hiddenKongIndex {
			hiddenKongVector[i] = tileToVec(nil)
		}
	}
	return chowVector, pungVector, kongVector, hiddenKongVector
}

func bonusTiles(t *mahjong.CombinationCollection) []int {
	return []int{
		numBonus(t, mahjong.FlowerPlumb),
		numBonus(t, mahjong.FlowerOrchid),
		numBonus(t, mahjong.FlowerChrysanthemum),
		numBonus(t, mahjong.FlowerBamboo),
		numBonus(t, mahjong.SeasonSpring),
		numBonus(t, mahjong.SeasonSummer),
		numBonus(t, mahjong.SeasonAutumn),
		numBonus(t, mahjong.SeasonWinter),
	}
}

func numBonus(t *mahjong.CombinationCollection, tile mahjong.Tile) int {
	if t.Contains(mahjong.BonusTile{Tile: tile}) {
		return 1
	}
	return 0
}

func discardingPlayerVec(table mahjong.Table, playerIndex int) []int {
	return DiscardingPlayerVectors[((table.GetActivePlayerIndex()-playerIndex)+4)%4]
}

func tileToVec(t *mahjong.Tile) []int {
	if t == nil {
		return []int{0, 0, 0}
	}
	if t.IsBonusTile() {
		panic("cannot express bonus tile as tile vector")
	}
	tile := int(*t)
	group := 0
	tpe := 0
	nr := 0
	if t.IsSuit() {
		tpe = tile / 10
		nr = tile % 10
	}
	if t.IsDragon() {
		group = 2
		tpe = tile % 10
	}
	if t.IsWind() {
		group = 3
		tpe = tile % 10
	}
	return []int{group, tpe, nr}
}
