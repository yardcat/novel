package world

type Fish struct {
}

func NewFish() *Fish {
	return &Fish{}
}

func (f *Fish) PassBy() {

}

func (f *Fish) Explore() int {
	return 0
}
