package combat

type CardActor struct {
	CombatableBase
	Energy int
}

func NewCardActor(base CombatableBase) *CardActor {
	a := &CardActor{
		CombatableBase: base,
	}
	a.CombatType = ACTOR
	a.Statuses = make([]*Status, 0)
	a.MaxLife = a.Life
	return a
}

func (a *CardActor) GetBase() *CombatableBase {
	return &a.CombatableBase
}
