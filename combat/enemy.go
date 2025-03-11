package combat

type Enemy struct {
	CombatableBase
}

func CreateEnemy(proto *Enemy) *Enemy {
	ret := &Enemy{
		CombatableBase: proto.CombatableBase,
	}
	ret.CombatType = ENEMY
	ret.AttackStep = 0
	return ret
}

func (a *Enemy) GetBase() *CombatableBase {
	return &a.CombatableBase
}

func (c *Enemy) OnKill(Combatable) {

}
