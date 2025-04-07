package combat

const (
	CARD_TYPE_ATTACK = iota
	CARD_TYPE_SKILL
	CARD_TYPE_EFFECT

	CARD_RARITY_COMMON = iota
	CARD_RARITY_UNCOMMON
	CARD_RARITY_RARE

	DAMAGE = iota
	VULNERABLE
	DEFEND
	ADD_CARD
	MULTI_DAMAGE
	DAMAGE_DEFENSE
	WEAK
	STRENGTH
)

type CardEffect struct {
	Effect int
	Value  int
}

type Card struct {
	Type    int
	Rarity  int
	Effects []CardEffect
}

type CardSystem struct {
	cards   map[string]*Card
	hand    []*Card
	discard []*Card
	remove  []*Card
}

func NewCardSystem() *CardSystem {
	return &CardSystem{
		cards: make(map[string]*Card),
	}
}

func (c *CardSystem) GetCard(name string) *Card {
	return c.cards[name]
}

func (c *CardSystem) AddCard() {
}
