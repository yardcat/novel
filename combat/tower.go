package combat

import (
	"context"
	"encoding/json"
	"fmt"
	"my_test/event"
	"my_test/log"
	"my_test/pb"
	"my_test/push"
	"my_test/util"
	"os"
	"path"
	"reflect"
	"slices"
	"strings"

	"github.com/bilibili/gengine/builder"
	gctx "github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"github.com/jinzhu/copier"
	"github.com/samber/lo"
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

type RestRoom struct {
}

func (r *RestRoom) Type() int {
	return ROOM_TYPE_REST

}

func (r *RestRoom) Heal(cbt Combatable) {

}

func (r *RestRoom) Update() {

}

type DestinyRoom struct {
	Event string
}

func (r *DestinyRoom) Type() int {
	return ROOM_TYPE_EVENT
}

type Floor struct {
	room Room
}

type Tower struct {
	FloorNum       int              `json:"floor_num"`
	RoomNum        int              `json:"room_num"`
	ShopNum        int              `json:"shop_num"`
	RestNum        int              `json:"rest_num"`
	EventNum       int              `json:"event_num"`
	EnemyGroupMap  map[int][]string `json:"group"`
	EnemyMap       map[string]*CardEnemy
	PotionMap      map[string]*Potion
	RelicMap       map[string]*Relic
	cardMap        map[string]*Card
	careerMap      map[string]*CardCareer
	eventMap       map[string]*CardEvent
	cardBindingMap map[string]reflect.Type

	currentCombat   *CardCombat
	cards           []*Card
	potions         []*Potion
	potionLimit     int
	relics          []*Relic
	actor           *CardActor
	floor           *Floor
	timingCallbacks map[int][]any
	resourceDir     string
	floorCount      int
	shopCount       int
	restCount       int
	eventCount      int

	// script
	effects     map[int][]*Effect
	effectRules strings.Builder
	engine      *engine.Gengine
	dataContext *gctx.DataContext
	ruleBuilder *builder.RuleBuilder
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
	towerInstance = t
	return t
}

var towerInstance *Tower

func GetTower() *Tower {
	if towerInstance == nil {
		panic("tower is not create")
	}
	return towerInstance
}

func (t *Tower) Init(params *TowerParams) {
	t.actor = params.Actor
	t.resourceDir = params.Path.GetPath("card")
	t.effects = make(map[int][]*Effect)
	t.timingCallbacks = make(map[int][]any)

	t.registerCardBindings()
	t.loadData()
	t.generateFloor()
	t.PrepareCard()
	t.initScript()
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
	t.timingCallbacks = make(map[int][]any)
	t.currentCombat = nil
	t.actor = nil
	t.initScript()
}

func (t *Tower) initScript() {
	t.dataContext = gctx.NewDataContext()
	t.dataContext.Add("actor", t.actor)
	t.dataContext.Add("t", t)
	t.dataContext.Add("log", fmt.Println)
	t.dataContext.Add("Sprintf", fmt.Sprintf)
	t.ruleBuilder = builder.NewRuleBuilder(t.dataContext)
	err := t.ruleBuilder.BuildRuleFromString(t.effectRules.String())

	if err != nil {
		panic("initScirpt error")
	}

	t.engine = engine.NewGengine()
}

func (t *Tower) EnterNextFloor() *Floor {
	t.generateFloor()
	t.floorCount++
	return t.floor
}

func (t *Tower) GetRoomTypeChoices() []int {
	choices := []int{ROOM_TYPE_FIGHT}
	candidates := []int{}
	if t.shopCount < t.ShopNum {
		candidates = append(candidates, ROOM_TYPE_SHOP)
	}
	if t.restCount < t.RestNum {
		candidates = append(candidates, ROOM_TYPE_REST)
	}
	if t.eventCount < t.EventNum {
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

func (t *Tower) HandleDestiny(ev string) {
	t.currentCombat.requestUpdateUI()
}

func (t *Tower) regiserTimingCallback(timing int, callback any) {
	if t.timingCallbacks[timing] == nil {
		t.timingCallbacks[timing] = make([]any, 0)
	}
	t.timingCallbacks[timing] = append(t.timingCallbacks[timing], callback)
}

func (t *Tower) getCombatBonus() CombatBonus {
	return CombatBonus{
		Cards:             []string{"flex", "cleave", "clash"},
		CardChooseCount:   1,
		Potions:           []string{"test"},
		PotionChooseCount: 1,
		Relics:            []string{"test"},
		RelicChooseCount:  1,
	}
}

func (t *Tower) startCardCombat() *CardCombat {
	params := CardCombatParams{
		Actors:             []*CardActor{t.actor},
		Enemies:            t.fightRoom().Enemy,
		ResourceDir:        t.resourceDir,
		CardCombatDelegate: t,
		Cards:              t.cards,
	}
	t.currentCombat = NewCardCombat(&params)
	t.dataContext.Add("combat", t.currentCombat)
	t.dataContext.Add("turn", t.currentCombat.turnInfo)
	t.currentCombat.Start()
	t.onStartCombat()

	return t.currentCombat
}

func (t *Tower) onStartCombat() {
	t.EffectOn(TIMING_COMBAT_START)
}

func (t *Tower) addBonus(name string, typ int32) {
	switch typ {
	case BONUS_TYPE_CARD:
		t.AddCard(name)
	case BONUS_TYPE_POTION:
		t.AddPotion(name)
	case BONUS_TYPE_RELIC:
		t.AddRelic(name)
	}
}

func (t *Tower) AddCard(name string) {
	card := t.GetCard(name)
	t.cards = append(t.cards, card)
}

func (t *Tower) UpgradeCard(card *Card) {
	if card.Upgraded {
		return
	}
	t.dataContext.Add("card", card)
	upgradeRule := card.Name + "_upgrade"
	err := t.engine.ExecuteSelectedRules(t.ruleBuilder, []string{upgradeRule})
	if err != nil {
		log.Error("upgrade card %s error: %s", card.Name, err)
		panic(err)
	}
	card.Upgraded = true
	t.EffectOn(TIMING_UPGRADE_CARD)
}

func (t *Tower) UpdatePotionUI() {
	ev := event.CardUpdatePotion{}
	copier.Copy(&ev.Potions, t.potions)
	push.PushEvent(ev)
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

	for _, v := range potion.Effects {
		if v.Timing != TIMING_IMMEDIATE {
			v.CasterID = potion.Name
			v.CasterType = CASTER_TYPE_POTION
			t.effects[v.Timing] = append(t.effects[v.Timing], v)
		}
	}

	t.UpdatePotionUI()
}

func (t *Tower) RemovePotion(name string) {
	potion, exist := t.PotionMap[name]
	if !exist {
		log.Error("potion %s not exist", name)
		panic("potion not exist")
	}
	for i, v := range t.potions {
		if v == potion {
			t.potions = slices.Delete(t.potions, i, i+1)
		}
	}

	t.UpdatePotionUI()
}

func (t *Tower) usePotion(name string) bool {
	if t.currentCombat != nil && t.currentCombat.finish {
		return false
	}

	potion, exsit := t.PotionMap[name]
	if !exsit {
		log.Error("potion %s not exist", name)
		panic("potion not exist")
	}
	for _, v := range potion.Effects {
		t.UseEffect(v)
		push.PushAction("use potion: %s", name)
	}

	t.RemovePotion(name)

	return true
}

func (t *Tower) discardPotion(name string) {
	t.RemovePotion(name)
}

func (t *Tower) UpdateRelicUI() {
	ev := event.CardUpdateRelic{}
	copier.Copy(&ev.Relics, t.relics)
	push.PushEvent(ev)
}

func (t *Tower) AddRelic(name string) {
	relic, exist := t.RelicMap[name]
	if !exist {
		log.Error("relic %s not exist", name)
		panic("relic not exist")
	}
	t.relics = append(t.relics, relic)
	for _, v := range relic.Effects {
		if v.Timing != TIMING_IMMEDIATE {
			v.CasterID = relic.Name
			v.CasterType = CASTER_TYPE_RELIC
			t.effects[v.Timing] = append(t.effects[v.Timing], v)
		} else {
			t.UseEffect(v)
		}
	}

	t.UpdateRelicUI()
}

func (t *Tower) RemoveRelic(name string) {
	relic, exist := t.RelicMap[name]
	if !exist {
		log.Error("relic %s not exist", name)
		panic("relic not exist")
	}
	for i, v := range t.relics {
		if relic == v {
			t.relics = slices.Delete(t.relics, i, i+1)
			break
		}
	}
	for _, v := range relic.Effects {
		arr := t.effects[v.Timing]
		t.effects[v.Timing] = lo.Filter(arr, func(v *Effect, i int) bool {
			return v.CasterType == CASTER_TYPE_RELIC && v.CasterID == name
		})
	}

	t.UpdateRelicUI()
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

func (t *Tower) generateShopRoom() *Shop {
	room := NewShop(t)
	return room
}

func (t *Tower) generateRestRoom() *RestRoom {
	room := &RestRoom{}
	return room
}

func (t *Tower) generateEventRoom() *DestinyRoom {
	room := &DestinyRoom{}
	return room
}

func (t *Tower) enterRoom(typ int) {
	room := t.generateRoom(typ)
	t.floor.room = room

	switch room.Type() {
	case ROOM_TYPE_FIGHT:
		t.startCardCombat()
	case ROOM_TYPE_SHOP:
		t.showShopItems()
	case ROOM_TYPE_REST:
	case ROOM_TYPE_EVENT:
	}

	push.PushEvent(event.CardEnterRoomDone{
		Type: room.Type(),
	})
}

func (t *Tower) showShopItems() {
	push.PushEvent(event.CardUpdateShopUI{
		Cards:   []string{"flex", "cleave", "clash"},
		Potions: []string{"test"},
		Relics:  []string{"test"},
	})
}

func (c *Tower) loadData() error {
	if err := c.loadTower(); err != nil {
		panic(fmt.Errorf("failed to load tower: %w", err))
	}
	if err := c.loadCard(); err != nil {
		panic(fmt.Errorf("failed to load cards: %w", err))
	}
	if err := c.loadCareer(); err != nil {
		panic(fmt.Errorf("failed to load careers: %w", err))
	}
	if err := c.loadEvent(); err != nil {
		panic(fmt.Errorf("failed to load events: %w", err))
	}
	if err := c.loadPotion(); err != nil {
		panic(fmt.Errorf("failed to load potions: %w", err))
	}
	if err := c.loadRelic(); err != nil {
		panic(fmt.Errorf("failed to load relics: %w", err))
	}
	if err := c.loadEnemy(); err != nil {
		panic(fmt.Errorf("failed to load enemy: %w", err))
	}
	if err := c.loadScript(); err != nil {
		panic(fmt.Errorf("failed to load script: %w", err))
	}

	return nil
}

func (t *Tower) loadTower() error {
	data, err := os.ReadFile(path.Join(t.resourceDir, "tower.json"))
	if err != nil {
		return fmt.Errorf("error reading tower.json: %w", err)
	}
	if err := json.Unmarshal(data, t); err != nil {
		return fmt.Errorf("error unmarshaling tower.json: %w", err)
	}
	return nil
}

func (t *Tower) loadCard() error {
	data, err := os.ReadFile(path.Join(t.resourceDir, "card.json"))
	if err != nil {
		return fmt.Errorf("error reading card.json: %w", err)
	}
	if err := json.Unmarshal(data, &t.cardMap); err != nil {
		return fmt.Errorf("error unmarshaling card.json: %w", err)
	}
	for k, v := range t.cardMap {
		v.Id = k
	}
	return nil
}

func (t *Tower) loadCareer() error {
	data, err := os.ReadFile(path.Join(t.resourceDir, "career.json"))
	if err != nil {
		return fmt.Errorf("error reading career.json: %w", err)
	}
	adapter := map[string]struct {
		Cards []string `json:"init_cards"`
	}{}
	if err := json.Unmarshal(data, &adapter); err != nil {
		return fmt.Errorf("error unmarshaling career.json: %w", err)
	}
	t.careerMap = make(map[string]*CardCareer, len(adapter))
	for k, v := range adapter {
		t.careerMap[k] = &CardCareer{
			InitCards: make([]*Card, 0, len(v.Cards)),
		}
		for _, card := range v.Cards {
			t.careerMap[k].InitCards = append(t.careerMap[k].InitCards, t.cardMap[card])
		}
	}
	return nil
}

func (t *Tower) loadEvent() error {
	data, err := os.ReadFile(path.Join(t.resourceDir, "event.json"))
	if err != nil {
		return fmt.Errorf("error reading event.json: %w", err)
	}
	if err := json.Unmarshal(data, &t.eventMap); err != nil {
		return fmt.Errorf("error unmarshaling event.json: %w", err)
	}
	return nil
}

func (t *Tower) loadPotion() error {
	data, err := os.ReadFile(path.Join(t.resourceDir, "potion.json"))
	if err != nil {
		return fmt.Errorf("error reading potion.json: %w", err)
	}
	if err := json.Unmarshal(data, &t.PotionMap); err != nil {
		return fmt.Errorf("error unmarshaling potion.json: %w", err)
	}
	return nil
}

func (t *Tower) loadRelic() error {
	fileName := "relic.json"
	data, err := os.ReadFile(path.Join(t.resourceDir, fileName))
	if err != nil {
		return fmt.Errorf("error reading %s: %w", fileName, err)
	}
	if err := json.Unmarshal(data, &t.RelicMap); err != nil {
		return fmt.Errorf("error unmarshaling %s : %w", fileName, err)
	}
	return nil
}

func (t *Tower) loadEnemy() error {
	fileName := "enemy.json"
	data, err := os.ReadFile(path.Join(t.resourceDir, fileName))
	if err != nil {
		return fmt.Errorf("error reading %s: %w", fileName, err)
	}

	if err := json.Unmarshal(data, &t.EnemyMap); err != nil {
		return fmt.Errorf("error unmarshaling %s: %w", fileName, err)
	}

	return nil
}

func (t *Tower) loadScript() error {
	data, err := os.ReadFile(path.Join(t.resourceDir, "rule.gengine"))
	if err != nil {
		return fmt.Errorf("error reading rule.gengine: %w", err)
	}
	t.effectRules.Write(data)

	data, err = os.ReadFile(path.Join(t.resourceDir, "upgrade.gengine"))
	if err != nil {
		return fmt.Errorf("error reading upgrade.gengine: %w", err)
	}
	t.effectRules.Write(data)

	data, err = os.ReadFile(path.Join(t.resourceDir, "enemy.gengine"))
	if err != nil {
		return fmt.Errorf("error reading enemy.gengine: %w", err)
	}
	t.effectRules.Write(data)

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

func (t *Tower) CanUse(card *Card) bool {
	if card.CanUse == "" {
		return true
	}
	err := t.engine.ExecuteSelectedRules(t.ruleBuilder, []string{card.CanUse})
	if err != nil {
		panic(err)
	}
	result, _ := t.engine.GetRulesResultMap()
	return result["can_use"].(bool)
}

func (t *Tower) UpgradeCardInCombat(card *Card) {
	if card.Upgraded {
		return
	}
	t.dataContext.Add("card", card)
	upgradeRule := card.Name + "_upgrade"
	err := t.engine.ExecuteSelectedRules(t.ruleBuilder, []string{upgradeRule})
	if err != nil {
		log.Error("upgrade card %s error: %s", card.Name, err)
		panic(err)
	}
	card.Upgraded = true
}

func (t *Tower) AddEnemyEffect(effect *Effect) {
	t.effects[effect.Timing] = append(t.effects[effect.Timing], effect)
}

func (t *Tower) EnableBuff() {
	for _, v := range t.effects {
		for _, effect := range v {
			if effect.Type == EFFECT_TYPE_BUFF {
				effect.Enabled = true
			}
		}
	}
}

func (t *Tower) TriggerEffect(effect *Effect, bindings map[string]any) {
	if len(bindings) > 0 {
		for k, v := range bindings {
			t.dataContext.Add(k, v)
		}
	}
	t.UseEffect(effect)
}

func (t *Tower) TriggerEnemyAction(action string, binding map[string]any) {
	if len(binding) > 0 {
		for k, v := range binding {
			t.dataContext.Add(k, v)
		}
	}
	err := t.engine.ExecuteSelectedRules(t.ruleBuilder, []string{action})
	if err != nil {
		log.Error("execute enemy action %s error: %s", action, err)
		panic(err)
	}
}

func (t *Tower) TriggerTiming(timing int, bindings map[string]any) {
	if len(bindings) > 0 {
		for k, v := range bindings {
			t.dataContext.Add(k, v)
		}
	}
	t.EffectOn(timing)
}

func (t *Tower) TriggerPrepareIntent() {
	t.dataContext.Add("ai", t.currentCombat.ai)
	for _, v := range t.currentCombat.enemies {
		t.dataContext.Add("enemy", v)
		err := t.engine.ExecuteSelectedRules(t.ruleBuilder, []string{v.Move})
		if err != nil {
			log.Error("execute enemy move %s error: %s", v.Move, err)
			panic(err)
		}
	}
}

func (t *Tower) OnLose() {
	push.PushEvent(event.CardCombatLose{})
	t.Reset()
}

func (t *Tower) OnWin() {
	bonus := t.getCombatBonus()
	ev := event.CardCombatWin{}
	copier.Copy(&ev.Bonus, bonus)
	ev.NextFloor = t.GetRoomTypeChoices()
	push.PushEvent(ev)

	t.EnterNextFloor()
}

func (t *Tower) OnUseCard(card *Card) {
	timing := TIMING_USE_CARD
	t.EffectOn(timing)
	for _, v := range t.timingCallbacks[timing] {
		callback := v.(func(*Card))
		callback(card)
	}
	if card.Binding != nil {
		card.Binding.Use(t, card)
	}
}

func (t *Tower) OnDiscardCard(card *Card) {
}

func (t *Tower) OnShuffle() {
}

func (t *Tower) OnExhaustCard(card *Card) {
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
	t.EffectOn(TIMING_ENEMY_TURN_END)
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
	t.HandleDestiny(request.Event)

	return &pb.WelcomeResponse{
		Result: "ok",
	}, nil
}

func (t *Tower) CanUseCard(ctx context.Context, request *pb.CanUseRequest) (*pb.CanUseResponse, error) {
	result := t.CanUse(t.currentCombat.hand[request.Card])
	return &pb.CanUseResponse{
		Result: result,
	}, nil
}

func (t *Tower) SendCard(ctx context.Context, request *pb.SendCardRequest) (*pb.SendCardResponse, error) {
	errList := make([]error, 0)
	for _, v := range request.Cards {
		err := t.currentCombat.UseCard(v, request.Target)
		if err != nil {
			errList = append(errList, err)
			continue
		}
	}

	if len(errList) > 0 {
		return &pb.SendCardResponse{
			Result: "error",
		}, nil
	}

	return &pb.SendCardResponse{
		Result: "ok",
	}, nil
}

func (t *Tower) DiscardCard(ctx context.Context,
	request *pb.DiscardCardRequest) (*pb.DiscardCardResponse, error) {
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
	t.enterRoom(int(request.Room))

	return &pb.NextFloorResponse{
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

func (t *Tower) UsePotion(ctx context.Context, request *pb.UsePotionRequest) (*pb.UsePotionResponse, error) {
	t.usePotion(request.Name)

	return &pb.UsePotionResponse{
		Result: "ok",
	}, nil
}

func (t *Tower) DiscardPotion(ctx context.Context, request *pb.DiscardPotionRequest) (*pb.DiscardPotionResponse, error) {
	t.discardPotion(request.Name)

	return &pb.DiscardPotionResponse{
		Result: "ok",
	}, nil
}

func (t *Tower) Buy(ctx context.Context, request *pb.BuyRequest) (*pb.BuyResponse, error) {
	return nil, nil
}
