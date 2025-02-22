package world

import "strconv"

type Player struct {
	Health    int
	Hunger    int
	Thirst    int
	Energy    int
	Inventory []string
	Location  string
	Bag       *Bag
	Story     *Story
}

func NewPlayer(story *Story) *Player {
	return &Player{
		Health:    100,
		Hunger:    100,
		Thirst:    100,
		Energy:    100,
		Inventory: make([]string, 0),
		Location:  "beach",
		Bag:       NewBag(),
		Story:     story,
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

func (s *Player) RegisterEventHander(maps map[string]any) {
	maps["ChangeStatus"] = s.OnChangeStatus
	maps["Bonus"] = s.OnBonus
}

func (s *Player) OnChangeStatus(typ string, value string) {
	switch typ {
	case "hp":
		v, _ := strconv.Atoi(value)
		s.Health += v
	}
}

func (s *Player) OnBonus(stuff string, count string) {
	n, _ := strconv.Atoi(count)
	itemId := s.Story.itemSystem.GetItemId(stuff)
	s.Bag.Add(itemId, n)
}
