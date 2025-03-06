package combat

import (
	"fmt"
	"my_test/event"
	"my_test/log"
	"my_test/user"
)

const (
	MAX_STEP = 100.0
	ACTOR    = 1
	ENEMY    = 2
)

type Combatable interface {
	GetAttackSpeed() int
	GetAttack() int
	CombatType() int
	IsAlive() bool
	OnAttack(defender Combatable)
	OnDamage(damage int, attacker Combatable)
	OnDead(killer Combatable)
}

type CombatClient interface {
	OnLost()
	OnWin()
	OnKill(killer Combatable)
	OnDead()
}

type Combat struct {
	attackStep  map[Combatable]float64
	combatables []Combatable
	actors      []*Actor
	enemies     []*Enemy
	record      string
	client      CombatClient
}

type CombatOnceResult struct {
	attackerDead bool
	defenderDead bool
}

func (c *Combat) Start() {
	for len(c.actors) > 0 && len(c.enemies) > 0 {
		attacker := c.ChooseAttacker()
		defender := c.ChooseDefender(attacker)
		isActorAttacker := attacker.CombatType() == ACTOR
		result := c.CombatOnce(attacker, defender, isActorAttacker)
		if result.attackerDead {
			fmt.Println(attacker.OnDead(defender), "dead")
			c.removeCombatable(attacker)
		}
		if result.defenderDead {
			fmt.Println(defender.OnDead(attacker), "dead")
			c.removeCombatable(defender)
		}
	}
	if len(c.actors) == 0 {
		fmt.Println(defender.GetName(), "dead")
		event.GetEventBus().OnEvent(event.Die, map[string]any{"player": player})
	} else if len(c.enemies) == 0 {
		fmt.Println(attacker.GetName(), "win")
		event.GetEventBus().OnEvent(event.Die, map[string]any{"player": player})
	}
}

func (c *Combat) ChooseAttacker() Combatable {
	fast := MAX_STEP
	fast_idx := 0
	for i, comb := range c.combatables {
		fast = min(fast, (MAX_STEP-c.attackStep[comb])/float64(comb.GetAttackSpeed()))
		fast_idx = i
	}
	for _, comb := range c.combatables {
		c.attackStep[comb] += float64(comb.GetAttackSpeed()) * fast
	}
	return c.combatables[fast_idx]
}

func (c *Combat) ChooseDefender(attacker Combatable) Combatable {
	if attacker.CombatType() == ACTOR {
		return c.enemies[0]
	} else if attacker.CombatType() == ENEMY {
		return c.actors[0]
	}
	log.Info("unknown attacker type %d", attacker.CombatType())
	return nil
}

func (c *Combat) CombatOnce(attacker Combatable, defender Combatable, isActorAttacker bool) CombatOnceResult {
	attacker.OnAttack(defender)
	damage_reduce_factor := getDamageFactor(attacker, op)
	damage := int(float32(attacker.GetAttack()) * damage_reduce_factor)
	if shouldDodge(attacker, defender) {
		damage = 0
	}
	defender.OnDamage(damage)
	return CombatOnceResult{
		attackerDead: attacker.IsAlive(),
		defenderDead: defender.IsAlive(),
	}
}

func getDamageFactor(player user.Fightable, op user.Fightable) float32 {
	return 1.0
}

func shouldDodge(player user.Fightable, op user.Fightable) bool {
	return false
}

func (c *Combat) removeCombatable(combatable Combatable) {
	switch combatable.CombatType() {
	case ACTOR:
		for i, actor := range c.actors {
			if actor == combatable {
				c.actors = append(c.actors[:i], c.actors[i+1:]...)
				return
			}
		}
	case ENEMY:
		for i, enemy := range c.enemies {
			if enemy == combatable {
				c.enemies = append(c.enemies[:i], c.enemies[i+1:]...)
				return
			}
		}
	default:
		log.Info("unknown combatable type %d", combatable.CombatType())
	}
}
