package combat

import (
	"encoding/json"
	"my_test/log"
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
	TIMING_UPGRADE_CARD
	TIMING_DISCARD_CARD
	TIMING_COMBAT_START
	TIMING_COMBAT_END
	TIMING_ENTER_REST
	TIMING_ENTER_SHOP
	TIMING_ENTER_EVENT
	TIMING_IMMEDIATE
	TIMING_ADD_ARMOR
	TIMING_ADD_DEBUFF
	TIMING_ATTACK
	TIMING_ENEMY_HURT
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
	e.Timing = TimingStr2Int(tmp.Timing)
	return nil
}

func TimingStr2Int(timing string) int {
	ret := TIMING_NONE
	switch timing {
	case "ACTOR_TURN_START":
		ret = TIMING_ACTOR_TURN_START
	case "ACTOR_TURN_END":
		ret = TIMING_ACTOR_TURN_END
	case "ENEMY_TURN_START":
		ret = TIMING_ENEMY_TURN_START
	case "ENEMY_TURN_END":
		ret = TIMING_ENEMY_TURN_END
	case "USE_CARD":
		ret = TIMING_USE_CARD
	case "DISCARD_CARD":
		ret = TIMING_DISCARD_CARD
	case "COMBAT_START":
		ret = TIMING_COMBAT_START
	case "COMBAT_END":
		ret = TIMING_COMBAT_END
	case "ENTER_REST":
		ret = TIMING_ENTER_REST
	case "ENTER_SHOP":
		ret = TIMING_ENTER_SHOP
	case "ENTER_EVENT":
		ret = TIMING_ENTER_EVENT
	case "TIMING_ENEMY_HURT":
		ret = TIMING_ENEMY_HURT
	default:
		ret = TIMING_NONE
	}
	return ret
}

func (t *Tower) EffectOn(timing int) {
	for _, v := range t.effects[timing] {
		t.UseEffect(v)
	}
}

func (t *Tower) UseEffect(effect *Effect) {
	err := t.engine.ExecuteSelectedRules(t.ruleBuilder, []string{effect.Rule})
	if err != nil {
		log.Error("use effect %s err %v", effect.Rule, err)
		panic(err)
	}

	if t.currentCombat != nil {
		t.currentCombat.requestUpdateUI()
	}
}
