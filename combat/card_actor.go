package combat

type CardActor struct {
	CombatableBase
	Magic  int
	Energy int
}

func NewCardActor(base CombatableBase) *CardActor {
	a := &CardActor{
		CombatableBase: base,
		Magic:          100,
	}
	a.Statuses = make([]*Status, 0)
	a.MaxLife = a.Life
	return a
}

func (a *CardActor) GetBase() *CombatableBase {
	return &a.CombatableBase
}
