package combat

import "my_test/log"

type ActorClient interface {
	OnCombatDone(result CombatResult)
}

type Actor struct {
	CombatableBase
	Magic  int
	Energy int
	client ActorClient
}

func NewActor(base CombatableBase, c ActorClient) *Actor {
	a := &Actor{
		CombatableBase: base,
		Magic:          100,
		client:         c,
	}
	a.Statuses = make(map[int]*Status)
	a.MaxLife = a.Life
	return a
}

func (a *Actor) GetBase() *CombatableBase {
	return &a.CombatableBase
}

func (a *Actor) OnKill(defender Combatable) {
	log.Info("%s kill %s", a.GetName(), defender.GetName())
}

func (a *Actor) OnCombatDone(result CombatResult) {
	a.client.OnCombatDone(result)
}
