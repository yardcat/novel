package skill

import "my_test/user"

type Skill struct {
	name   string
	damage int
	cost   int
	caster user.Fightable
}

func (*Skill) NewSkill(name string, damage int, cost int, caster user.Fightable) *Skill {
	return &Skill{
		name:   name,
		damage: damage,
		cost:   cost,
		caster: caster,
	}
}

func (s *Skill) Cast(target user.Fightable) {

}

func (s *Skill) Update() {

}
