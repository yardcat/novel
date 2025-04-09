package combat

import "my_test/log"

const (
	ENEMY_BEHAVIOR_ATTACK = iota
	ENEMY_BEHAVIOR_DEFEND
	ENEMY_BEHAVIOR_BUFF
	ENEMY_BEHAVIOR_DEBUFF
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
	enemy     *Enemy
	behaviors []EnemyBehavior
	combat    Combat
}

func NewEnemyAI(enemy *Enemy, combat Combat) *EnemyAI {
	return &EnemyAI{
		enemy:     enemy,
		behaviors: make([]EnemyBehavior, 0),
		combat:    combat,
	}
}

func (ai *EnemyAI) AddBehavior(behavior EnemyBehavior) {
	ai.behaviors = append(ai.behaviors, behavior)
}

func (ai *EnemyAI) ChooseAction() EnemyIntent {
	var intent EnemyIntent

	if ai.enemy.GetBase().Life < ai.enemy.GetBase().MaxLife/2 {
		for _, behavior := range ai.behaviors {
			if behavior.Type == ENEMY_BEHAVIOR_DEFEND {
				intent.Behavior = behavior
				intent.Target = ai.enemy
				return intent
			}
		}
	}

	if ai.enemy.GetBase().GetStatusValue(STATUS_STRENGTH) > 0 {
		for _, behavior := range ai.behaviors {
			if behavior.Type == ENEMY_BEHAVIOR_ATTACK {
				intent.Behavior = behavior
				if len(ai.combat.Actors()) > 0 {
					intent.Target = ai.combat.Actors()[0]
				}
				return intent
			}
		}
	}

	for _, behavior := range ai.behaviors {
		intent.Behavior = behavior
		if behavior.Type == ENEMY_BEHAVIOR_ATTACK && len(ai.combat.Actors()) > 0 {
			intent.Target = ai.combat.Actors()[0]
		} else {
			intent.Target = ai.enemy
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
	case ENEMY_BEHAVIOR_BUFF:
		log.Info("buff")
	case ENEMY_BEHAVIOR_DEBUFF:
		log.Info("debuff")
	}
}
