package combat

type Relic struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rarity      int       `json:"rarity"`
	Price       int       `json:"price"`
	Effects     []*Effect `json:"effects"`
}
