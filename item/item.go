package item

const (
	Weapon = iota
	Armor
	Consumable
)

type Item struct {
	Name string
	Type int
}

type WeaponItem struct {
	Item
	Attack int
}

type ArmorItem struct {
	Item
	Defence int
	Dodge   int
}

const (
	Life = iota
	Maigc
)

type PotionItem struct {
	Item
	Effect int
	Value  int
}
