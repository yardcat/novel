package scene

import (
	"math/rand"
	"my_test/combat"
	"my_test/user"
)

type LineScene struct {
	players        []*user.Player
	monster_points []*MonsterPoint
	distance       int
}

type MonsterPoint struct {
	monster *user.Monster
	x       int
}

func NewLineScene(group []*user.Player, monsters []*user.Monster, distance int) *LineScene {
	return &LineScene{
		players:        group,
		distance:       distance,
		monster_points: PositeMonsters(monsters, distance),
	}
}

func PositeMonsters(monsters []*user.Monster, distance int) []*MonsterPoint {
	monster_size := len(monsters)
	monster_points := make([]*MonsterPoint, monster_size)
	last_x := 0
	for i := 0; i < monster_size; i++ {
		offset_x := rand.Intn(distance - last_x)
		last_x = offset_x
		monster_points[i] = &MonsterPoint{
			monster: monsters[i],
			x:       offset_x,
		}
	}
	return monster_points
}

func (l *LineScene) DoCombat() {
	monsters := make([]*user.Monster, len(l.monster_points))
	for i, monster_point := range l.monster_points {
		monsters[i] = monster_point.monster
	}

	var fightables []user.Fightable
	for _, player := range monsters {
		fightables = append(fightables, player)
	}
	combat := combat.NewCombat(l.players, fightables)
	combat.Start()
}
