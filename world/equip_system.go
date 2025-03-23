package world

import "my_test/equip"

type PlayerEquips struct {
	Armor     equip.Armor
	Bracelets [2]equip.Bracelet
	Gloves    [2]equip.Glove
	Helmet    equip.Helmet
	Necklace  equip.Necklace
	Rings     [2]equip.Ring
	Shoes     [2]equip.Shoe
	Trouser   equip.Trouser
	Weapon    [2]equip.Weapon
}

type EquipSystem struct {
}

func NewEquipSystem() *EquipSystem {
	return &EquipSystem{}
}
