package combat

import (
	"encoding/json"
	"maps"
	"my_test/util"
	"reflect"

	"github.com/jinzhu/copier"
	"github.com/samber/lo"
)

const (
	CARD_TYPE_ATTACK = iota
	CARD_TYPE_SKILL
	CARD_TYPE_POWER
)

const (
	CARD_RANGE_SINGLE = iota
	CARD_RANGE_ALL
	CARD_RANGE_RANDOM
)

const (
	CARD_RARITY_COMMON = iota
	CARD_RARITY_UNCOMMON
	CARD_RARITY_RARE
)

const (
	CARD_CHOOSE_FROM_HAND = iota
	CARD_CHOOSE_FROM_DRAW
	CARD_CHOOSE_FROM_DISCARD
	CARD_CHOOSE_FROM_EXHAUST
)

type Card struct {
	Id          string
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        int            `json:"type"`
	Values      map[string]any `json:"values"`
	Exhaust     bool           `json:"exhaust"`
	CanUse      string         `json:"can_use"`
	Range       int            `json:"range"`
	Rarity      int            `json:"rarity"`
	Cost        int            `json:"cost"`
	Price       int            `json:"price"`
	Effects     []*Effect      `json:"effects"`
	Upgraded    bool
	Binding     CardBinding
}

func (c *Card) UnmarshalJSON(data []byte) error {
	type TmpCard Card
	tmp := struct {
		*TmpCard
		HasBind bool `json:"bind"`
	}{
		TmpCard: (*TmpCard)(c),
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if tmp.HasBind {
		t := GetTower()
		c.Binding = reflect.New(t.cardBindingMap[c.Name]).Interface().(CardBinding)
	}

	for key, value := range c.Values {
		if floatVal, ok := value.(float64); ok {
			c.Values[key] = int(floatVal)
		}
	}
	return nil
}

func (c *Card) Copy() *Card {
	newCard := &Card{}
	copier.Copy(newCard, c)
	newCard.Values = make(map[string]any, len(c.Values))
	newCard.Effects = make([]*Effect, len(c.Effects))
	maps.Copy(newCard.Values, c.Values)
	for i, v := range c.Effects {
		newCard.Effects[i] = &Effect{}
		copier.Copy(newCard.Effects[i], v)
	}
	return newCard
}

func (c *Card) HasValue(key string) bool {
	return lo.HasKey(c.Values, key)
}

func (c *Card) GetValue(key string) int {
	return c.Values[key].(int)
}

func (c *Card) GetValueString(key string) string {
	return c.Values[key].(string)
}

func (c *Card) GetDescription() string {
	return util.FormatString(c.Description, c.Values)
}

func (c *Card) IsAttack() bool {
	return c.Type == CARD_TYPE_ATTACK
}

func (c *Card) IsSkill() bool {
	return c.Type == CARD_TYPE_SKILL
}

func (c *Card) IsPower() bool {
	return c.Type == CARD_TYPE_POWER
}
