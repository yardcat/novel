package world

import (
	"encoding/json"
	"my_test/combat"
	"my_test/log"
	"os"
)

type CombatSystem struct {
	Monsters map[string]combat.Enemy
	Dungeons map[string]combat.Dungeon
	story    *Story
}

func NewCombatSystem() *CombatSystem {
	c := &CombatSystem{
		story: GetStory(),
	}
	c.loadData()
	return c
}

func (c *CombatSystem) GetEnemy(name string) combat.Enemy {
	return c.Monsters[name]
}

func (c *CombatSystem) StartCombat(actors []*combat.Actor, enemies []*combat.Enemy) {
	combat.NewCombat(actors, enemies, c).Start()
}

// OnDead implements combat.CombatClient.
func (c *CombatSystem) OnDead(combat.Combatable) {
	panic("unimplemented")
}

// OnDraw implements combat.CombatClient.
func (c *CombatSystem) OnDraw() {
	panic("unimplemented")
}

// OnKill implements combat.CombatClient.
func (c *CombatSystem) OnKill(combat.Combatable) {
	panic("unimplemented")
}

// OnLose implements combat.CombatClient.
func (c *CombatSystem) OnLose() {
	panic("unimplemented")
}

// OnWin implements combat.CombatClient.
func (c *CombatSystem) OnWin() {
	panic("unimplemented")
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
	jsonData, err := os.ReadFile(c.story.GetResources().GetPath("dungeon/test.json"))
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &c.Dungeons)
	if err != nil {
		return err
	}
	return nil
}
