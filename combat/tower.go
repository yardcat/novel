package combat

import (
	"encoding/json"
	"os"
	"path"
)

const (
	ROOM_TYPE_FIGHT = iota
	ROOM_TYPE_SHOP
	ROOM_TYPE_EVENT
	ROOM_TYPE_REST
)

type Room interface {
	Type() int
	SubRooms() []Room
	AddSubRoom(room Room)
}

type RoomBase struct {
	subRooms []Room
}

func (r *RoomBase) SubRooms() []Room {
	return r.subRooms
}

func (r *RoomBase) AddSubRoom(room Room) {
	r.subRooms = append(r.subRooms, room)
}

type FightRoom struct {
	RoomBase
	Enemy []*Enemy
	Bouns []string
}

func (f *FightRoom) Type() int {
	return ROOM_TYPE_FIGHT
}

type ShopRoom struct {
	RoomBase
	Potions []Potion
	Relics  []Relic
}

func (f *ShopRoom) Type() int {
	return ROOM_TYPE_SHOP
}

type EventRoom struct {
	RoomBase
	Event string
}

func (f *EventRoom) Type() int {
	return ROOM_TYPE_EVENT
}

type Floor struct {
	Rooms []Room
}

func (f *Floor) GetRooms() []int {
	ret := make([]int, 0)
	for _, v := range f.Rooms {
		ret = append(ret, v.Type())
	}
	return ret
}

type Tower struct {
	Floors     []*Floor
	FloorNum   int `json:"floor_num"`
	RoomNum    int `json:"room_num"`
	ShopNum    int `json:"shop_num"`
	RestNum    int `json:"rest_num"`
	EventNum   int `json:"event_num"`
	shopCount  int
	restCount  int
	eventCount int
}

func NewTower(path PathProvider) *Tower {
	t := &Tower{}
	t.loadData(path.GetPath("card"))
	t.Floors = make([]*Floor, 0, t.FloorNum)

	t.generateFloor()
	return t
}

func (t *Tower) EnterNextFloor() *Floor {
}

func (t *Tower) generateFloor() {
	fl := &Floor{}
	t.Floors = append(t.Floors, fl)
	topRoom := t.generateRoom()
	fl.Rooms = append(fl.Rooms, topRoom)
}

func (t *Tower) generateRoom() Room {
	return &FightRoom{
		RoomBase: RoomBase{
			subRooms: make([]Room, 0),
		},
	}
}

func (t *Tower) loadData(dir string) error {
	data, err := os.ReadFile(path.Join(dir, "tower.json"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, t)
	if err != nil {
		return err
	}
	return nil
}
