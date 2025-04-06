package combat

import (
	"my_test/log"
)

type LineLayout struct {
	combat Combat
}

func NewLineCombat(combat Combat) *LineLayout {
	return &LineLayout{
		combat: combat,
	}
}

func (c *LineLayout) ChooseDefender(attacker Combatable) Combatable {
	if attacker.GetCombatType() == ACTOR {
		return c.combat.Enemies()[0]
	} else if attacker.GetCombatType() == ENEMY {
		return c.combat.Actors()[0]
	} else {
		log.Info("unknown attacker type %d", attacker.GetCombatType())
	}
	return nil
}
