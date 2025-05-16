package combat

import (
	"encoding/json"
	"maps"
	"reflect"

	"github.com/jinzhu/copier"
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

type Card struct {
	Id          string
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        int            `json:"type"`
	Values      map[string]int `json:"values"`
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
	return nil
}

func (c *Card) Copy() *Card {
	newCard := &Card{}
	copier.Copy(newCard, c)
	maps.Copy(newCard.Values, c.Values)
	return newCard
}
