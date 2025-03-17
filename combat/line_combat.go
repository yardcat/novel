package combat

import (
	"my_test/log"
)

type LineCombat struct {
	Combat
}

func NewLineCombat(actors []*Actor, enemies []*Enemy, client CombatClient) *LineCombat {
	return &LineCombat{
		Combat: *NewCombat(actors, enemies, client),
	}
}

func (c *LineCombat) ChooseAttacker() Combatable {
	fast := MAX_STEP
	fast_idx := 0
	for i, comb := range c.combatables {
		speed := (MAX_STEP - comb.GetBase().AttackStep) / float64(comb.GetAttackSpeed())
		if speed < fast {
			fast = speed
			fast_idx = i
		}
	}
	for _, comb := range c.combatables {
		comb.GetBase().AttackStep += float64(comb.GetAttackSpeed()) * fast
	}
	c.combatables[fast_idx].GetBase().AttackStep = 0
	return c.combatables[fast_idx]
}

func (c *LineCombat) ChooseDefender(attacker Combatable) Combatable {
	if attacker.GetCombatType() == ACTOR {
		return c.enemies[0]
	} else if attacker.GetCombatType() == ENEMY {
		return c.actors[0]
	}
	log.Info("unknown attacker type %d", attacker.GetCombatType())
	return nil
}
