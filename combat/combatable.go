package combat

const (
	ACTOR = 1
	ENEMY = 2
)

type Combatable interface {
	GetBase() *CombatableBase
	GetName() string
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
	Name        string
	AttackSpeed int
	AttackStep  float64
	Attack      int
	Life        int
	Defense     int
	Dodge       int
	CombatType  int
}

func NewCombatableBase(id int, name string) *CombatableBase {
	return &CombatableBase{
		Name: name,
	}
}

func (c *CombatableBase) GetName() string {
	return c.Name
}

func (c *CombatableBase) GetAttackSpeed() int {
	return c.AttackSpeed
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
	c.Life -= damage
}

func (c *CombatableBase) OnDead(Combatable) {

}

func (c *CombatableBase) OnKill(Combatable) {

}
