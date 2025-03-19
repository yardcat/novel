package world

import "my_test/combat"

type Pet struct {
	combat.CombatableBase
	Name string
	Type int
}

func NewPet(name string, petType int) *Pet {
	return &Pet{
		Name: name,
		Type: petType,
	}
}

func CreatePet(proto *Pet) *Pet {
	ret := &Pet{
		CombatableBase: proto.CombatableBase,
	}
	ret.CombatType = combat.ACTOR
	ret.AttackStep = 0
	return ret
}

func (p *Pet) GetCombatableBase() combat.CombatableBase {
	return p.CombatableBase
}
