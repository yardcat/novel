package combat

import (
	"slices"

	"github.com/samber/lo"
)

const (
	ACTOR = 1
	ENEMY = 2
)

type Combatable interface {
	GetBase() *CombatableBase
	GetName() string
	GetAttackRange() int
	GetAttackSpeed() int
	GetAttack() int
	GetLife() int
	GetDefense() int
	GetDodge() int

	GetCombatType() int
	IsAlive() bool
	OnAttack(defender Combatable)
	OnDamage(damage int, attacker Combatable)
	OnDead(Combatable)
	OnKill(Combatable)
}

type CombatableBase struct {
	Name             string
	AttackSpeed      int
	AttackRange      int
	AttackStep       float64
	Attack           int
	Life             int
	MaxLife          int
	Dodge            int
	CombatType       int
	Strength         int
	Defense          int
	Buffs            []*Status
	WeakFactor       int
	VulnerableFactor int
}

func NewCombatableBase(id int, name string) *CombatableBase {
	return &CombatableBase{
		Name:  name,
		Buffs: make([]*Status, 0),
	}
}

func (c *CombatableBase) GetName() string {
	return c.Name
}

func (c *CombatableBase) GetAttackSpeed() int {
	return c.AttackSpeed
}

func (c *CombatableBase) GetAttackRange() int {
	return c.AttackRange
}

func (c *CombatableBase) GetAttack() int {
	return c.Attack
}
func (c *CombatableBase) GetLife() int {
	return c.Life
}

func (c *CombatableBase) GetDefense() int {
	return c.Defense
}

func (c *CombatableBase) GetDodge() int {
	return c.Dodge
}

func (c *CombatableBase) GetCombatType() int {
	return c.CombatType
}

func (c *CombatableBase) IsAlive() bool {
	return c.Life > 0
}

func (c *CombatableBase) OnAttack(defender Combatable) {
}

func (c *CombatableBase) OnDamage(damage int, attacker Combatable) {
}

func (c *CombatableBase) OnDead(Combatable) {

}

func (c *CombatableBase) OnKill(Combatable) {

}

func (c *CombatableBase) AddStatus(status Status) {
	var idx int = -1
	for i, v := range c.Buffs {
		if status.Type == v.Type {
			idx = i
		}
	}
	if idx != -1 {
		v := c.Buffs[idx]
		switch v.Type {
		case STATUS_VULNERABLE:
		case STATUS_WEAK:
			v.Turn = max(v.Turn, status.Turn)
		case STATUS_STRENGTH:
		case STATUS_ARMOR:
			v.Value += status.Value
		case STATUS_POISON:
			idx = -1
		}
	}
	if idx == -1 {
		c.Buffs = append(c.Buffs, &status)
	}
}

func (c *CombatableBase) RemoveStatus(statusType int) {
	for i, v := range c.Buffs {
		if v.Type == statusType {
			c.Buffs = slices.Delete(c.Buffs, i, i+1)
		}
	}
}

func (c *CombatableBase) HasStatus(statusType int) bool {
	return lo.ContainsBy(c.Buffs, func(v *Status) bool {
		return v.Type == statusType
	})
}

func (c *CombatableBase) GetArmorStatus() *Status {
	for _, v := range c.Buffs {
		if v.Type == STATUS_ARMOR {
			return v
		}
	}
	return nil
}

func (c *CombatableBase) GetStatusValue(statusType int) int {
	for _, v := range c.Buffs {
		if v.Type == statusType {
			return v.Value
		}
	}
	return 0
}

func (c *CombatableBase) UpdateStatus() {
	c.Buffs = lo.Filter(c.Buffs, func(v *Status, _ int) bool {
		return v.Turn > 0
	})
}

func (c *CombatableBase) GetArmor() int {
	for _, v := range c.Buffs {
		if v.Type == STATUS_ARMOR {
			return v.Value
		}
	}
	return 0
}
