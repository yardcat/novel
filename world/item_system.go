package world

import (
	"encoding/json"
	"my_test/log"
	"os"
)

type Item interface {
	GetId() int
	GetName() string
	GetDescription() string
}

type ItemSystem struct {
	ItemMap   map[int]Item
	Path2Id   map[string]int
	Resources *Resources
	idInc     int
}

type ItemData struct {
	Id          int
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (i *ItemData) GetId() int {
	return i.Id
}

func (i *ItemData) GetName() string {
	return i.Name
}

func (i *ItemData) GetDescription() string {
	return i.Description
}

type BuildItem struct {
	ItemData
	Energy int
}

type FoodItem struct {
	ItemData
	Effects map[string]int
}

type WeaponItem struct {
	ItemData
	Attributes map[string]int
}

func NewItemSystem(resources *Resources) *ItemSystem {
	itemSystem := &ItemSystem{
		ItemMap:   make(map[int]Item),
		Path2Id:   make(map[string]int),
		Resources: resources,
		idInc:     0,
	}
	itemSystem.loadStuff()
	return itemSystem
}

func (s *ItemSystem) GetItemId(path string) int {
	return s.Path2Id[path]
}

func (s *ItemSystem) GetItemById(id int) Item {
	return s.ItemMap[id]
}

func (s *ItemSystem) GetItemByName(path string) Item {
	id := s.Path2Id[path]
	return s.ItemMap[id]
}

func (s *ItemSystem) loadStuff() error {
	err := s.loadBuild()
	if err != nil {
		log.Info("load build err %v", err)
	}
	err = s.loadFood()
	if err != nil {
		log.Info("load food err %v", err)
	}
	err = s.loadWeapon()
	if err != nil {
		log.Info("load weapon err %v", err)
	}
	return nil
}

func (s *ItemSystem) loadBuild() error {
	stuffBytes, err := os.ReadFile(s.Resources.GetPath("item/build.json"))
	if err != nil {
		return err
	}

	var items map[string]BuildItem
	if err := json.Unmarshal(stuffBytes, &items); err != nil {
		return err
	}
	for k, v := range items {
		v.Id = s.AllocId("build", k)
		s.ItemMap[v.Id] = &v
	}
	return nil
}

func (s *ItemSystem) loadFood() error {
	stuffBytes, err := os.ReadFile(s.Resources.GetPath("item/food.json"))
	if err != nil {
		return err
	}

	var items map[string]FoodItem
	if err := json.Unmarshal(stuffBytes, &items); err != nil {
		return err
	}
	for k, v := range items {
		v.Id = s.AllocId("food", k)
		s.ItemMap[v.Id] = &v
	}
	return nil
}

func (s *ItemSystem) loadWeapon() error {
	stuffBytes, err := os.ReadFile(s.Resources.GetPath("item/weapon.json"))
	if err != nil {
		return err
	}

	var items map[string]WeaponItem
	if err := json.Unmarshal(stuffBytes, &items); err != nil {
		return err
	}
	for k, v := range items {
		v.Id = s.AllocId("weapon", k)
		s.ItemMap[v.Id] = &v
	}
	return nil
}

func (s *ItemSystem) AllocId(category string, name string) int {
	id := s.idInc
	s.Path2Id[category+"."+name] = s.idInc
	s.idInc++
	return id
}

func (s *ItemSystem) Craft() int {
}
