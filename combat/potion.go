package combat

const (
	POTION_COMMON = iota
	POTION_UNCOMMON
	POTION_RARE
)

type Potion struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rarity      int       `json:"rarity"`
	Price       int       `json:"price"`
	Effects     []*Effect `json:"effects"`
}

func (p *Potion) CanUse() bool {
	return false
}
