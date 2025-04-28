package combat

const (
	CARD_TYPE_ATTACK = iota
	CARD_TYPE_SKILL
	CARD_TYPE_EFFECT
)

const (
	CARD_RARITY_COMMON = iota
	CARD_RARITY_UNCOMMON
	CARD_RARITY_RARE
)

const (
	EFFECT_DAMAGE = iota
	EFFECT_VULNERABLE
	EFFECT_DEFEND
	EFFECT_ADD_CARD
	EFFECT_MULTI_DAMAGE
	EFFECT_DAMAGE_DEFENSE
	EFFECT_WEAK
	EFFECT_STRENGTH
	EFFECT_HEAL
	EFFECT_MAX_HEALTH
)

const (
	UI_ACTOR_HP = iota
	UI_ACTOR_MAX_HP
	UI_ENEMY_HP
	UI_ENEMY_MAX_HP
)

const (
	STATUS_VULNERABLE = iota
	STATUS_WEAK
	STATUS_STRENGTH
	STATUS_ARMOR
	STATUS_POISON
)

const (
	CARD_INIT   = 5
	CARD_MAX    = 10
	ENERGY_INIT = 3
	ENERGY_MAX  = 20
)

const (
	TIMING_START = iota
	TIMING_END
)

const (
	SKILL_ADD_STRENGTH = iota
)

type CardSkill struct {
	Effects      []*CardEffect
	Timing       int
	TurnInterval int
	TurnCount    int
}

type CardEffect struct {
	Effect string `json:"effect"`
	Value  any    `json:"value"`
}

type Card struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Type        int          `json:"type"`
	Rarity      int          `json:"rarity"`
	Cost        int          `json:"cost"`
	Upgrade     []*Card      `json:"upgrade,omitempty"`
	Effects     []CardEffect `json:"effects"`
}

type CardCareer struct {
	InitCards []*Card
}

type CardEvent struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Effects     []CardEffect `json:"effects"`
}

type CardTurnInfo struct {
	Cards        []string
	DrawCount    int
	DiscardCount int
	RemoveCount  int
	Energy       int
}
