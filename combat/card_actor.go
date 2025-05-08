package combat

type CardActor struct {
	CombatableBase
	Energy        int
	InitEnergy    int
	InitCardCount int
}

func NewCardActor(base CombatableBase) *CardActor {
	a := &CardActor{
		CombatableBase: base,
		Energy:         3,
		InitEnergy:     3,
	}
	a.CombatType = ACTOR
	a.Statuses = make([]*Status, 0)
	a.MaxLife = a.Life
	return a
}

func (a *CardActor) GetBase() *CombatableBase {
	return &a.CombatableBase
}
