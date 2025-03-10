package combat

type EnemyGroup struct {
	Name    string
	Enemies []*Enemy
}

type Dungeon struct {
	Name   string
	Groups []EnemyGroup
}
