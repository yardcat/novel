package world

type Home struct {
	Width  int
	Height int
}

func NewHome() *Home {
	return &Home{
		Width:  5,
		Height: 5,
	}
}
