package world

import (
	"encoding/json"
	"my_test/log"
	"os"
)

type CraftNeedValue struct {
	need  string
	value int
}

type CraftNeed struct {
	skill []CraftNeedValue
	item  []CraftNeedValue
}

type CraftItem struct {
	Name        string
	Description string
	Energy      int
	Need        CraftNeed
	Attribute   map[string]int
	Effects     map[string]int
}

type CraftSystem struct {
	Recipes map[string]CraftItem `json:"Recipe"`
	story   *Story
}

func NewCraftSystem(story *Story) *CraftSystem {
	c := &CraftSystem{
		Recipes: make(map[string]CraftItem),
		story:   story,
	}
	c.loadData()
	return c
}

func (c *CraftSystem) Craft(name string) {
}

func (c *CraftSystem) loadData() error {
	file := c.story.GetResources().GetPath("item/craft.json")
	jsonData, err := os.ReadFile(file)
	if err != nil {
		log.Error("load config file err: %v", err)
		return err
	}

	err = json.Unmarshal(jsonData, &c)
	if err != nil {
		log.Error("unmarshal config file err: %v", err)
		return err
	}
	return nil
}
