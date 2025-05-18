package combat

import (
	"container/list"
	"my_test/log"
)

type CardEnemyBehavior struct {
}

type EnemyAction struct {
	Action      string
	ActionValue string
	Description string
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

func (e *EnemyAI) SetAction(enemy *CardEnemy, rule string, value string) {
	action := EnemyAction{
		ActionValue: value,
		Action:      rule,
		Target:      0,
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
