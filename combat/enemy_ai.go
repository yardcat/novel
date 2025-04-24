package combat

import "my_test/log"

const (
	ENEMY_BEHAVIOR_ATTACK = iota
	ENEMY_BEHAVIOR_DEFEND
	ENEMY_BEHAVIOR_ESCAPE
	ENEMY_BEHAVIOR_SKILL
	ENEMY_BEHAVIOR_COUNT
)

type EnemyBehavior struct {
	Type        int
	Value       int
	Description string
}

type EnemyIntent struct {
	Behavior EnemyBehavior
	Target   Combatable
}

type EnemyAI struct {
	behaviors []EnemyBehavior
}

func NewEnemyAI() *EnemyAI {
	return &EnemyAI{
		behaviors: make([]EnemyBehavior, 0),
	}
}

func (ai *EnemyAI) EnemyAction(enemy *Enemy, actors []*Actor) EnemyIntent {
	var intent EnemyIntent

	if enemy.GetBase().Life < enemy.GetBase().MaxLife/2 {
		for _, behavior := range ai.behaviors {
			if behavior.Type == ENEMY_BEHAVIOR_DEFEND {
				intent.Behavior = behavior
				intent.Target = enemy
				return intent
			}
		}
	}

	if enemy.GetBase().GetStatusValue(STATUS_STRENGTH) > 0 {
		for _, behavior := range ai.behaviors {
			if behavior.Type == ENEMY_BEHAVIOR_ATTACK {
				intent.Behavior = behavior
				if len(actors) > 0 {
					intent.Target = actors[0]
				}
				return intent
			}
		}
	}

	for _, behavior := range ai.behaviors {
		intent.Behavior = behavior
		if behavior.Type == ENEMY_BEHAVIOR_ATTACK && len(actors) > 0 {
			intent.Target = actors[0]
		} else {
			intent.Target = enemy
		}
		return intent
	}

	return intent
}

func (ai *EnemyAI) ExecuteAction(intent EnemyIntent) {
	switch intent.Behavior.Type {
	case ENEMY_BEHAVIOR_ATTACK:
		log.Info("ai attack")
	case ENEMY_BEHAVIOR_DEFEND:
		log.Info("ai defend")
	}
}
