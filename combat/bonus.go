package combat

const (
	BONUS_TYPE_CARD = iota
	BONUS_TYPE_POTION
	BONUS_TYPE_RELIC
)

type CombatBonus struct {
	Cards             []string
	CardChooseCount   int
	Potions           []string
	PotionChooseCount int
	Relics            []string
	RelicChooseCount  int
}
