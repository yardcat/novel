package combat

import (
	"container/list"
	"my_test/log"
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
	Description string
	Target      Combatable
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

func (e *EnemyAI) PrepareAction(enemies []*CardEnemy, actors []*CardActor) {
	var action EnemyAction

	for _, enemy := range enemies {
		if enemy.Life > enemy.MaxLife/2 {
			action.Action = ENEMY_BEHAVIOR_ATTACK
			action.Target = actors[0]
			action.ActionValue = enemy.Attack + enemy.Strength
		} else if enemy.Life < enemy.MaxLife/2 {
			action.Action = ENEMY_BEHAVIOR_DEFEND
			action.Target = actors[0]
			action.ActionValue = 5
		}
		e.currentTurnAction[enemy] = action
		e.history[enemy].PushFront(action)
	}
}

func (e *EnemyAI) EnemyAction(enemy *CardEnemy) EnemyAction {
	if e.history[enemy].Front() == nil {
		log.Error("no history")
	}
	return e.history[enemy].Front().Value.(EnemyAction)
}

func (e *EnemyAI) EnemyActions() []EnemyAction {
	actions := []EnemyAction{}
	for _, enemy := range e.enemies {
		actions = append(actions, e.EnemyAction(enemy))
	}
	return actions
}

func (e *EnemyAI) onEnemyTurnFinish() {
	e.currentTurnAction = make(map[*CardEnemy]EnemyAction)
}
