package combat

const (
	ROOM_TYPE_FIGHT = iota
	ROOM_TYPE_SHOP
	ROOM_TYPE_EVENT
)

type Room interface {
	Type() int
}

type FightRoom struct {
	Enemy []*Enemy
	Bouns []string
}

func (f *FightRoom) Type() int {
	return ROOM_TYPE_FIGHT
}

type ShopRoom struct {
	Potions []Potion
	Relics  []Relic
}

func (f *ShopRoom) Type() int {
	return ROOM_TYPE_SHOP
}

type EventRoom struct {
	Event string
}

func (f *EventRoom) Type() int {
	return ROOM_TYPE_EVENT
}

type Tower struct {
	RoomMap []*Room
}
