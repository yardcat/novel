package combat

type Actor struct {
	CombatableBase
}

func NewActor(id int, name string) *Actor {
	return &Actor{
		CombatableBase: CombatableBase{
			name: name,
		},
	}
}
