package combat

const (
	ACTOR = 1
	ENEMY = 2
)

type Combatable interface {
	GetBase() *CombatableBase
	GetName() string
	GetAttackRange() int
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
	AttackRange int
	AttackStep  float64
	Attack      int
	Life        int
	MaxLife     int
	Dodge       int
	CombatType  int
	Strength    int
	Defense     int
	Statuses    []Status
}

func NewCombatableBase(id int, name string) *CombatableBase {
	return &CombatableBase{
		Name:     name,
		Statuses: []Status{},
		MaxLife:  100,
		Life:     100,
	}
}

func (c *CombatableBase) GetName() string {
	return c.Name
}

func (c *CombatableBase) GetAttackSpeed() int {
	return c.AttackSpeed
}

func (c *CombatableBase) GetAttackRange() int {
	return c.AttackRange
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
	vulnerable := c.GetStatusValue(STATUS_VULNERABLE)
	if vulnerable > 0 {
		damage = int(float64(damage) * 1.5)
	}

	weak := c.GetStatusValue(STATUS_WEAK)
	if weak > 0 {
		damage = int(float64(damage) * 0.75)
	}

	defense := c.GetStatusValue(STATUS_DEFENSE)
	damage = max(0, damage-defense)

	c.Life -= damage
}

func (c *CombatableBase) OnDead(Combatable) {

}

func (c *CombatableBase) OnKill(Combatable) {

}

func (c *CombatableBase) AddStatus(status Status) {
	c.Statuses = append(c.Statuses, status)
}

func (c *CombatableBase) RemoveStatus(statusType int) {
	for i := len(c.Statuses) - 1; i >= 0; i-- {
		if c.Statuses[i].Type == statusType {
			c.Statuses = append(c.Statuses[:i], c.Statuses[i+1:]...)
		}
	}
}

func (c *CombatableBase) GetStatusValue(statusType int) int {
	value := 0
	for _, status := range c.Statuses {
		if status.Type == statusType {
			value += status.Turn
		}
	}
	return value
}

func (c *CombatableBase) UpdateStatus() {
	for i := len(c.Statuses) - 1; i >= 0; i-- {
		c.Statuses[i].Turn--
		if c.Statuses[i].Turn <= 0 {
			c.Statuses = append(c.Statuses[:i], c.Statuses[i+1:]...)
		}
	}
}
