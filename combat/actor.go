package combat

type Actor struct {
	CombatableBase
}

func NewActor(id int, name string) *Actor {
	return &Actor{
		CombatableBase: CombatableBase{
			Name:        name,
			CombatType:  ACTOR,
			Life:        100,
			Attack:      10,
			Defense:     10,
			Dodge:       10,
			AttackSpeed: 10,
		},
	}
}
