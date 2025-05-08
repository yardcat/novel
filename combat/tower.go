package combat

import (
	"context"
	"encoding/json"
	"my_test/event"
	pb "my_test/event"
	"my_test/log"
	"my_test/push"
	"my_test/util"
	"os"
	"path"

	"github.com/jinzhu/copier"
)

const (
	ROOM_TYPE_NONE = iota
	ROOM_TYPE_FIGHT
	ROOM_TYPE_SHOP
	ROOM_TYPE_REST
	ROOM_TYPE_EVENT
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
	FloorNum      int                   `json:"floor_num"`
	RoomNum       int                   `json:"room_num"`
	ShopNum       int                   `json:"shop_num"`
	RestNum       int                   `json:"rest_num"`
	EventNum      int                   `json:"event_num"`
	EnemyMap      map[string]*CardEnemy `json:"enemies"`
	EnemyGroupMap map[int][]string      `json:"group"`
	PotionMap     map[string]*Potion    `json:"potion"`
	RelicMap      map[string]*Relic     `json:"relic"`
	cardMap       map[string]*Card
	careerMap     map[string]*CardCareer
	eventMap      map[string]*CardEvent

	currentCombat *CardCombat
	cards         []*Card
	potions       []*Potion
	potionLimit   int
	relics        []*Relic
	actor         *CardActor
	floor         *Floor
	effects       map[int][]*Effect
	resourceDir   string
	floorCount    int
	shopCount     int
	restCount     int
	eventCount    int
}

type TowerParams struct {
	Actor *CardActor
	Path  PathProvider
}

func NewTower() *Tower {
	t := &Tower{
		floorCount:  1,
		potionLimit: 3,
	}
	return t
}

func (t *Tower) Init(params *TowerParams) {
	t.actor = params.Actor
	t.resourceDir = params.Path.GetPath("card")

	t.loadData()
	t.generateFloor()
	t.PrepareCard()
}

func (t *Tower) Reset() {
	t.floorCount = 0
	t.shopCount = 0
	t.restCount = 0
	t.eventCount = 0
	t.PrepareCard()
	t.relics = []*Relic{}
	t.potions = []*Potion{}
	t.effects = make(map[int][]*Effect)
	t.currentCombat = nil
	t.actor = nil
}

func (t *Tower) EnterNextFloor() *Floor {
	t.generateFloor()
	t.floorCount++
	return t.floor
}

func (t *Tower) GetRoomTypeChoices() []int32 {
	choices := []int32{ROOM_TYPE_FIGHT}
	candidates := []int32{}
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

func (t *Tower) GetWelcomeEvents() []string {
	return []string{"strength", "max_health", "draw_card"}
}

func (t *Tower) PrepareCard() {
	career := t.careerMap["kongfu"]
	t.cards = make([]*Card, 0)
	t.cards = append(t.cards, career.InitCards...)
}

func (t *Tower) HandleEvent(ev string) {
	for _, effect := range t.eventMap[ev].Effects {
		t.currentCombat.handCardEffect(&effect, t.actor)
	}
	t.currentCombat.requestUpdateUI()
}

func (t *Tower) GetCombatBonus() CombatBonus {
	return CombatBonus{
		Cards:             []string{"dang", "quan", "jiao"},
		CardChooseCount:   1,
		Potions:           []string{""},
		PotionChooseCount: 1,
		Relics:            []string{""},
		RelicChooseCount:  1,
	}
}

func (t *Tower) StartCardCombat() *CardCombat {
	params := CardCombatParams{
		Actors:             []*CardActor{t.actor},
		Enemies:            t.fightRoom().Enemy,
		ResourceDir:        t.resourceDir,
		CardCombatDelegate: t,
		Cards:              t.cards,
	}
	t.currentCombat = NewCardCombat(&params)
	t.currentCombat.Start()
	return t.currentCombat
}

func (t *Tower) addBonus(bonus string, typ int32) {
	switch typ {
	case BONUS_TYPE_CARD:
		t.addBonus(bonus, typ)
	case BONUS_TYPE_POTION:
		t.addBonus(bonus, typ)
	case BONUS_TYPE_RELIC:
	}
}

func (t *Tower) AddCard(name string) {
	card := t.GetCard(name)
	t.cards = append(t.cards, card)
}

func (t *Tower) AddPotion(name string) {
	potion, exist := t.PotionMap[name]
	if !exist {
		log.Error("potion %s not exist", name)
		panic("potion not exist")
	}
	if len(t.potions) >= t.potionLimit {
		log.Info("reach potion limit %d", t.potionLimit)
		return
	}
	t.potions = append(t.potions, potion)
}

func (t *Tower) AddRelic(name string) {
	relic, exist := t.RelicMap[name]
	if !exist {
		log.Error("relic %s not exist", name)
		panic("relic not exist")
	}
	t.relics = append(t.relics, relic)
}

func (t *Tower) generateFloor() {
	fl := &Floor{}
	t.floor = fl
}

func (t *Tower) fightRoom() *FightRoom {
	return t.floor.room.(*FightRoom)
}

func (t *Tower) generateRoom(typ int) Room {
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
	return room
}

func (t *Tower) generateFightRoom() *FightRoom {
	room := &FightRoom{
		Enemy: []*CardEnemy{},
	}
	gr := t.EnemyGroupMap[t.floorCount]
	for _, v := range gr {
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

func (t *Tower) enterRoom(typ int) {
	room := t.generateRoom(typ)
	t.floor.room = room

	switch room.Type() {
	case ROOM_TYPE_FIGHT:
		t.StartCardCombat()
	case ROOM_TYPE_SHOP:
	case ROOM_TYPE_REST:
	case ROOM_TYPE_EVENT:
	}

	push.PushEvent(event.CardEnterRoomDone{
		Type: room.Type(),
	})
}

func (c *Tower) loadData() error {
	err := c.loadTower()
	if err != nil {
		return nil
	}
	err = c.loadCard()
	if err != nil {
		return err
	}
	err = c.loadCareer()
	if err != nil {
		return err
	}
	err = c.loadEvent()
	if err != nil {
		return err
	}
	err = c.loadPotion()
	if err != nil {
		return err
	}
	err = c.loadRelic()
	if err != nil {
		return err
	}
	return nil
}

func (t *Tower) loadTower() error {
	data, err := os.ReadFile(path.Join(t.resourceDir, "tower.json"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tower) loadCard() error {
	data, err := os.ReadFile(path.Join(t.resourceDir, "card.json"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &t.cardMap)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tower) loadCareer() error {
	data, err := os.ReadFile(path.Join(t.resourceDir, "career.json"))
	if err != nil {
		return err
	}
	adapter := map[string]struct {
		Cards []string `json:"init_cards"`
	}{}
	err = json.Unmarshal(data, &adapter)
	t.careerMap = make(map[string]*CardCareer, len(adapter))
	for k, v := range adapter {
		t.careerMap[k] = &CardCareer{
			InitCards: make([]*Card, 0, len(v.Cards)),
		}
		for _, card := range v.Cards {
			t.careerMap[k].InitCards = append(t.careerMap[k].InitCards, t.cardMap[card])
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (t *Tower) loadEvent() error {
	data, err := os.ReadFile(path.Join(t.resourceDir, "event.json"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &t.eventMap)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tower) loadPotion() error {
	data, err := os.ReadFile(path.Join(t.resourceDir, "potion.json"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &t.PotionMap)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tower) loadRelic() error {
	data, err := os.ReadFile(path.Join(t.resourceDir, "relic.json"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &t.RelicMap)
	if err != nil {
		return err
	}
	return nil
}

// CardCombatDelegate
func (t *Tower) GetCard(name string) *Card {
	_, exist := t.cardMap[name]
	if !exist {
		log.Error("get card %s not exist", name)
		panic("card not exist")
	}
	return t.cardMap[name]
}

func (t *Tower) OnLose() {
	push.PushEvent(event.CardCombatLose{})
	t.Reset()
}

func (t *Tower) OnWin() {
	bonus := t.GetCombatBonus()
	ev := event.CardCombatWin{}
	copier.Copy(&ev.Bonus, bonus)
	push.PushEvent(ev)

	t.EnterNextFloor()
}

func (t *Tower) OnPlayCard(card *Card) {

}

func (t *Tower) OnDiscardCard(card *Card) {
}

func (t *Tower) OnShuffle() {
}

func (t *Tower) OnRemoveCard(card *Card) {
}

func (t *Tower) OnDrawCard(card *Card) {
}

func (t *Tower) OnAddCard(card *Card) {
}

func (t *Tower) OnEnemyDead(enemy *CardEnemy) {
}

func (t *Tower) OnEnemyTurnStart() {
}

func (t *Tower) OnEnemyTurnEnd() {
}

func (t *Tower) OnEnenyDamage(enemy *CardEnemy, damage int) {

}

func (t *Tower) OnActorTurnEnd() {

}

func (t *Tower) OnActorTurnStart() {
	t.EffectOn(TIMING_ENEMY_TURN_START)
}

// grpc
func (t *Tower) Welcome(ctx context.Context, request *pb.WelcomeRequest) (*pb.WelcomeResponse, error) {
	t.enterRoom(ROOM_TYPE_FIGHT)
	t.HandleEvent(request.Event)

	return &pb.WelcomeResponse{
		Result: "ok",
	}, nil
}

func (t *Tower) SendCard(ctx context.Context, request *pb.SendCardRequest) (*pb.SendCardResponse, error) {
	t.currentCombat.UseCards(request.Cards, request.Target)

	return &pb.SendCardResponse{
		Result: "ok",
	}, nil
}

func (t *Tower) DiscardCard(ctx context.Context, request *pb.DiscardCardRequest) (*pb.DiscardCardResponse, error) {
	t.currentCombat.DiscardCards(request.Cards)

	return &pb.DiscardCardResponse{
		Result: "ok",
	}, nil
}

func (t *Tower) EndTurn(ctx context.Context, request *pb.EndTurnRequest) (*pb.EndTurnResponse, error) {
	t.currentCombat.EndTurn()

	return &pb.EndTurnResponse{
		Result: "ok",
	}, nil
}

func (t *Tower) NextFloor(ctx context.Context, request *pb.NextFloorRequest) (*pb.NextFloorResponse, error) {
	return &pb.NextFloorResponse{
		RoomChoices: t.GetRoomTypeChoices(),
	}, nil
}

func (t *Tower) EnterRoom(ctx context.Context, request *pb.EnterRoomRequest) (*pb.EnterRoomResponse, error) {
	t.enterRoom(int(request.Type))

	return &pb.EnterRoomResponse{
		Result: "ok",
	}, nil
}

func (t *Tower) ChooseBonus(ctx context.Context, request *pb.ChooseBonusRequest) (*pb.ChooseBonusResponse, error) {
	bonus := request.Bonus

	// TODO:validate
	for _, v := range bonus {
		t.addBonus(v.Name, v.Type)
	}

	return &pb.ChooseBonusResponse{
		Result: "ok",
	}, nil
}
