package combat

import (
	"encoding/json"
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
	TIMING_COMBAT_START
	TIMING_COMBAT_END
	TIMING_ENTER_REST
	TIMING_ENTER_SHOP
	TIMING_ENTER_EVENT
	TIMING_IMMEDIATE
)

type Effect struct {
	CasterType int
	CasterID   string
	Timing     int
	Condition  string
	Modifier   string
}

func (e *Effect) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Condition string
		Modifier  string
		Timing    string
	}{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.Condition = tmp.Condition
	e.Modifier = tmp.Modifier

	// TODO: 通过editor生成json, 使用init而非string
	switch tmp.Timing {
	case "ACTOR_TURN_START":
		e.Timing = TIMING_ACTOR_TURN_START
	case "ACTOR_TURN_END":
		e.Timing = TIMING_ACTOR_TURN_END
	case "ENEMY_TURN_START":
		e.Timing = TIMING_ENEMY_TURN_START
	case "ENEMY_TURN_END":
		e.Timing = TIMING_ENEMY_TURN_END
	case "PLAY_CARD":
		e.Timing = TIMING_PLAY_CARD
	case "DISCARD_CARD":
		e.Timing = TIMING_DISCARD_CARD
	case "COMBAT_START":
		e.Timing = TIMING_COMBAT_START
	case "COMBAT_END":
		e.Timing = TIMING_COMBAT_END
	case "ENTER_REST":
		e.Timing = TIMING_ENTER_REST
	case "ENTER_SHOP":
		e.Timing = TIMING_ENTER_SHOP
	case "ENTER_EVENT":
		e.Timing = TIMING_ENTER_EVENT
	}
	return nil
}

func (t *Tower) EffectOn(timing int) {
	for _, v := range t.effects[timing] {
		t.UseEffect(v)
	}
}

func (t *Tower) UseEffect(effect *Effect) {
	ok := true
	if effect.Condition != "" {
		expr, err := govaluate.NewEvaluableExpression(effect.Condition)
		if err != nil {
			log.Error("evaluate condition err %v", err)
			return
		}
		ret, err := expr.Eval(t.export)
		if err != nil {
			log.Error("evaluate condition err %v", err)
			return
		}
		ok = ret.(bool)
	}
	if ok {
		expr, err := govaluate.NewEvaluableExpression(effect.Modifier)
		parameters := make(map[string]interface{})
		parameters["c"] = &t.export
		if err != nil {
			log.Error("evaluate modifier err %v", err)
			return
		}
		ret, err := expr.Evaluate(parameters)
		if err != nil {
			log.Error("evaluate modifier err %v", err)
			return
		}
		if ret.(bool) {
			log.Error("modifier run fail")
		}
	}
	if t.currentCombat != nil {
		t.currentCombat.requestUpdateUI()
	}
}
