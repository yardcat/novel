package combat

type EnemyGroup struct {
	Enemies []Enemy
}

type Dungeon struct {
	name   string
	groups []EnemyGroup
}
