package world

import (
	"encoding/json"
	"fmt"
	"my_test/combat"
	"my_test/log"
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
	Npcs     []Npc
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
		Npcs:     make([]Npc, 0),
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

func (p *Player) RegisterEventHander(maps map[string]any) {
	// day event
	maps["ChangeStatus"] = p.OnChangeStatus
	maps["Bonus"] = p.OnBonus

	// user event
	maps["Collect"] = p.Collect
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

func (p *Player) GetCombatableBase() combat.CombatableBase {
	return combat.CombatableBase{
		Name:        "player",
		CombatType:  combat.ACTOR,
		Life:        p.Health,
		Attack:      10,
		Defense:     2,
		Dodge:       10,
		AttackSpeed: 10,
		AttackRange: 6,
		AttackStep:  0,
	}
}

func (p *Player) AddPet(name string) {
	proto := p.Story.PetSystem.GetPet(name)
	pet := CreatePet(&proto)
	p.Pets = append(p.Pets, *pet)
}

func (p *Player) GetPet(name string) *Pet {
	for _, pet := range p.Pets {
		if pet.Name == name {
			return &pet
		}
	}
	return nil
}

func (p *Player) AddNpc(name string) {
	proto := p.Story.NpcSystem.GetNpc(name)
	npc := CreateNpc(&proto)
	p.Npcs = append(p.Npcs, *npc)
}

func (p *Player) GetNpc(name string) *Npc {
	for _, npc := range p.Npcs {
		if npc.Name == name {
			return &npc
		}
	}
	return nil
}

func (p *Player) OnChangeStatus(event ChangeStatusEvent) {
	switch event.Type {
	case "hp":
		p.Health += event.Value
	}
}

func (p *Player) OnBonus(event BonusEvent) {
	p.Bag.Add(p.Story.ItemSystem.GetItemByName(event.Item), event.Count)
}

func (p *Player) ToString() string {
	return fmt.Sprintf(`Health: %d, Hunger: %d, Thirst: %d, Energy: %d`,
		p.Health, p.Hunger, p.Thirst, p.Energy)
}

func (p *Player) ToJson() string {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(data)
}

func (p *Player) OnCombatDone(result combat.CombatResult) {
	log.Info("player combat done %v", result)
}
