package combat

import (
	"my_test/log"
)

type LineLayout struct {
	*AutoCombat
}

func NewLineLayout(combat *AutoCombat) *LineLayout {
	return &LineLayout{
		combat,
	}
}

func (c *LineLayout) ChooseDefender(attacker Combatable) Combatable {
	if attacker.GetCombatType() == ACTOR {
		return c.enemies[0]
	} else if attacker.GetCombatType() == ENEMY {
		return c.actors[0]
	} else {
		log.Info("unknown attacker type %d", attacker.GetCombatType())
	}
	return nil
}
