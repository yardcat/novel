package combat

import (
	"fmt"
	"my_test/log"
	"my_test/util"
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

type Record struct {
	actorCastDamage  int
	actorIncurDamage int
	turns            int
}

type CombatResult struct {
	LifeCost  int
	MagicCost int
}

type CombatStrategy interface {
	ChooseDefender(attacker Combatable) Combatable
}

type Combat struct {
	combatables []Combatable
	actors      []*Actor
	enemies     []*Enemy
	client      CombatClient
	strategy    CombatStrategy
	Record
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
	c.strategy = NewGridCombat(c)
	return c
}

func (c *Combat) Start() {
	for len(c.actors) > 0 && len(c.enemies) > 0 {
		attacker := c.ChooseAttacker()
		defender := c.strategy.ChooseDefender(attacker)
		if defender == nil {
			log.Info("attacker %s can't find defender", attacker.GetName())
			continue
		}
		isActorAttacker := attacker.GetCombatType() == ACTOR
		result := c.CombatOnce(attacker, defender, isActorAttacker)
		c.turns++
		if result.attackerDead {
			defender.OnKill(attacker)
			attacker.OnDead(defender)
			c.removeCombatable(attacker)
		}
		if result.defenderDead {
			attacker.OnKill(defender)
			defender.OnDead(attacker)
			c.removeCombatable(defender)
		}
	}
	if len(c.actors) != 0 {
		fmt.Println("win")
		c.client.OnWin()
	} else if len(c.enemies) != 0 {
		fmt.Println("lose")
		c.client.OnLose()
	} else {
		fmt.Println("draw")
		c.client.OnDraw()
	}
	c.onCombatFinish()
}

func (c *Combat) ChooseAttacker() Combatable {
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

func (c *Combat) CombatOnce(attacker Combatable, defender Combatable, isActorAttacker bool) CombatOnceResult {
	attacker.OnAttack(defender)
	damage := c.cacDamage(attacker, defender)
	if c.shouldDodge(attacker, defender) {
		log.Info("%s dodge damage %d on %s", attacker.GetName(), damage, defender.GetName())
		return CombatOnceResult{
			attackerDead: false,
			defenderDead: false,
		}
	}
	defender.OnDamage(damage, attacker)
	if isActorAttacker {
		c.actorCastDamage += damage
	} else {
		c.actorIncurDamage += damage

	}
	log.Info("%s cast damage %d on %s", attacker.GetName(), damage, defender.GetName())
	return CombatOnceResult{
		attackerDead: !attacker.IsAlive(),
		defenderDead: !defender.IsAlive(),
	}
}

func (c *Combat) cacDamage(attacker Combatable, defender Combatable) int {
	attack := attacker.GetAttack()
	defense := defender.GetDefense()
	damage_reduce_factor := 0.0
	damage := int(float64(attack-defense) * (1 - damage_reduce_factor))
	return max(damage, 0)
}

func (c *Combat) shouldDodge(_ Combatable, defender Combatable) bool {
	randomNumber := util.GetRandomInt(100)
	dodge := defender.GetDodge()
	return randomNumber < dodge
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
			return
		}
	}
}

func (c *Combat) onCombatFinish() {
	log.Info("combat finish, turns %d, actor cast %d damaage, actor incur %d damage", c.turns, c.actorCastDamage, c.actorIncurDamage)
	for _, actor := range c.actors {
		result := CombatResult{LifeCost: c.actorIncurDamage}
		actor.OnCombatDone(result)
	}
}

func (c *Combat) getEnemyAsCombatable() []Combatable {
	enemies := make([]Combatable, len(c.enemies))
	for i, enemy := range c.enemies {
		enemies[i] = enemy
	}
	return enemies
}

func (c *Combat) getActorAsCombatable() []Combatable {
	actors := make([]Combatable, len(c.actors))
	for i, actor := range c.actors {
		actors[i] = actor
	}
	return actors
}
