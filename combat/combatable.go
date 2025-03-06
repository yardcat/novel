package combat

type Combatable interface {
	GetName() string
	GetAttackSpeed() int
	GetAttack() int
	GetLife() int
	GetDefense() int
	GetDodge() int

	CombatType() int
	IsAlive() bool
	OnAttack(defender Combatable)
	OnDamage(damage int, attacker Combatable)
	OnDead(Combatable)
	OnKill(Combatable)
}

type CombatableBase struct {
	name        string
	attackSpeed int
	attack      int
	life        int
	defense     int
	dodge       int
	combatType  int
}

func NewCombatableBase(id int, name string) *CombatableBase {
	return &CombatableBase{
		name: name,
	}
}

func (c *CombatableBase) GetName() string {
	return c.name
}

func (c *CombatableBase) GetAttackSpeed() int {
	return c.attackSpeed
}
func (c *CombatableBase) GetAttack() int {
	return c.attack
}
func (c *CombatableBase) GetLife() int {
	return c.life
}

func (c *CombatableBase) GetDefense() int {
	return c.defense
}

func (c *CombatableBase) GetDodge() int {
	return c.dodge
}

func (c *CombatableBase) CombatType() int {
	return c.combatType
}

func (c *CombatableBase) IsAlive() bool {
	return c.life > 0
}

func (c *CombatableBase) OnAttack(defender Combatable) {
}

func (c *CombatableBase) OnDamage(damage int, attacker Combatable) {
}

func (c *CombatableBase) OnDead(Combatable) {

}

func (c *CombatableBase) OnKill(Combatable) {

}
