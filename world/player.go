package world

import (
	"encoding/json"
	"fmt"
	"my_test/combat"
)

type Player struct {
	Id       string
	Health   int    `json:"Health"`
	Hunger   int    `json:"Hunger"`
	Thirst   int    `json:"Thirst"`
	Energy   int    `json:"Energy"`
	Location string `json:"-"`
	Bag      *Bag   `json:"-"`
	Story    *Story `json:"-"`
	Pets     []Pet
}

func NewPlayer(story *Story, id string) *Player {
	return &Player{
		Id:       id,
		Health:   100,
		Hunger:   100,
		Thirst:   100,
		Energy:   100,
		Location: "beach",
		Bag:      NewBag(),
		Story:    story,
		Pets:     make([]Pet, 0),
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
	// day event
	maps["ChangeStatus"] = s.OnChangeStatus
	maps["Bonus"] = s.OnBonus

	// user event
	maps["Collect"] = s.Collect
}

func (p *Player) Collect(event CollectEvent) CollectEventReply {
	reply := CollectEventReply{}
	for _, i := range event.Items {
		if p.Energy >= 10 {
			p.Energy -= 10
			p.Bag.Add(p.Story.ItemSystem.GetItemByName(i.Item), i.Count)
			reply.EnergyCost += 10
			reply.Items = append(reply.Items, i)
		}
	}
	return reply
}

func (s *Player) GetCombatableBase() combat.CombatableBase {
	return combat.CombatableBase{
		Name:        "player",
		CombatType:  combat.ACTOR,
		Life:        s.Health,
		Attack:      10,
		Defense:     2,
		Dodge:       10,
		AttackSpeed: 10,
		AttackRange: 6,
		AttackStep:  0,
	}
}

func (s *Player) AddPet(name string) {
	proto := s.Story.PetSystem.GetPet(name)
	pet := CreatePet(&proto)
	s.Pets = append(s.Pets, *pet)
}

func (s *Player) OnChangeStatus(event ChangeStatusEvent) {
	switch event.Type {
	case "hp":
		s.Health += event.Value
	}
}

func (s *Player) OnBonus(event BonusEvent) {
	s.Bag.Add(s.Story.ItemSystem.GetItemByName(event.Item), event.Count)
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
