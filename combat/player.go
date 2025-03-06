package combat

type Actor struct {
}

func NewActor(id int, name string) *Actor {
	return &Actor{}
}

func (u *Actor) AttackSpeed() int {
	return 1
}

func (u *Actor) CombatType() int {
	return ENEMY
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
