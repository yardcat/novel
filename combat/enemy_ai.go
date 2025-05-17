package combat

import (
	"container/list"
	"my_test/log"
	"strings"
)

const (
	ENEMY_BEHAVIOR_ATTACK = "attack"
	ENEMY_BEHAVIOR_DEFEND = "defend"
	ENEMY_BEHAVIOR_ESCAPE = "escape"
	ENEMY_BEHAVIOR_SKILL  = "skill"
)

type CardEnemyBehavior struct {
}

type EnemyAction struct {
	Action      string
	ActionValue int
	Target      int
}

type EnemyAI struct {
	enemies           []*CardEnemy
	history           map[*CardEnemy]*list.List
	currentTurnAction map[*CardEnemy]EnemyAction
}

func NewEnemyAI(enemy []*CardEnemy) *EnemyAI {
	ai := &EnemyAI{
		enemies:           enemy,
		history:           make(map[*CardEnemy]*list.List),
		currentTurnAction: make(map[*CardEnemy]EnemyAction),
	}
	for _, v := range enemy {
		ai.history[v] = list.New()
	}

	return ai
}

func (e *EnemyAI) SetAction(enemy *CardEnemy, rule string, value int) {
	action := EnemyAction{
		ActionValue: value,
		Action:      rule,
		Target:      0,
	}
	if strings.Contains(rule, "attack") {
		action.ActionValue = enemy.Values["attack"] + enemy.Strength
	}
	e.currentTurnAction[enemy] = action
	e.history[enemy].PushFront(action)
}

func (e *EnemyAI) GetAction(enemy *CardEnemy) EnemyAction {
	if e.history[enemy].Front() == nil {
		log.Error("no history")
		return EnemyAction{}
	}
	return e.history[enemy].Front().Value.(EnemyAction)
}

func (e *EnemyAI) onEnemyTurnFinish() {
	e.currentTurnAction = make(map[*CardEnemy]EnemyAction)
}
