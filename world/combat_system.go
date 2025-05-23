package world

import (
	"context"
	"encoding/json"
	"fmt"
	"my_test/combat"
	"my_test/event"
	"my_test/log"
	"my_test/pb"
	"my_test/push"
	"my_test/util"
	"net"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
)

type CombatSystem struct {
	Monsters map[string]*combat.Enemy
	Dungeons map[string]*combat.Dungeon
	story    *Story
	tower    *combat.Tower
	server   *grpc.Server
}

func NewCombatSystem() *CombatSystem {
	c := &CombatSystem{
		story: GetStory(),
	}
	c.loadData()
	err := c.initGrpc()
	if err != nil {
		log.Info("init grpc err %v", err)
		return nil
	}
	return c
}

func (c *CombatSystem) initGrpc() error {
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return err
	}
	c.server = grpc.NewServer()
	pb.RegisterWorldServer(c.server, c)

	c.tower = combat.NewTower()
	pb.RegisterCardServer(c.server, c.tower)

	go c.server.Serve(lis)
	return nil
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

// func (c *CombatSystem) HandleWelcome(ev *event.CardWelcomeEvent) *event.CardWelcomeReply {
// 	reply := &event.CardWelcomeReply{
// 		Results: make(map[string]any),
// 	}
// 	c.tower.EnterRoom(combat.ROOM_TYPE_FIGHT)
// 	c.cardCombat = c.tower.StartCardCombat()
// 	c.tower.HandleEvent(ev.Event)

// 	return reply
// }

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

// grpc
func (c *CombatSystem) StartCard(context.Context, *pb.StartCardRequest) (*pb.StartCardResponse, error) {
	if c.tower.HasInit() {
		c.tower.Reset()
	}
	player := c.story.GetPlayer("0")
	player.AddCareer("doctor")
	actor := combat.NewCardActor(player.GetCombatableBase())
	actor.Name = "winter"
	params := &combat.TowerParams{
		Actor: actor,
		Path:  c,
	}
	c.tower.Init(params)

	response := &pb.StartCardResponse{
		Choices: c.tower.GetWelcomeEvents(),
	}
	return response, nil
}
