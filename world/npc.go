package world

import (
	"my_test/combat"
	"my_test/log"
)

type Npc struct {
	combat.CombatableBase
	Description string
}

func CreateNpc(proto *Npc) *Npc {
	ret := &Npc{
		CombatableBase: proto.CombatableBase,
	}
	ret.CombatType = combat.ACTOR
	ret.AttackStep = 0
	return ret
}

func (p *Npc) GetCombatableBase() combat.CombatableBase {
	return p.CombatableBase
}

func (p *Npc) OnCombatDone(result combat.CombatResult) {
	log.Info("npc combat done %v", result)
}
