package combat

const (
	RELIC_STARTER = iota
	RELIC_COMMON
	RELIC_UNCOMMON
	RELIC_RARE
	RELIC_BOSS
	RELIC_SHOP
)

type RelicEffect struct {
	Type  int
	Value int
}

type Relic struct {
	Name        string
	Description string
	Rarity      int
	Effects     []RelicEffect
}

type RelicSystem struct {
	relics []*Relic
}

func NewRelicSystem() *RelicSystem {
	return &RelicSystem{
		relics: make([]*Relic, 0),
	}
}

func (rs *RelicSystem) AddRelic(relic *Relic) {
	rs.relics = append(rs.relics, relic)
}

func (rs *RelicSystem) OnTurnStart(cbt *Combat) {
	for _, relic := range rs.relics {
		for _, effect := range relic.Effects {
			switch effect.Type {
			case RELIC_EFFECT_ENERGY:
				// cbt.energy += effect.Value
			case RELIC_EFFECT_DRAW:
				// for i := 0; i < effect.Value; i++ {
				// 	cbt.cardSystem.DrawCard()
				// }
			case RELIC_EFFECT_STRENGTH:
				// for _, actor := range cbt.actors {
				// 	actor.GetBase().Status.AddStatus(Status{
				// 		Type:  STATUS_STRENGTH,
				// 		Value: effect.Value,
				// 		Turn:  1,
				// 	})
				// }
			}
		}
	}
}

func (rs *RelicSystem) OnCardPlayed(card *Card, combatSystem *Combat) {
	for _, relic := range rs.relics {
		for _, effect := range relic.Effects {
			switch effect.Type {
			case RELIC_EFFECT_BLOCK:
				// for _, actor := range combatSystem.actors {
				// 	actor.GetBase().Status.AddStatus(Status{
				// 		Type:  STATUS_DEFENSE,
				// 		Value: effect.Value,
				// 		Turn:  1,
				// 	})
				// }
			}
		}
	}
}

func (rs *RelicSystem) OnEnemyDamaged(enemy Combatable, damage int) {
	for _, relic := range rs.relics {
		for _, effect := range relic.Effects {
			switch effect.Type {
			case RELIC_EFFECT_HEAL:
				// TODO: Implement healing
			}
		}
	}
}

const (
	RELIC_EFFECT_ENERGY = iota
	RELIC_EFFECT_DRAW
	RELIC_EFFECT_STRENGTH
	RELIC_EFFECT_BLOCK
	RELIC_EFFECT_HEAL
)

// Common relics
var (
	BurningBlood = &Relic{
		Name:        "Burning Blood",
		Description: "Heal 6 HP at the end of combat",
		Rarity:      RELIC_BOSS,
		Effects: []RelicEffect{
			{Type: RELIC_EFFECT_HEAL, Value: 6},
		},
	}

	RingOfTheSnake = &Relic{
		Name:        "Ring of the Snake",
		Description: "Draw 2 additional cards at the start of each combat",
		Rarity:      RELIC_STARTER,
		Effects: []RelicEffect{
			{Type: RELIC_EFFECT_DRAW, Value: 2},
		},
	}

	Vajra = &Relic{
		Name:        "Vajra",
		Description: "Start each combat with 1 Strength",
		Rarity:      RELIC_COMMON,
		Effects: []RelicEffect{
			{Type: RELIC_EFFECT_STRENGTH, Value: 1},
		},
	}
)
