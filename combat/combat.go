package combat

import (
	"fmt"
	"my_test/event"
	"my_test/user"
)

type Combat struct {
	group  []*user.Player
	ops    []user.Fightable
	turn   int
	record string
}

func NewCombat(group []*user.Player, ops []user.Fightable) *Combat {
	return &Combat{
		group: group,
		ops:   ops,
	}
}

func (c *Combat) Start() {
	for _, duel := range c.ops {
		player := c.group[0]
		fmt.Println(player.GetName(), " vs ", duel.GetName())
		c.Fight(player, duel)
	}
}

func (c *Combat) Fight(player user.Fightable, op user.Fightable) {
	attacker := player
	defender := op
	for attacker.IsAlive() && defender.IsAlive() {
		attacker, defender = defender, attacker
		c.turn = c.turn + 1
		c.FightOnce(attacker, defender)
	}
	if !defender.IsAlive() {
		fmt.Println(defender.GetName(), "dead")
		event.GetEventBus().OnEvent(event.Die, map[string]any{"player": player})
	}
}

func (c *Combat) ChooseAttacker() user.Fightable {
	return nil
}

func (c *Combat) ChooseDefender() user.Fightable {
	return nil
}

func (c *Combat) FightOnce(player user.Fightable, op user.Fightable) {
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
