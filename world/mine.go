package world

import (
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

func (m *Mine) Explore(times int) {
	for i := 0; i < times; i++ {
		action := m.SelectAction()
		if action == MINE_ACTION_RESOURCE {
			resource := m.SelectResource()
			if resource == "" {
				log.Info("no resource")
			} else {
				log.Info("get resource %s", resource)
			}
		} else if action == MINE_ACTION_ENEMY {
			enemy := m.SelectEnemy()
			if enemy == "" {
				log.Info("no enemy")
			} else {
				log.Info("get enemy %s", enemy)
			}
		} else {
			log.Info("action do nothing")
		}
	}
}

func (m *Mine) SelectAction() int {
	rate := util.GetRandomInt(RATE_MAX)
	for i, r := range m.ActionRate {
		rate -= r
		if rate < 0 {
			return i
		}
	}
	return MINE_ACTION_NONE
}

func (m *Mine) SelectResource() string {
	rate := util.GetRandomInt(RATE_MAX)
	for k, v := range m.Resource {
		rate -= v.Rate
		if rate < 0 {
			return k
		}
	}
	return ""
}

func (m *Mine) SelectEnemy() string {
	rate := util.GetRandomInt(100)
	for k, v := range m.Enemy {
		rate -= v.Rate
		if rate < 0 {
			return k
		}
	}
	return ""
}
