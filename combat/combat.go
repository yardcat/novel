package combat

const (
	MAX_STEP = 100.0
)

const (
	STATUS_BUFF = iota
	STATUS_DEBUFF
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

type Buff struct {
	Type  int
	Name  string
	Value int
	Turn  int
}
