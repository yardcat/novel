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
		InitCardCount:  5,
	}
	a.CombatType = ACTOR
	a.Buffs = make([]*Status, 0)
	a.MaxLife = a.Life
	a.Defense = 0
	a.Strength = 0
	a.WeakFactor = 75
	a.VulnerableFactor = 150
	return a
}

func (a *CardActor) GetBase() *CombatableBase {
	return &a.CombatableBase
}
