package world

type Grid struct {
}

type MapSystem struct {
	Width  int
	Height int
	Grids  []*Grid
}
