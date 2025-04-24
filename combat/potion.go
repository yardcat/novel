package combat

const (
	POTION_COMMON = iota
	POTION_UNCOMMON
	POTION_RARE
)

type PotionEffect struct {
	Type  int
	Value int
}

type Potion struct {
	Name        string
	Description string
	Rarity      int
	Effects     []PotionEffect
}

type PotionSystem struct {
	potions []*Potion
	maxSize int
}

func NewPotionSystem(maxSize int) *PotionSystem {
	return &PotionSystem{
		potions: make([]*Potion, 0),
		maxSize: maxSize,
	}
}

func (ps *PotionSystem) AddPotion(potion *Potion) bool {
	if len(ps.potions) >= ps.maxSize {
		return false
	}
	ps.potions = append(ps.potions, potion)
	return true
}

func (ps *PotionSystem) UsePotion(index int, cbt Combat) bool {
	if index < 0 || index >= len(ps.potions) {
		return false
	}

	potion := ps.potions[index]
	for _, effect := range potion.Effects {
		switch effect.Type {
		case POTION_EFFECT_HEAL:
			for _, actor := range cbt.Actors() {
				actor.GetBase().Life = min(actor.GetBase().Life+effect.Value, actor.GetBase().MaxLife)
			}
		case POTION_EFFECT_ENERGY:
			// cbt.Energy += effect.Value
		case POTION_EFFECT_STRENGTH:
			for _, actor := range cbt.Actors() {
				actor.GetBase().AddStatus(Status{
					Type: STATUS_STRENGTH,
					Turn: effect.Value,
				})
			}
		case POTION_EFFECT_DEXTERITY:
			for _, actor := range cbt.Actors() {
				actor.GetBase().AddStatus(Status{
					Type: STATUS_ARMOR,
					Turn: effect.Value,
				})
			}
		}
	}

	// Remove used potion
	ps.potions = append(ps.potions[:index], ps.potions[index+1:]...)
	return true
}

const (
	POTION_EFFECT_HEAL = iota
	POTION_EFFECT_ENERGY
	POTION_EFFECT_STRENGTH
	POTION_EFFECT_DEXTERITY
)

// Common potions
var (
	BloodPotion = &Potion{
		Name:        "Blood Potion",
		Description: "Heal 20% of max HP",
		Rarity:      POTION_COMMON,
		Effects: []PotionEffect{
			{Type: POTION_EFFECT_HEAL, Value: 20}, // Value is percentage
		},
	}

	EnergyPotion = &Potion{
		Name:        "Energy Potion",
		Description: "Gain 2 Energy",
		Rarity:      POTION_COMMON,
		Effects: []PotionEffect{
			{Type: POTION_EFFECT_ENERGY, Value: 2},
		},
	}

	StrengthPotion = &Potion{
		Name:        "Strength Potion",
		Description: "Gain 2 Strength",
		Rarity:      POTION_UNCOMMON,
		Effects: []PotionEffect{
			{Type: POTION_EFFECT_STRENGTH, Value: 2},
		},
	}

	DexterityPotion = &Potion{
		Name:        "Dexterity Potion",
		Description: "Gain 2 Dexterity",
		Rarity:      POTION_UNCOMMON,
		Effects: []PotionEffect{
			{Type: POTION_EFFECT_DEXTERITY, Value: 2},
		},
	}
)
