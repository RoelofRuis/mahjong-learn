package main

func InitGame(id uint64) Game {
	players := make(map[int]Player, 4)

	tiles := NewTileCollection()

	return Game{
		Id: id,
		HasEnded: false,
		Wall: tiles,
		Players: players,
	}
}

func NewTileCollection() TileCollection {
	c := TileCollection{Tiles: make(map[Tile]int)}
	c.Add(Bamboo1, 4)
	c.Add(Bamboo2, 4)
	c.Add(Bamboo3, 4)
	c.Add(Bamboo4, 4)
	c.Add(Bamboo5, 4)
	c.Add(Bamboo6, 4)
	c.Add(Bamboo7, 4)
	c.Add(Bamboo8, 4)
	c.Add(Bamboo9, 4)
	c.Add(Circles1, 4)
	c.Add(Circles2, 4)
	c.Add(Circles3, 4)
	c.Add(Circles4, 4)
	c.Add(Circles5, 4)
	c.Add(Circles6, 4)
	c.Add(Circles7, 4)
	c.Add(Circles8, 4)
	c.Add(Circles9, 4)
	c.Add(Characters1, 4)
	c.Add(Characters2, 4)
	c.Add(Characters3, 4)
	c.Add(Characters4, 4)
	c.Add(Characters5, 4)
	c.Add(Characters6, 4)
	c.Add(Characters7, 4)
	c.Add(Characters8, 4)
	c.Add(Characters9, 4)
	c.Add(RedDragon, 4)
	c.Add(GreenDragon, 4)
	c.Add(WhiteDragon, 4)
	c.Add(EastWind, 4)
	c.Add(SouthWind, 4)
	c.Add(NorthWind, 4)
	c.Add(WestWind, 4)
	c.Add(FlowerPlumb, 1)
	c.Add(FlowerOrchid, 1)
	c.Add(FlowerChrysanthemum, 1)
	c.Add(FlowerBamboo, 1)
	c.Add(SeasonSpring, 1)
	c.Add(SeasonSummer, 1)
	c.Add(SeasonAutumn, 1)
	c.Add(SeasonWinter, 1)
	return c
}

func (t *TileCollection) Add(tile Tile, n int) {
	t.Tiles[tile] = n
}