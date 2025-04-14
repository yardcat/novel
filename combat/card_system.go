package combat

import (
	"encoding/json"
	"my_test/util"
	"os"
	"path/filepath"
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
	HEALTH

	STATUS_VULNERABLE = iota
	STATUS_WEAK
	STATUS_STRENGTH
	STATUS_DEFENSE
)

type CardEffect struct {
	Effect string `json:"effect"`
	Value  any    `json:"value"`
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

type CardCareer struct {
	InitCards []*Card
}

type CardEvent struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Effects     []CardEffect `json:"effects"`
}

type CardTurnInfo struct {
	Cards        []string
	DrawCount    int
	DiscardCount int
	RemoveCount  int
	Energy       int
}

func (c *CardCombat) GetCard(name string) *Card {
	return c.cardMap[name]
}

func (c *CardCombat) GetCardTurnInfo() *CardTurnInfo {
	info := &CardTurnInfo{
		Cards:        make([]string, len(c.Hand)),
		DrawCount:    len(c.deck),
		DiscardCount: len(c.discard),
		RemoveCount:  len(c.remove),
		Energy:       c.energy,
	}
	for i, card := range c.Hand {
		info.Cards[i] = card.Name
	}
	return info
}

func (c *CardCombat) PrepareCard() {
	career := c.careerMap["kongfu"]
	c.deck = append(c.deck, career.InitCards...)
	c.ShuffleDeck()
	c.DrawCard(c.maxCard)
}

func (c *CardCombat) GenerateChooseEvents() []string {
	return []string{"strength", "max_health", "draw_card"}
}

func (c *CardCombat) HandleChooseEvents(event string) {
	for _, e := range c.eventMap[event].Effects {
		switch c.EffectFromString(e.Effect) {
		case STRENGTH:
			c.actors[0].Strength += e.Value.(int)
		case HEALTH:
			c.actors[0].Life += e.Value.(int)
		case ADD_CARD:
			cards := e.Value.([]string)
			for _, card := range cards {
				c.AddCard(c.GetCard(card))
			}
		}
	}
}

func (c *CardCombat) AddCard(card *Card) {
	c.Hand = append(c.Hand, card)
}

func (c *CardCombat) DrawCard(n int) []*Card {
	cards := make([]*Card, 0, n)

	for i := 0; i < n; i++ {
		if len(c.deck) == 0 {
			if len(c.discard) == 0 {
				break
			}
			c.deck = append(c.deck, c.discard...)
			c.discard = make([]*Card, 0)
			c.ShuffleDeck()
		}
		card := c.deck[0]
		c.deck = c.deck[1:]
		c.Hand = append(c.Hand, card)
		cards = append(cards, card)
	}

	return cards
}

func (c *CardCombat) DiscardCard(card *Card) {
	for i, v := range c.Hand {
		if v == card {
			c.Hand = append(c.Hand[:i], c.Hand[i+1:]...)
			c.discard = append(c.discard, card)
			return
		}
	}
}

func (c *CardCombat) RemoveCard(card *Card) {
	c.remove = append(c.remove, card)
}

func (c *CardCombat) ShuffleDeck() {
	for i := len(c.deck) - 1; i > 0; i-- {
		j := util.GetRandomInt(i + 1)
		c.deck[i], c.deck[j] = c.deck[j], c.deck[i]
	}
}

func (c *CardCombat) EffectFromString(effect string) int {
	switch effect {
	case "damage":
		return DAMAGE
	case "add_card":
		return ADD_CARD
	case "multi_damage":
		return MULTI_DAMAGE
	case "damage_defense":
		return DAMAGE_DEFENSE
	case "vulnerable":
		return VULNERABLE
	case "defend":
		return DEFEND
	case "weak":
		return WEAK
	case "strength":
		return STRENGTH
	case "health":
		return HEALTH
	default:
		return 0
	}
}

func (c *CardCombat) Use(cards []*Card, actors []Combatable, enemies []Combatable) {
	for _, card := range cards {
		for _, effect := range card.Effects {
			switch c.EffectFromString(effect.Effect) {
			case DAMAGE:
				for _, enemy := range enemies {
					enemy.OnDamage(effect.Value.(int), nil)
				}
			case VULNERABLE:
				for _, enemy := range enemies {
					enemy.GetBase().AddStatus(Status{
						Type:  STATUS_VULNERABLE,
						Value: effect.Value.(int),
						Turn:  2,
					})
				}
			case DEFEND:
				for _, actor := range actors {
					actor.GetBase().AddStatus(Status{
						Type:  STATUS_DEFENSE,
						Value: effect.Value.(int),
						Turn:  1,
					})
				}
			case ADD_CARD:
			case MULTI_DAMAGE:
				n := effect.Value.(int)
				for i := 0; i < n; i++ {
					for _, enemy := range enemies {
						enemy.OnDamage(effect.Value.(int), nil)
					}
				}
			case DAMAGE_DEFENSE:
				for _, enemy := range enemies {
					enemy.OnDamage(effect.Value.(int), nil)
				}
				for _, actor := range actors {
					actor.GetBase().AddStatus(Status{
						Type:  STATUS_DEFENSE,
						Value: effect.Value.(int),
						Turn:  1,
					})
				}
			case WEAK:
				for _, enemy := range enemies {
					enemy.GetBase().AddStatus(Status{
						Type:  STATUS_WEAK,
						Value: effect.Value.(int),
						Turn:  2,
					})
				}
			case STRENGTH:
				for _, actor := range actors {
					actor.GetBase().AddStatus(Status{
						Type:  STATUS_STRENGTH,
						Value: effect.Value.(int),
						Turn:  1,
					})
				}
			}
		}
	}
}

func (c *CardCombat) loadData(dir string) error {
	err := c.loadCard(filepath.Join(dir, "card.json"))
	if err != nil {
		return err
	}
	err = c.loadCareer(filepath.Join(dir, "career.json"))
	if err != nil {
		return err
	}
	err = c.loadEvent(filepath.Join(dir, "event.json"))
	if err != nil {
		return err
	}
	return nil
}

func (c *CardCombat) loadCard(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &c.cardMap)
	if err != nil {
		return err
	}
	return nil
}

func (c *CardCombat) loadCareer(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	adapter := map[string]struct {
		Cards []string `json:"init_cards"`
	}{}
	err = json.Unmarshal(data, &adapter)
	for k, v := range adapter {
		c.careerMap[k] = &CardCareer{
			InitCards: make([]*Card, 0, len(v.Cards)),
		}
		for _, card := range v.Cards {
			c.careerMap[k].InitCards = append(c.careerMap[k].InitCards, c.cardMap[card])
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *CardCombat) loadEvent(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &c.eventMap)
	if err != nil {
		return err
	}
	return nil
}
