package combat

type Shop struct {
	Potions []*Potion
	Relics  []*Relic
	Cards   []*Card
}

func (r *Shop) Type() int {
	return ROOM_TYPE_SHOP
}
