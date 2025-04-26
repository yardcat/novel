package combat

type Enemy struct {
	CombatableBase
	behaviors int
}

func CreateEnemy(proto *Enemy) *Enemy {
	ret := &Enemy{
		CombatableBase: proto.CombatableBase,
	}
	ret.CombatType = ENEMY
	ret.AttackStep = 0
	ret.MaxLife = ret.Life
	ret.Statuses = make(map[int]*Status)
	return ret
}

func (a *Enemy) GetBase() *CombatableBase {
	return &a.CombatableBase
}

func (c *Enemy) OnKill(Combatable) {

}
