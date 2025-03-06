package combat

type Actor struct {
	name string
}

func NewActor(id int, name string) *Actor {
	return &Actor{
		name: name,
	}
}

func (u *Actor) AttackSpeed() int {
	return 1
}

func (u *Actor) CombatType() int {
	return ENEMY
}

func (u *Actor) OnAttack(op Combatable) {
}

func (u *Actor) OnDamage(damage int, op Combatable) {
	u.property.Life -= damage
}

func (u *Actor) GetName() string {
	return u.name
}

func (u *Actor) GetLife() int {
	return u.property.Life
}

func (u *Actor) GetAttack() int {
	return u.property.Attack
}

func (u *Actor) GetDefense() int {
	return u.property.Defense
}

func (u *Actor) GetDodge() int {
	return u.property.Dodge
}

func (u *Actor) IsAlive() bool {
	return u.property.Life > 0
}
