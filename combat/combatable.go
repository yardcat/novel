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
	Statuses    map[int]*Status
}

func NewCombatableBase(id int, name string) *CombatableBase {
	return &CombatableBase{
		Name:     name,
		Statuses: map[int]*Status{},
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

	defense := c.GetStatusValue(STATUS_ARMOR)
	damage = max(0, damage-defense)

	c.Life -= damage
}

func (c *CombatableBase) OnDead(Combatable) {

}

func (c *CombatableBase) OnKill(Combatable) {

}

func (c *CombatableBase) AddStatus(status Status) {
	v, exist := c.Statuses[status.Type]
	if exist {
		switch v.Type {
		case STATUS_VULNERABLE:
		case STATUS_WEAK:
			v.Turn = max(v.Turn, status.Turn)
		case STATUS_STRENGTH:
		case STATUS_ARMOR:
			v.Value += status.Value
		}
	}
	if !exist {
		ns := &Status{}
		*ns = status
		c.Statuses[status.Type] = ns
	}
}

func (c *CombatableBase) RemoveStatus(statusType int) {
	_, exist := c.Statuses[statusType]
	if exist {
		delete(c.Statuses, statusType)
	}
}

func (c *CombatableBase) GetStatusValue(statusType int) int {
	_, exist := c.Statuses[statusType]
	if !exist {
		return 0
	}
	return c.Statuses[statusType].Value
}

func (c *CombatableBase) UpdateStatus() {
	for k, v := range c.Statuses {
		v.Turn--
		if v.Turn <= 0 {
			delete(c.Statuses, k)
		}
	}
}
