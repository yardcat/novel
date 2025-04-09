package combat

import (
	"my_test/log"
	"my_test/util"
)

const (
	GRID_WIDTH  = 6
	GRID_HEIGHT = 6
)

type GridLayout struct {
	combat   Combat
	pos2comb map[int]Combatable
	comb2pos map[Combatable]int
}

func NewGridLayout(combat Combat) *GridLayout {
	g := &GridLayout{
		combat,
		make(map[int]Combatable),
		make(map[Combatable]int),
	}
	g.placeCombatables()
	return g
}

func (g *GridLayout) getComb(x, y int) Combatable {
	idx := y*GRID_WIDTH + x
	return g.pos2comb[idx]
}

func (g *GridLayout) ChooseDefender(attacker Combatable) Combatable {
	var near Combatable
	if attacker.GetCombatType() == ACTOR {
		near = g.getNearDefender(attacker, g.combat.Enemies())
	} else if attacker.GetCombatType() == ENEMY {
		near = g.getNearDefender(attacker, g.combat.Enemies())
	} else {
		log.Info("unknown attacker type %d", attacker.GetCombatType())
	}
	return near
}

func (g *GridLayout) getNearDefender(attacker Combatable, defenders []Combatable) Combatable {
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

func (g *GridLayout) getDistance(attacker Combatable, defender Combatable) int {
	aX, aY := g.index2Cord(g.comb2pos[attacker])
	dX, dY := g.index2Cord(g.comb2pos[defender])
	return util.Abs(aX-dX) + util.Abs(aY-dY)
}

func (g *GridLayout) index2Cord(idx int) (int, int) {
	return idx % GRID_WIDTH, idx / GRID_WIDTH
}

func (g *GridLayout) cord2Index(x, y int) int {
	return y*GRID_WIDTH + x
}

func (g *GridLayout) placeCombatables() {
	for i := 0; i < len(g.combat.Combatables()); i++ {
		pos := util.GetRandomInt(GRID_WIDTH*GRID_HEIGHT*0.5 - 1)
		combatables := g.combat.Combatables()
		if combatables[i].GetCombatType() == ENEMY {
			pos += GRID_WIDTH * GRID_HEIGHT * 0.5
		}
		_, exist := g.pos2comb[pos]
		if exist {
			i--
			continue
		}
		g.pos2comb[pos] = combatables[i]
		g.comb2pos[combatables[i]] = pos
	}
}
