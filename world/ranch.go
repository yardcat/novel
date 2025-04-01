package world

import (
	"encoding/json"
	"my_test/log"
	"os"
)

type Animal struct {
	Name        string
	GrowProcess int
}

type Ranch struct {
	AnimalMap map[string]Animal
	Width     int
	Height    int
	animals   map[*Animal]int
	story     *Story
}

func NewRanch(story *Story) *Ranch {
	log.Info("new farm")
	f := &Ranch{
		AnimalMap: make(map[string]Animal),
		animals:   make(map[*Animal]int),
		story:     story,
	}
	f.loadData()
	return f
}

func (f *Ranch) PassBy() {

}

func (f *Ranch) Explore() int {
	return 0
}

func (f *Ranch) CreateAnimal(name string) *Animal {
	plant := f.AnimalMap[name]
	return &plant
}

func (f *Ranch) Update() {
	for plant := range f.animals {
		plant.GrowProcess++
		if plant.GrowProcess >= 100 {
			log.Info("plant matured %s", plant.Name)
			f.RemoveAnimal(plant)
		}
	}
}

func (f *Ranch) AddAnimal(plant *Animal) {
	f.animals[plant]++
}

func (f *Ranch) RemoveAnimal(key *Animal) {
	if _, exist := f.animals[key]; exist {
		f.animals[key]--
		if f.animals[key] == 0 {
			delete(f.animals, key)
		}
	}
}

func (f *Ranch) Water() {
	for plant := range f.animals {
		plant.GrowProcess++
	}
}

func (f *Ranch) loadData() error {
	file := f.story.GetResources().GetPath("scene/ranch.json")
	jsonData, err := os.ReadFile(file)
	if err != nil {
		log.Error("load config file err: %v", err)
		return err
	}

	err = json.Unmarshal(jsonData, &f)
	if err != nil {
		log.Error("unmarshal config file err: %v", err)
		return err
	}
	return nil
}
