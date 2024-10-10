package user

type Player struct {
	id       int
	name     string
	property Property
}

func NewPlayer(id int, name string) *Player {
	return &Player{
		id:   id,
		name: name,
		property: Property{
			Life:    100,
			Attack:  10,
			Defense: 5,
			Dodge:   10,
		},
	}
}

func (u *Player) Attack(op Fightable) {
}

func (u *Player) Damage(damage int) {
	u.property.Life -= damage
}

func (u *Player) GetName() string {
	return u.name
}

func (u *Player) GetLife() int {
	return u.property.Life
}

func (u *Player) GetAttack() int {
	return u.property.Attack
}

func (u *Player) GetDefense() int {
	return u.property.Defense
}

func (u *Player) GetDodge() int {
	return u.property.Dodge
}

func (u *Player) IsAlive() bool {
	return u.property.Life > 0
}
