package world

import (
	"encoding/json"
	"my_test/combat"
	"my_test/event"
	"my_test/log"
	"my_test/push"
	"my_test/util"
	"os"
	"path/filepath"
	"time"
)

type CombatSystem struct {
	Monsters   map[string]*combat.Enemy
	Dungeons   map[string]*combat.Dungeon
	story      *Story
	cardCombat *combat.CardCombat
	tower      *combat.Tower
}

func NewCombatSystem() *CombatSystem {
	c := &CombatSystem{
		story: GetStory(),
	}
	c.loadData()
	return c
}

func (c *CombatSystem) GetEnemy(name string) *combat.Enemy {
	return c.Monsters[name]
}

func (c *CombatSystem) ChallengeDungeon(name string) error {
	player := c.story.GetPlayer("0")
	player.AddCareer("doctor")
	actor := combat.NewActor(player.GetCombatableBase(), player)
	petName := "dog_pet"
	player.AddPet(petName)
	pet := player.GetPet(petName)
	npcName := "SunWuKong"
	player.AddNpc(npcName)
	npc := player.GetNpc(npcName)
	npcActor := combat.NewActor(npc.GetCombatableBase(), npc)
	petActor := combat.NewActor(pet.GetCombatableBase(), pet)
	actors := []*combat.Actor{actor, petActor, npcActor}
	dg := c.Dungeons[name]
	for _, group := range dg.Groups {
		log.Info("start combat group %s", group.Name)
		enemies := combat.CreateEnemyGroup(group)
		params := combat.AutoCombatParams{
			Actors:  actors,
			Enemies: enemies.Enemies,
			Client:  c,
		}
		combat.NewAutoCombat(&params).Start()
	}
	return nil
}

func (c *CombatSystem) ChallengeTower(ev *event.CardStartEvent) *event.CardStartEventReply {
	reply := &event.CardStartEventReply{}

	player := c.story.GetPlayer("0")
	player.AddCareer("doctor")
	actor := combat.NewCardActor(player.GetCombatableBase())
	actor.Name = "winter"
	params := &combat.TowerParams{
		Actor: actor,
		Path:  c,
	}
	c.tower = combat.NewTower(params)
	reply.Events = c.tower.GetWelcomeEvents()

	return reply
}

func (c *CombatSystem) SendCards(ev *event.CardSendCards) *event.CardSendCardsReply {
	return c.cardCombat.UseCards(ev)
}

func (c *CombatSystem) DiscardCards(ev *event.CardDiscardCards) *event.CardDiscardCardsReply {
	return c.cardCombat.DiscardCards(ev)
}

func (c *CombatSystem) EndTurn(ev *event.CardTurnEndEvent) *event.CardTurnEndEventReply {
	return c.cardCombat.EndTurn(ev)
}

func (c *CombatSystem) NextFloor(ev *event.CardNextFloorEvent) *event.CardNextFloorReply {
	reply := &event.CardNextFloorReply{}
	reply.ChooseRoom = c.tower.GetRoomTypeChoices()
	return reply
}

func (c *CombatSystem) EnterRoom(ev *event.CardEnterRoomEvent) *event.CardEnterRoomReply {
	reply := &event.CardEnterRoomReply{}
	c.tower.EnterRoom(ev.RoomType)
	if ev.RoomType == combat.ROOM_TYPE_FIGHT {
		c.tower.StartCardCombat()
	}
	return reply
}

func (c *CombatSystem) HandleWelcome(ev *event.CardWelcomeEvent) *event.CardWelcomeReply {
	reply := &event.CardWelcomeReply{
		Results: make(map[string]any),
	}
	c.tower.EnterRoom(combat.ROOM_TYPE_FIGHT)
	c.cardCombat = c.tower.StartCardCombat()
	c.tower.HandleEvent(ev.Event)

	return reply
}

// OnKill implements combat.CombatClient.
func (c *CombatSystem) OnKill(combat.Combatable) {
	log.Info("OnKill is unimplemented")
}

// OnLose implements combat.CombatClient.
func (c *CombatSystem) OnLose() {
	log.Info("OnLose is unimplemented")
}

// OnWin implements combat.CombatClient.
func (c *CombatSystem) OnWin() {
	log.Info("OnWin is unimplemented")
	go func() {
		time.Sleep(1 * time.Second)
		push.PushEvent(event.CombatWinEvent{Result: "win"})
	}()
}

func (c *CombatSystem) GetPath(path string) string {
	return c.story.GetResources().GetPath(path)
}

func (c *CombatSystem) loadData() error {
	err := c.loadMonsters()
	if err != nil {
		log.Info("load monsters err %v", err)
	}
	err = c.loadDungeons()
	if err != nil {
		log.Info("load dungeons err %v", err)
	}
	return nil
}

func (c *CombatSystem) loadMonsters() error {
	jsonData, err := os.ReadFile(c.story.GetResources().GetPath("enemy/monster.json"))
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &c.Monsters)
	if err != nil {
		return err
	}
	return nil
}

func (c *CombatSystem) loadDungeons() error {
	dungeonFiles, err := filepath.Glob(c.story.GetResources().GetPath("dungeon/*.json"))
	if err != nil {
		return err
	}

	c.Dungeons = make(map[string]*combat.Dungeon)
	for _, file := range dungeonFiles {
		jsonData, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		jsonDungeon := make(map[string]struct {
			Name    string
			Members []string
		})

		err = json.Unmarshal(jsonData, &jsonDungeon)
		if err != nil {
			return err
		}
		dungeon := &combat.Dungeon{
			Name:   file,
			Groups: make([]combat.EnemyGroup, len(jsonDungeon)),
		}
		groupID := 0
		for k, v := range jsonDungeon {
			dungeon.Groups[groupID].Name = k
			dungeon.Groups[groupID].Enemies = make([]*combat.Enemy, len(v.Members))
			for i, name := range v.Members {
				dungeon.Groups[groupID].Enemies[i] = c.Monsters[name]
			}
			groupID++
		}
		dungeonName := util.GetFileNameWithoutExt(file)
		c.Dungeons[dungeonName] = dungeon
	}

	return nil
}
