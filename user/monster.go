package user

type Monster struct {
	name     string
	property Property
}

func NewMonster(name string, property Property) *Monster {
	return &Monster{
		name:     name,
		property: property,
	}
}

func (u *Monster) Attack(op Fightable) {
}

func (u *Monster) Damage(damage int) {
	u.property.Life -= damage
}

func (u *Monster) GetName() string {
	return u.name
}

func (u *Monster) GetLife() int {
	return u.property.Life
}

func (u *Monster) GetAttack() int {
	return u.property.Attack
}

func (u *Monster) GetDefense() int {
	return u.property.Defense
}

func (u *Monster) GetDodge() int {
	return u.property.Dodge
}

func (m *Monster) IsAlive() bool {
	return m.property.Life > 0
}
