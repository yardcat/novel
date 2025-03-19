package combat

import (
	"my_test/log"
	"my_test/util"
)

const (
	GRID_WIDTH  = 6
	GRID_HEIGHT = 6
)

type GridCombat struct {
	*Combat
	pos2comb map[int]Combatable
	comb2pos map[Combatable]int
}

func NewGridCombat(combat *Combat) *GridCombat {
	g := &GridCombat{
		combat,
		make(map[int]Combatable),
		make(map[Combatable]int),
	}
	g.placeCombatables()
	return g
}

func (g *GridCombat) getComb(x, y int) Combatable {
	idx := y*GRID_WIDTH + x
	return g.pos2comb[idx]
}

func (g *GridCombat) ChooseDefender(attacker Combatable) Combatable {
	var near Combatable
	if attacker.GetCombatType() == ACTOR {
		near = g.getNearDefender(attacker, g.getEnemyAsCombatable())
	} else if attacker.GetCombatType() == ENEMY {
		near = g.getNearDefender(attacker, g.getActorAsCombatable())
	} else {
		log.Info("unknown attacker type %d", attacker.GetCombatType())
	}
	return near
}

func (g *GridCombat) getNearDefender(attacker Combatable, defenders []Combatable) Combatable {
	attackRange := attacker.GetAttackRange()
	if attackRange == 0 {
		attackRange = GRID_WIDTH + GRID_HEIGHT
	}
	var near Combatable
	var minDistance int
	for _, defender := range defenders {
		distance := g.getDistance(attacker, defender)
		if distance <= attackRange && distance > minDistance {
			minDistance = distance
			near = defender
		}
	}
	return near
}

func (g *GridCombat) getDistance(attacker Combatable, defender Combatable) int {
	aX, aY := g.index2Cord(g.comb2pos[attacker])
	dX, dY := g.index2Cord(g.comb2pos[defender])
	return util.Abs(aX-dX) + util.Abs(aY-dY)
}

func (g *GridCombat) index2Cord(idx int) (int, int) {
	return idx % GRID_WIDTH, idx / GRID_WIDTH
}

func (g *GridCombat) cord2Index(x, y int) int {
	return y*GRID_WIDTH + x
}

func (g *GridCombat) placeCombatables() {
	for i := 0; i < len(g.combatables); i++ {
		pos := util.GetRandomInt(GRID_WIDTH*GRID_HEIGHT*0.5 - 1)
		if g.combatables[i].GetCombatType() == ENEMY {
			pos += GRID_WIDTH * GRID_HEIGHT * 0.5
		}
		_, exist := g.pos2comb[pos]
		if exist {
			i--
			continue
		}
		g.pos2comb[pos] = g.combatables[i]
		g.comb2pos[g.combatables[i]] = pos
	}
}
