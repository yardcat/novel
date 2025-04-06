package combat

const (
	MAX_STEP = 100.0
)

type CombatClient interface {
	OnLose()
	OnWin()
	OnDraw()
	OnKill(Combatable)
	OnDead(Combatable)
}

type Record struct {
	actorCastDamage  int
	actorIncurDamage int
	turns            int
}

type CombatResult struct {
	LifeCost  int
	MagicCost int
}

type CombatLayout interface {
	ChooseDefender(attacker Combatable) Combatable
}

type CombatParams struct {
	Actors  []*Actor
	Enemies []*Enemy
	Client  CombatClient
}

type CombatOnceResult struct {
	attackerDead bool
	defenderDead bool
}

type Combat interface {
	Start()
	Actors() []*Actor
	Enemies() []*Enemy
	Combatables() []Combatable
}
