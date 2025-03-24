package world

import (
	"my_test/combat"
	"my_test/equip"
	"my_test/util"
)

type PlayerAttr struct {
	combat.CombatableBase
	Exp         int
	Health      int
	Hunger      int
	Thirst      int
	Energy      int
	LevelExp    int
	LevelExpInc int
	LevelUp     map[string]int
}

func (p *PlayerAttr) UpdateFromAttr(attrs []equip.Attr, percent bool) {
	for _, attr := range attrs {
		if percent && attr.Value.Type != util.Percent {
			continue
		} else if !percent && attr.Value.Type == util.Percent {
			continue
		}
		switch attr.Name {
		case "Life":
			p.Life = attr.Value.Int()
		case "Attack":
			p.Attack = attr.Value.Int()
		case "Defense":
			p.Defense = attr.Value.Int()
		case "Dodge":
			p.Dodge = attr.Value.Int()
		case "AttackSpeed":
			p.AttackSpeed = attr.Value.Int()
		case "AttackRange":
			p.AttackRange = attr.Value.Int()
		case "Exp":
			p.Exp = attr.Value.Int()
		case "Health":
			p.Health = attr.Value.Int()
		case "Hunger":
			p.Hunger = attr.Value.Int()
		case "Thirst":
			p.Thirst = attr.Value.Int()
		case "Energy":
			p.Energy = attr.Value.Int()
		case "LevelExp":
			p.LevelExp = attr.Value.Int()
		case "LevelExpInc":
			p.LevelExpInc = attr.Value.Int()
		}
	}
}

func (p *PlayerAttr) UpdateFinal(add *PlayerAttr, percent *PlayerAttr) {
	p.Attack = CacAttr(p.Attack, add.Attack, percent.Attack)
	p.Life = CacAttr(p.Life, add.Life, percent.Life)
	p.Defense = CacAttr(p.Defense, add.Defense, percent.Defense)
	p.Dodge = CacAttr(p.Dodge, add.Dodge, percent.Dodge)
	p.AttackSpeed = CacAttr(p.AttackSpeed, add.AttackSpeed, percent.AttackSpeed)
	p.AttackRange = CacAttr(p.AttackRange, add.AttackRange, percent.AttackRange)
	p.Exp = CacAttr(p.Exp, add.Exp, percent.Exp)
	p.Health = CacAttr(p.Health, add.Health, percent.Health)
	p.Hunger = CacAttr(p.Hunger, add.Hunger, percent.Hunger)
	p.Thirst = CacAttr(p.Thirst, add.Thirst, percent.Thirst)
	p.Energy = CacAttr(p.Energy, add.Energy, percent.Energy)
	p.LevelExp = CacAttr(p.LevelExp, add.LevelExp, percent.LevelExp)
	p.LevelExpInc = CacAttr(p.LevelExpInc, add.LevelExpInc, percent.LevelExpInc)
}

func CacAttr[T int | float64](value T, add T, percent T) T {
	value = value * (percent + 100) / 100
	value += add
	return value
}
