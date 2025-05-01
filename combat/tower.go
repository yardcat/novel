package combat

import (
	"encoding/json"
	"my_test/util"
	"os"
	"path"
)

const (
	ROOM_TYPE_FIGHT = iota
	ROOM_TYPE_SHOP
	ROOM_TYPE_EVENT
	ROOM_TYPE_REST
	ROOM_TYPE_COUNT
)

type Room interface {
	Type() int
}

type FightRoom struct {
	Enemy []*CardEnemy
	Bouns []string
}

func (r *FightRoom) Type() int {
	return ROOM_TYPE_FIGHT
}

type ShopRoom struct {
	Potions []Potion
	Relics  []Relic
}

type RestRoom struct {
}

func (r *RestRoom) Type() int {
	return ROOM_TYPE_REST

}

func (r *RestRoom) Heal(cbt Combatable) {

}

func (r *RestRoom) Update() {

}

func (r *ShopRoom) Type() int {
	return ROOM_TYPE_SHOP
}

type EventRoom struct {
	Event string
}

func (r *EventRoom) Type() int {
	return ROOM_TYPE_EVENT
}

type Floor struct {
	room Room
}

type Tower struct {
	FloorNum      int                    `json:"floor_num"`
	RoomNum       int                    `json:"room_num"`
	ShopNum       int                    `json:"shop_num"`
	RestNum       int                    `json:"rest_num"`
	EventNum      int                    `json:"event_num"`
	EnemyMap      map[string]*CardEnemy  `json:"enemies"`
	EnemyGroupMap map[int]CardEnemyGroup `json:"group"`

	actor      *CardActor
	floor      *Floor
	path       PathProvider
	floorCount int
	shopCount  int
	restCount  int
	eventCount int
}

type TowerParams struct {
	Actor *CardActor
	Path  PathProvider
}

func NewTower(params *TowerParams) *Tower {
	t := &Tower{
		actor: params.Actor,
		path:  params.Path,
	}
	t.loadData(params.Path.GetPath("card"))
	t.generateFloor()
	return t
}

func (t *Tower) EnterNextFloor() *Floor {
	t.generateFloor()
	t.floorCount++
	return t.floor
}

func (t *Tower) GetRoomTypeChoices() []int {
	choices := []int{ROOM_TYPE_FIGHT}
	candidates := []int{}
	if t.shopCount != 0 {
		candidates = append(candidates, ROOM_TYPE_SHOP)
	}
	if t.restCount != 0 {
		candidates = append(candidates, ROOM_TYPE_REST)
	}
	if t.eventCount != 0 {
		candidates = append(candidates, ROOM_TYPE_EVENT)
	}
	dice := util.GetRandomInt(len(candidates))
	choices = append(choices, candidates[dice])
	return choices
}

func (t *Tower) generateFloor() {
	fl := &Floor{}
	t.floor = fl
}

func (t *Tower) EnterRoom(typ int) Room {
	var room Room
	switch typ {
	case ROOM_TYPE_FIGHT:
		room = t.generateFightRoom()
	case ROOM_TYPE_SHOP:
		room = t.generateShopRoom()
		t.shopCount++
	case ROOM_TYPE_REST:
		room = t.generateRestRoom()
		t.restCount++
	case ROOM_TYPE_EVENT:
		room = t.generateEventRoom()
		t.eventCount++
	}
	t.floor.room = room
	return room
}

func (t *Tower) StartCardCombat() {
	params := CardCombatParams{
		Actors:  []*CardActor{t.actor},
		Enemies: t.fightRoom().Enemy,
		Path:    t.path,
		Client:  t,
	}
	cardCombat := NewCardCombat(&params)
	cardCombat.Start()
	cardCombat.GetCardTurnInfo()
}

func (t *Tower) fightRoom() *FightRoom {
	return t.floor.room.(*FightRoom)
}

func (t *Tower) generateFightRoom() *FightRoom {
	room := &FightRoom{
		Enemy: []*CardEnemy{},
	}
	gr := t.EnemyGroupMap[t.floorCount]
	for _, v := range gr.group {
		enemy := t.EnemyMap[v]
		room.Enemy = append(room.Enemy, NewCardEnemy(enemy))
	}

	return room
}

func (t *Tower) generateShopRoom() *ShopRoom {
	room := &ShopRoom{}
	return room
}

func (t *Tower) generateRestRoom() *RestRoom {
	room := &RestRoom{}
	return room
}

func (t *Tower) generateEventRoom() *EventRoom {
	room := &EventRoom{}
	return room
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

func (t *Tower) OnLose() {

}
func (t *Tower) OnWin() {

}
func (t *Tower) OnKill(Combatable) {

}
