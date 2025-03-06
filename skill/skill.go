package skill

type Caster interface {
}

type Skill struct {
	name   string
	damage int
	cost   int
	caster Caster
}

func (*Skill) NewSkill(name string, damage int, cost int, caster Caster) *Skill {
	return &Skill{
		name:   name,
		damage: damage,
		cost:   cost,
		caster: caster,
	}
}

func (s *Skill) Cast(target Caster) {

}

func (s *Skill) Update() {

}
