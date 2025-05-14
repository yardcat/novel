package combat

import (
	"encoding/json"
)

const (
	CASTER_TYPE_CARD = iota
	CASTER_TYPE_POTION
	CASTER_TYPE_RELIC
	CASTER_TYPE_ENEMY
)

const (
	TIMING_NONE = iota
	TIMING_ACTOR_TURN_START
	TIMING_ACTOR_TURN_END
	TIMING_ENEMY_TURN_START
	TIMING_ENEMY_TURN_END
	TIMING_USE_CARD
	TIMING_DISCARD_CARD
	TIMING_COMBAT_START
	TIMING_COMBAT_END
	TIMING_ENTER_REST
	TIMING_ENTER_SHOP
	TIMING_ENTER_EVENT
	TIMING_IMMEDIATE
	TIMING_ADD_ARMOR
	TIMING_ADD_BUFF
	TIMING_ADD_DEBUFF
	TIMING_ATTACK
)

type Effect struct {
	CasterType int
	CasterID   string
	Timing     int
	Rule       string
}

func (e *Effect) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Timing string `json:"timing"`
		Rule   string `json:"rule"`
	}{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.Rule = tmp.Rule

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
	case "USE_CARD":
		e.Timing = TIMING_USE_CARD
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
	default:
		e.Timing = TIMING_NONE
	}
	return nil
}

func (t *Tower) EffectOn(timing int) {
	for _, v := range t.effects[timing] {
		t.UseEffect(v)
	}
}

func (t *Tower) UseEffect(effect *Effect) {
	t.engine.ExecuteSelectedRules(t.ruleBuilder, []string{effect.Rule})

	if t.currentCombat != nil {
		t.currentCombat.requestUpdateUI()
	}
}
