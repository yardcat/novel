package world

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Player struct {
	Health    int      `json:"Health"`
	Hunger    int      `json:"Hunger"`
	Thirst    int      `json:"Thirst"`
	Energy    int      `json:"Energy"`
	Inventory []string `json:"-"`
	Location  string   `json:"-"`
	Bag       *Bag     `json:"-"`
	Story     *Story   `json:"-"`
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

func (s *Player) RegisterEventHander(maps map[string]any) {
	maps["ChangeStatus"] = s.OnChangeStatus
	maps["Bonus"] = s.OnBonus
}

func (s *Player) OnChangeStatus(params map[string]string) {
	typ := params["type"]
	value := params["value"]
	switch typ {
	case "hp":
		v, _ := strconv.Atoi(value)
		s.Health += v
	}
}

func (s *Player) OnBonus(params map[string]string) {
	item := params["item"]
	count := params["count"]
	n, _ := strconv.Atoi(count)
	itemId := s.Story.itemSystem.GetItemId(item)
	s.Bag.Add(itemId, n)
}

func (s *Player) ToString() string {
	return fmt.Sprintf(`Health: %d, Hunger: %d, Thirst: %d, Energy: %d`,
		s.Health, s.Hunger, s.Thirst, s.Energy)
}

func (s *Player) ToJson() string {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(data)
}
