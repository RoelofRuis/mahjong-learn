package model

type Action interface {
	// ActionIndex, has to be unique among all defined actions (to guarantee a stable sorting)
	ActionIndex() int
}

type ByIndex []Action

func (a ByIndex) Len() int           { return len(a) }
func (a ByIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIndex) Less(i, j int) bool { return a[i].ActionIndex() < a[j].ActionIndex() }


type Discard struct {
	Tile Tile
}

func (d Discard) ActionIndex() int { return int(d.Tile) }

type DoNothing struct{}

func (d DoNothing) ActionIndex() int { return 100 }
