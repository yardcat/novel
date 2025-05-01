package combat

type CardEnemy struct {
	CombatableBase
}

type CardEnemyGroup struct {
	group []string
}

func NewCardEnemy(proto *CardEnemy) *CardEnemy {
	ret := &CardEnemy{
		CombatableBase: proto.CombatableBase,
	}
	ret.CombatType = ENEMY
	ret.MaxLife = ret.Life
	ret.Statuses = make([]*Status, 0)
	return ret
}

func (a *CardEnemy) GetBase() *CombatableBase {
	return &a.CombatableBase
}
