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

func (p *Potion) CanUse() bool {
	return false
}

func (ps *Potion) Use() {
}
