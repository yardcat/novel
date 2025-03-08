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
	attackStep  map[Combatable]float64
	combatables []Combatable
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
		attackStep:  make(map[Combatable]float64),
	}
	for _, actor := range actors {
		c.combatables = append(c.combatables, actor)
	}
	for _, enemy := range enemies {
		c.combatables = append(c.combatables, enemy)
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
		fast = min(fast, (MAX_STEP-c.attackStep[comb])/float64(comb.GetAttackSpeed()))
		fast_idx = i
	}
	for _, comb := range c.combatables {
		c.attackStep[comb] += float64(comb.GetAttackSpeed()) * fast
	}
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
	damage := int(float32(attacker.GetAttack()) * damage_reduce)
	if shouldDodge(attacker, defender) {
		damage = 0
	}
	defender.OnDamage(damage, attacker)
	return CombatOnceResult{
		attackerDead: attacker.IsAlive(),
		defenderDead: defender.IsAlive(),
	}
}

func getDamageReduce(attacker Combatable, enemy Combatable) float32 {
	return 1.0
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
		log.Info("unknown combatable type %d", combatable.GetCombatType())
	}
}
