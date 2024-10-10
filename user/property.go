package user

type Fightable interface {
	Attack(target Fightable)
	Damage(damage int)

	GetName() string
	GetLife() int
	GetAttack() int
	GetDefense() int
	GetDodge() int
	IsAlive() bool
}

type Property struct {
	Life        int
	Attack      int
	Defense     int
	Dodge       int
	AttackSpeed int
}
