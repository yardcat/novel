package island

import "my_test/world"

type Player struct {
	Health    int
	Hunger    int
	Thirst    int
	Energy    int
	Inventory []string
	Location  string
	Bag       *world.Bag
}

func NewPlayer() *Player {
	return &Player{
		Health:    100,
		Hunger:    100,
		Thirst:    100,
		Energy:    100,
		Inventory: make([]string, 0),
		Location:  "beach",
		Bag:       world.NewBag(),
	}
}

func (p *Player) CollectStuff(stuff string) {
	if p.Energy >= 10 {
		p.Energy -= 10
		p.Inventory = append(p.Inventory, stuff)
	}
}

func (p *Player) Update() {
	// Decrease stats over time
	p.Hunger -= 1
	p.Thirst -= 2
	p.Energy += 1

	// Cap stats at min/max values
	if p.Hunger < 0 {
		p.Hunger = 0
	}
	if p.Thirst < 0 {
		p.Thirst = 0
	}
	if p.Energy > 100 {
		p.Energy = 100
	}

	// Take damage if hunger/thirst are depleted
	if p.Hunger == 0 || p.Thirst == 0 {
		p.Health -= 5
	}
}

func (p *Player) GetInfoAsJSON() string {
	return ""
}
