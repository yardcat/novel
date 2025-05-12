package combat

type fuhua struct {
}

func (f *fuhua) Init() {

}

func (f *fuhua) GetCost(t *Tower) int {
	return max(3-t.currentCombat.hurtCount, 1)
}

func (f *fuhua) GetDamage(t *Tower) int {
	return t.currentCombat.turnCount
}

func (f *fuhua) Modify(t *Tower) {

}

func (f *fuhua) CanPlay() string {
	return "fuhua"
}
