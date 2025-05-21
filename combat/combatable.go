package combat

import (
	"github.com/samber/lo"
	"github.com/samber/lo/mutable"
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
	Buffs            []*Buff
	WeakFactor       int
	VulnerableFactor int
}

func NewCombatableBase(id int, name string) *CombatableBase {
	return &CombatableBase{
		Name:  name,
		Buffs: make([]*Buff, 0),
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

func (c *CombatableBase) AddBuff(buff Buff) {
	var idx int = -1
	for i, v := range c.Buffs {
		if buff.Type == v.Type {
			idx = i
		}
	}
	if idx != -1 {
		v := c.Buffs[idx]
		switch v.Name {
		case BUFF_VULNERABLE:
			fallthrough
		case BUFF_WEAK:
			v.Turn = max(v.Turn, buff.Turn)
		case BUFF_STRENGTH:
			fallthrough
		case BUFF_ARMOR:
			v.Value += buff.Value
		case BUFF_POISON:
			idx = -1
		}
	}
	if idx == -1 {
		if buff.Turn == 0 {
			buff.Turn = 1
		}
		c.Buffs = append(c.Buffs, &buff)
	}
}

func (c *CombatableBase) RemoveBuff(name string) {
	c.Buffs = mutable.Filter(c.Buffs, func(v *Buff) bool {
		return v.Name != name
	})
}

func (c *CombatableBase) HasBuff(name string) bool {
	return lo.ContainsBy(c.Buffs, func(v *Buff) bool {
		return v.Name == name
	})
}

func (c *CombatableBase) GetArmorBuff() *Buff {
	for _, v := range c.Buffs {
		if v.Name == BUFF_ARMOR {
			return v
		}
	}
	return nil
}

func (c *CombatableBase) GetBuffValue(name string) int {
	for _, v := range c.Buffs {
		if v.Name == name {
			return v.Value
		}
	}
	return 0
}

func (c *CombatableBase) UpdateBuffs() {
	c.Buffs = lo.Filter(c.Buffs, func(v *Buff, _ int) bool {
		v.Turn--
		return v.Turn > 0
	})
}

func (c *CombatableBase) GetArmor() int {
	for _, v := range c.Buffs {
		if v.Name == BUFF_ARMOR {
			return v.Value
		}
	}
	return 0
}
