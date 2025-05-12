package combat

import "github.com/samber/lo"

type Shop struct {
	Potions   []*Potion
	Relics    []*Relic
	Cards     []*Card
	PotionNum int
	RelicNum  int
	CardNum   int
}

func (r *Shop) Type() int {
	return ROOM_TYPE_SHOP
}

func NewShop(t *Tower) *Shop {
	s := &Shop{
		PotionNum: 3,
		RelicNum:  3,
		CardNum:   4,
	}
	s.Potions = generateItem(t.PotionMap, s.PotionNum)
	s.Relics = generateItem(t.RelicMap, s.RelicNum)
	s.Cards = generateItem(t.cardMap, s.CardNum)

	return s
}

func generateItem[T any](mp map[string]T, num int) []T {
	keys := lo.Values(mp)
	return lo.Samples(keys, num)
}
