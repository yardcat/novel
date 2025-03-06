package combat

type Enemy struct {
	CombatableBase
}

func NewEnemy(id int, name string) *Enemy {
	return &Enemy{
		CombatableBase: CombatableBase{
			name: name,
		},
	}
}
