package world

import (
	"encoding/json"
	"fmt"
	"my_test/career"
	"my_test/combat"
	"my_test/equip"
	"my_test/event"
	"my_test/log"
	"os"
)

type Player struct {
	PlayerAttr
	Id             string
	Bag            *Bag
	Story          *Story
	Pets           []Pet
	Npcs           []Npc
	Equips         PlayerEquips
	Career         *career.Career
	attr           PlayerAttr
	attrAdd        PlayerAttr
	attrAddPercent PlayerAttr
	needUpdateAttr bool
}

func NewPlayer(story *Story, id string) *Player {
	p := &Player{
		Id:             id,
		Bag:            NewBag(),
		Story:          story,
		Pets:           make([]Pet, 0),
		Npcs:           make([]Npc, 0),
		Equips:         PlayerEquips{},
		needUpdateAttr: true,
	}
	p.loadPlayerData()
	return p
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

func (p *Player) Collect(ev event.CollectEvent) event.CollectEventReply {
	reply := event.CollectEventReply{}
	for _, i := range ev.Items {
		if p.Energy >= 10 {
			p.Energy -= 10
			p.Bag.Add(p.Story.itemSystem.GetItemByName(i.Item), i.Count)
			reply.EnergyCost += 10
			reply.Items = append(reply.Items, i)
		}
	}
	return reply
}

func (p *Player) UpdateAttr() {
	p.caculateAttrFromEquip()
}

func (p *Player) caculateAttrFromEquip() {
	attrs := make([]equip.Attr, 0)
	for _, v := range p.Equips.Equips {
		attrs = append(attrs, v.GetAttrs()...)
	}
	p.attrAdd.UpdateFromAttr(attrs, false)
	p.attrAddPercent.UpdateFromAttr(attrs, true)
	p.UpdateFinal(&p.attrAdd, &p.attrAddPercent)
}

func (p *Player) GetCombatableBase() combat.CombatableBase {
	// TODO : 优化获取属性逻辑, 使用Enum或者value
	if p.needUpdateAttr {
		p.UpdateAttr()
	}
	return p.attr.CombatableBase
}

func (p *Player) AddCareer(name string) {
	p.Career = p.Story.careerSystem.GetCareer(name)
}

func (p *Player) AddPet(name string) {
	proto := p.Story.petSystem.GetPet(name)
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
	proto := p.Story.npcSystem.GetNpc(name)
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

func (p *Player) OnChangeStatus(event event.ChangeStatusEvent) {
	switch event.Type {
	case "hp":
		p.Health += event.Value
	}
}

func (p *Player) OnBonus(event event.BonusEvent) {
	p.Bag.Add(p.Story.itemSystem.GetItemByName(event.Item), event.Count)
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

func (p *Player) loadPlayerData() {
	p.LoadAttrFromJson()
}

func (p *Player) LoadAttrFromJson() {
	data, err := os.ReadFile(p.Story.GetResources().GetPath("player/player.json"))
	if err != nil {
		log.Error("Failed to read player.json: %v", err)
		return
	}

	err = json.Unmarshal(data, &p.attr)
	if err != nil {
		log.Error("Failed to unmarshal player.json: %v", err)
		return
	}
}
