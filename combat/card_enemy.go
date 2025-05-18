package combat

import (
	"maps"

	"github.com/jinzhu/copier"
)

type CardEnemy struct {
	CombatableBase
	Values  map[string]int `json:"values"`
	Effects []*Effect      `json:"effects"`
	Move    string         `json:"move"`
}

func NewCardEnemy(proto *CardEnemy) *CardEnemy {
	ret := &CardEnemy{
		CombatableBase: proto.CombatableBase,
		Values:         make(map[string]int, len(proto.Values)),
		Effects:        make([]*Effect, len(proto.Effects)),
		Move:           proto.Move,
	}
	ret.CombatType = ENEMY
	ret.MaxLife = ret.Life
	ret.Defense = 0
	ret.Strength = 0
	ret.Statuses = make([]*Status, 0)
	ret.WeakFactor = 75
	ret.VulnerableFactor = 150
	maps.Copy(ret.Values, proto.Values)
	for i, v := range ret.Effects {
		copier.Copy(v, proto.Effects[i])
	}
	return ret
}

func (a *CardEnemy) GetBase() *CombatableBase {
	return &a.CombatableBase
}

func (a *CardEnemy) GetValue(key string) int {
	return a.Values[key]
}
