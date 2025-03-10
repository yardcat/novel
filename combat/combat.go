package combat

import (
	"fmt"
	"my_test/log"
)

const (
	MAX_STEP = 100.0
)

type CombatClient interface {
	OnLose()
	OnWin()
	OnDraw()
	OnKill(Combatable)
	OnDead(Combatable)
}

type Combat struct {
	combatables []Combatable
	attackStep  []float64
	actors      []*Actor
	enemies     []*Enemy
	client      CombatClient
}

type CombatOnceResult struct {
	attackerDead bool
	defenderDead bool
}

func NewCombat(actors []*Actor, enemies []*Enemy, client CombatClient) *Combat {
	c := &Combat{
		actors:      actors,
		enemies:     enemies,
		client:      client,
		combatables: make([]Combatable, len(actors)+len(enemies)),
		attackStep:  make([]float64, len(actors)+len(enemies)),
	}
	i := 0
	for _, actor := range actors {
		c.combatables[i] = actor
		i++
	}
	for _, enemy := range enemies {
		c.combatables[i] = enemy
		i++
	}
	return c
}

func (c *Combat) Start() {
	for len(c.actors) > 0 && len(c.enemies) > 0 {
		attacker := c.ChooseAttacker()
		defender := c.ChooseDefender(attacker)
		isActorAttacker := attacker.GetCombatType() == ACTOR
		result := c.CombatOnce(attacker, defender, isActorAttacker)
		if result.attackerDead {
			defender.OnKill(attacker)
			attacker.OnDead(defender)
			c.removeCombatable(attacker)
		}
		if result.defenderDead {
			defender.OnKill(attacker)
			defender.OnDead(attacker)
			c.removeCombatable(defender)
		}
	}
	if len(c.actors) != 0 {
		fmt.Println("win")
		c.client.OnLose()
	} else if len(c.enemies) != 0 {
		fmt.Println("lose")
		c.client.OnWin()
	} else {
		fmt.Println("draw")
		c.client.OnDraw()
	}
}

func (c *Combat) ChooseAttacker() Combatable {
	fast := MAX_STEP
	fast_idx := 0
	for i, comb := range c.combatables {
		speed := (MAX_STEP - c.attackStep[i]) / float64(comb.GetAttackSpeed())
		if speed < fast {
			fast = speed
			fast_idx = i
		}
	}
	for i, comb := range c.combatables {
		c.attackStep[i] += float64(comb.GetAttackSpeed()) * fast
	}
	c.attackStep[fast_idx] = 0
	return c.combatables[fast_idx]
}

func (c *Combat) ChooseDefender(attacker Combatable) Combatable {
	if attacker.GetCombatType() == ACTOR {
		return c.enemies[0]
	} else if attacker.GetCombatType() == ENEMY {
		return c.actors[0]
	}
	log.Info("unknown attacker type %d", attacker.GetCombatType())
	return nil
}

func (c *Combat) CombatOnce(attacker Combatable, defender Combatable, isActorAttacker bool) CombatOnceResult {
	attacker.OnAttack(defender)
	damage_reduce := getDamageReduce(attacker, defender)
	damage := int(float32(attacker.GetAttack()) * (1 - damage_reduce))
	if shouldDodge(attacker, defender) {
		damage = 0
	}
	defender.OnDamage(damage, attacker)
	log.Info("%s cast damage %d on %s", attacker.GetName(), damage, defender.GetName())
	return CombatOnceResult{
		attackerDead: !attacker.IsAlive(),
		defenderDead: !defender.IsAlive(),
	}
}

func getDamageReduce(attacker Combatable, enemy Combatable) float32 {
	return 0
}

func shouldDodge(attacker Combatable, enemy Combatable) bool {
	return false
}

func (c *Combat) removeCombatable(combatable Combatable) {
	switch combatable.GetCombatType() {
	case ACTOR:
		for i, actor := range c.actors {
			if actor == combatable {
				c.actors = append(c.actors[:i], c.actors[i+1:]...)
				break
			}
		}
	case ENEMY:
		for i, enemy := range c.enemies {
			if enemy == combatable {
				c.enemies = append(c.enemies[:i], c.enemies[i+1:]...)
				break
			}
		}
	default:
		log.Info("unknown combatable type %d", combatable.GetCombatType())
	}
	for i, v := range c.combatables {
		if v == combatable {
			c.combatables = append(c.combatables[:i], c.combatables[i+1:]...)
			c.attackStep = append(c.attackStep[:i], c.attackStep[i+1:]...)
			return
		}
	}
}
