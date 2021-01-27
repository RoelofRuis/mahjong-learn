package model

type Action interface {
	// ActionIndex, has to be unique among all defined actions (to guarantee a stable sorting)
	ActionIndex() int
}

type Discard struct {
	Tile Tile
}

func (d Discard) ActionIndex() int { return int(d.Tile) }

type DoNothing struct{}

func (d DoNothing) ActionIndex() int { return 100 }
