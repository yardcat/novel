package world

import (
	"my_test/combat"
	"my_test/log"
	"my_test/util"
)

const (
	MINE_ACTION_NONE     = -1
	MINE_ACTION_RESOURCE = 0
	MINE_ACTION_ENEMY    = 1
	RATE_MAX             = 100
)

type MineResource struct {
	Store int `json:"store"`
	Rate  int `json:"rate"`
}

type MineEnemy struct {
	Rate int `json:"rate"`
}

type Mine struct {
	Resource   map[string]MineResource `json:"resource"`
	Enemy      map[string]MineEnemy    `json:"enemy"`
	ActionRate []int                   `json:"action_rate"`
}

func NewMine() *Mine {
	return &Mine{}
}

func (m *Mine) Explore() ExploreResult {
	action := m.selectAction()
	if action == MINE_ACTION_RESOURCE {
		resource := m.selectResource()
		if resource == "" {
			log.Info("no resource")
		} else {
			log.Info("get resource %s", resource)
			return ExploreResult{Type: EXPLORE_RESULT_ITEM, itemResult: []string{resource}}
		}
	} else if action == MINE_ACTION_ENEMY {
		enemy := m.selectEnemy()
		if enemy == "" {
			log.Info("no enemy")
		} else {
			log.Info("get enemy %s", enemy)
			return ExploreResult{Type: EXPLORE_RESULT_COMBAT, combatResult: &combat.CombatResult{}}
		}
	} else {
		log.Info("action do nothing")
	}
	return ExploreResult{Type: EXPLORE_RESULT_NONE}
}

func (m *Mine) PassBy() {
	log.Info("pass by unimplemented")
}

func (m *Mine) selectAction() int {
	rate := util.GetRandomInt(RATE_MAX)
	for i, r := range m.ActionRate {
		rate -= r
		if rate < 0 {
			return i
		}
	}
	return MINE_ACTION_NONE
}

func (m *Mine) selectResource() string {
	rate := util.GetRandomInt(RATE_MAX)
	for k, v := range m.Resource {
		rate -= v.Rate
		if rate < 0 {
			return k
		}
	}
	return ""
}

func (m *Mine) selectEnemy() string {
	rate := util.GetRandomInt(100)
	for k, v := range m.Enemy {
		rate -= v.Rate
		if rate < 0 {
			return k
		}
	}
	return ""
}
