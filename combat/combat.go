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
	AttackSpeed() int
	CombatType() int
}

type Combat struct {
	attackStep  map[Combatable]float64
	combatables []Combatable
	actors      []*Actor
	enemies     []*Enemy
	record      string
}

func (c *Combat) Start() {
	for len(c.actors) > 0 && len(c.enemies) > 0 {
		attacker := c.ChooseAttacker()
		defender := c.ChooseDefender(attacker)
		if attacker.CombatType() == ACTOR {
			c.ActorAttack(attacker.(*Actor), defender.(*Enemy))
		} else if attacker.CombatType() == ENEMY {
			c.EnemyAttack(defender.(*Enemy), attacker.(*Actor))
		}
	}
	if len {
		fmt.Println(defender.GetName(), "dead")
		event.GetEventBus().OnEvent(event.Die, map[string]any{"player": player})
	}
}

func (c *Combat) ChooseAttacker() Combatable {
	fast := MAX_STEP
	fast_idx := 0
	for i, comb := range c.combatables {
		fast = min(fast, (MAX_STEP-c.attackStep[comb])/float64(comb.AttackSpeed()))
		fast_idx = i
	}
	for _, comb := range c.combatables {
		c.attackStep[comb] += float64(comb.AttackSpeed()) * fast
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

func (c *Combat) ActorAttack(actor *Actor, enemy *Enemy) {
	player.Attack(op)
	damage_reduce_factor := getDamageFactor(player, op)
	damage := int(float32(player.GetAttack()) * damage_reduce_factor)
	if shouldDodge(player, op) {
		damage = 0
	}
	op.Damage(damage)
}

func (c *Combat) EnemyAttack(enemy *Enemy, actor *Actor) {
	player.Attack(op)
	damage_reduce_factor := getDamageFactor(player, op)
	damage := int(float32(player.GetAttack()) * damage_reduce_factor)
	if shouldDodge(player, op) {
		damage = 0
	}
	op.Damage(damage)
}

func getDamageFactor(player user.Fightable, op user.Fightable) float32 {
	return 1.0
}

func shouldDodge(player user.Fightable, op user.Fightable) bool {
	return false
}
