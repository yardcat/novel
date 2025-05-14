package combat

import (
	"encoding/json"
	"reflect"

	"github.com/jinzhu/copier"
)

const (
	CARD_TYPE_ATTACK = iota
	CARD_TYPE_SKILL
	CARD_TYPE_EFFECT
)

const (
	CARD_RARITY_COMMON = iota
	CARD_RARITY_UNCOMMON
	CARD_RARITY_RARE
)

type Card struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Type        int           `json:"type"`
	Rarity      int           `json:"rarity"`
	Cost        int           `json:"cost"`
	Effects     []*CardEffect `json:"effects"`
	Binding     CardBinding
}

func (c *Card) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Name        string        `json:"name"`
		Description string        `json:"description"`
		Type        int           `json:"type"`
		Rarity      int           `json:"rarity"`
		Cost        int           `json:"cost"`
		Effects     []*CardEffect `json:"effects"`
		HasBind     bool          `json:"bind"`
	}{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	copier.Copy(c, &tmp)

	if tmp.HasBind {
		t := GetTower()
		c.Binding = reflect.New(t.cardBindingMap[c.Name]).Interface().(CardBinding)
	}
	return nil
}
