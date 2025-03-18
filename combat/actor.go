package combat

import "my_test/log"

type Actor struct {
	CombatableBase
	Magic int
}

func NewActor(id int, name string) *Actor {
	return &Actor{
		CombatableBase: CombatableBase{
			Name:        name,
			CombatType:  ACTOR,
			Life:        100,
			Attack:      10,
			Defense:     2,
			Dodge:       10,
			AttackSpeed: 10,
			AttackRange: 6,
			AttackStep:  0,
		},
		Magic: 100,
	}
}

func (a *Actor) GetBase() *CombatableBase {
	return &a.CombatableBase
}

func (a *Actor) OnKill(defender Combatable) {
	log.Info("%s kill %s", a.GetName(), defender.GetName())
}
