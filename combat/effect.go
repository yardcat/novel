package combat

import (
	"my_test/log"

	"github.com/Knetic/govaluate"
)

const (
	CASTER_TYPE_CARD = iota
	CASTER_TYPE_POTION
	CASTER_TYPE_RELIC
	CASTER_TYPE_ENEMY
)

const (
	TIMING_ACTOR_TURN_START = iota
	TIMING_ACTOR_TURN_END
	TIMING_ENEMY_TURN_START
	TIMING_ENEMY_TURN_END
	TIMING_PLAY_CARD
	TIMING_DISCARD_CARD
)

type Effect struct {
	CasterType int
	CasterID   string
	Timing     int
	Condition  string
	Modifier   string
}

func (t *Tower) EffectOn(timing int) {
	for _, v := range t.effects[timing] {
		ok := false
		if v.Condition != "" {
			exp, err := govaluate.NewEvaluableExpression(v.Condition)
			if err != nil {
				log.Error("evaluate condition err %v", err)
				continue
			}
			ret, err := exp.Evaluate(nil)
			if err != nil {
				log.Error("evaluate condition err %v", err)
				continue
			}
			ok = ret.(bool)
		}
		if ok {

		}
	}
}
