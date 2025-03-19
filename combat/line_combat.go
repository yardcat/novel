package combat

import (
	"my_test/log"
)

type LineCombat struct {
	*Combat
}

func NewLineCombat(combat *Combat) *LineCombat {
	return &LineCombat{
		combat,
	}
}

func (c *LineCombat) ChooseDefender(attacker Combatable) Combatable {
	if attacker.GetCombatType() == ACTOR {
		return c.enemies[0]
	} else if attacker.GetCombatType() == ENEMY {
		return c.actors[0]
	} else {
		log.Info("unknown attacker type %d", attacker.GetCombatType())
	}
	return nil
}
