package combat

type EnemyGroup struct {
	Name    string
	Enemies []*Enemy
}

type Dungeon struct {
	Name   string
	Groups []EnemyGroup
}

func CreateEnemyGroup(prototype EnemyGroup) EnemyGroup {
	ret := EnemyGroup{
		Name:    prototype.Name,
		Enemies: make([]*Enemy, len(prototype.Enemies)),
	}

	for i, v := range prototype.Enemies {
		ret.Enemies[i] = CreateEnemy(v)
	}
	return ret
}
