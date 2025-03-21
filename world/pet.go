package world

import (
	"my_test/combat"
	"my_test/log"
)

type Pet struct {
	combat.CombatableBase
	Type int
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

func (p *Pet) OnCombatDone(result combat.CombatResult) {
	log.Info("pet combat done %v", result)
}
