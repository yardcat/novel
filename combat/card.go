package combat

import (
	"encoding/json"
	"reflect"
)

const (
	CARD_TYPE_ATTACK = iota
	CARD_TYPE_SKILL
	CARD_TYPE_EFFECT
)

const (
	CARD_ATTACK_SINGLE = iota
	CARD_ATTACK_ALL
	CARD_RANDOM
)

const (
	CARD_RARITY_COMMON = iota
	CARD_RARITY_UNCOMMON
	CARD_RARITY_RARE
)

type Card struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        int            `json:"type"`
	Values      map[string]int `json:"values"`
	Disposal    bool           `json:"disposal"`
	Rarity      int            `json:"rarity"`
	Cost        int            `json:"cost"`
	Effects     []*Effect      `json:"effects"`
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
