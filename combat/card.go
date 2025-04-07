package combat

import (
	"my_test/util"
)

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

	STATUS_VULNERABLE = iota
	STATUS_WEAK
	STATUS_STRENGTH
	STATUS_DEFENSE
)

type CardEffect struct {
	Effect int `json:"effect"`
	Value  int `json:"value"`
}

type Card struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Type        int          `json:"type"`
	Rarity      int          `json:"rarity"`
	Cost        int          `json:"cost"`
	Upgrade     []*Card      `json:"upgrade,omitempty"`
	Effects     []CardEffect `json:"effects"`
}

type CardSystem struct {
	cards   map[string]*Card
	deck    []*Card
	hand    []*Card
	discard []*Card
	remove  []*Card
}

func NewCardSystem() *CardSystem {
	return &CardSystem{
		cards: make(map[string]*Card),
		deck:  make([]*Card, 0),
	}
}

func (c *CardSystem) GetCard(name string) *Card {
	return c.cards[name]
}

func (c *CardSystem) AddCard() {
}

func (c *CardSystem) DrawCard(n int) []*Card {
	cards := make([]*Card, 0, n)

	for i := 0; i < n; i++ {
		if len(c.deck) == 0 {
			if len(c.discard) == 0 {
				break
			}
			c.ShuffleDiscardPile()
		}
		card := c.deck[0]
		c.deck = c.deck[1:]
		c.hand = append(c.hand, card)
		cards = append(cards, card)
	}

	return cards
}

func (c *CardSystem) DiscardCard(card *Card) {
	for i, v := range c.hand {
		if v == card {
			c.hand = append(c.hand[:i], c.hand[i+1:]...)
			c.discard = append(c.discard, card)
			return
		}
	}
}

func (c *CardSystem) ShuffleDiscardPile() {
	c.deck = append(c.deck, c.discard...)
	c.discard = nil

	for i := len(c.deck) - 1; i > 0; i-- {
		j := util.GetRandomInt(i + 1)
		c.deck[i], c.deck[j] = c.deck[j], c.deck[i]
	}
}

func (c *CardSystem) Use(cards []*Card, actors []Combatable, enemies []Combatable) {
	for _, card := range cards {
		for _, effect := range card.Effects {
			switch effect.Effect {
			case DAMAGE:
				for _, enemy := range enemies {
					enemy.OnDamage(effect.Value, nil)
				}
			case VULNERABLE:
				for _, enemy := range enemies {
					enemy.GetBase().AddStatus(Status{
						Type:  STATUS_VULNERABLE,
						Value: effect.Value,
						Turn:  2,
					})
				}
			case DEFEND:
				for _, actor := range actors {
					actor.GetBase().AddStatus(Status{
						Type:  STATUS_DEFENSE,
						Value: effect.Value,
						Turn:  1,
					})
				}
			case ADD_CARD:
			case MULTI_DAMAGE:
				n := effect.Value
				for i := 0; i < n; i++ {
					for _, enemy := range enemies {
						enemy.OnDamage(effect.Value, nil)
					}
				}
			case DAMAGE_DEFENSE:
				for _, enemy := range enemies {
					enemy.OnDamage(effect.Value, nil)
				}
				for _, actor := range actors {
					actor.GetBase().AddStatus(Status{
						Type:  STATUS_DEFENSE,
						Value: effect.Value,
						Turn:  1,
					})
				}
			case WEAK:
				for _, enemy := range enemies {
					enemy.GetBase().AddStatus(Status{
						Type:  STATUS_WEAK,
						Value: effect.Value,
						Turn:  2,
					})
				}
			case STRENGTH:
				for _, actor := range actors {
					actor.GetBase().AddStatus(Status{
						Type:  STATUS_STRENGTH,
						Value: effect.Value,
						Turn:  1,
					})
				}
			}
		}
	}
}
