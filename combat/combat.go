package combat

const (
	MAX_STEP = 100.0
)

type CombatClient interface {
	OnLose()
	OnWin()
	OnKill(Combatable)
}

type CombatResult struct {
	LifeCost  int
	MagicCost int
}

type CombatLayout interface {
	ChooseDefender(attacker Combatable) Combatable
}

type PathProvider interface {
	GetPath(path string) string
}

type AutoCombatParams struct {
	Actors  []*Actor
	Enemies []*Enemy
	Path    PathProvider
	Client  CombatClient
}

type CombatOnceResult struct {
	attackerDead bool
	defenderDead bool
}

type Combat interface {
	Start()
	Actors() []Combatable
	Enemies() []Combatable
	Combatables() []Combatable
}

type Status struct {
	Type  int
	Value int
	Turn  int
}
