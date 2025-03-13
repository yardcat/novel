package world

import (
	"encoding/json"
	"my_test/combat"
	"my_test/log"
	"my_test/push"
	"my_test/util"
	"os"
	"path/filepath"
	"time"
)

type CombatSystem struct {
	Monsters map[string]*combat.Enemy
	Dungeons map[string]*combat.Dungeon
	story    *Story
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

func (c *CombatSystem) StartCombat(actors []*combat.Actor, enemies []*combat.Enemy) {
	combat.NewCombat(actors, enemies, c).Start()
}

func (c *CombatSystem) ChallengeDungeon(name string) error {
	actor := combat.NewActor(0, "player")
	actors := []*combat.Actor{actor}
	dg := c.Dungeons[name]
	for _, group := range dg.Groups {
		log.Info("start combat group %s", group.Name)
		enemies := combat.CreateEnemyGroup(group)
		combat.NewCombat(actors, enemies.Enemies, c).Start()
	}
	return nil
}

// OnDead implements combat.CombatClient.
func (c *CombatSystem) OnDead(combat.Combatable) {
	log.Info("OnDead is unimplemented")
}

// OnDraw implements combat.CombatClient.
func (c *CombatSystem) OnDraw() {
	log.Info("OnDraw is unimplemented")
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
		push.PushEvent(CombatWinEvent{Result: "win"})
	}()
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
		dungeonName := util.GetPureFileName(file)
		c.Dungeons[dungeonName] = dungeon
	}

	return nil
}
