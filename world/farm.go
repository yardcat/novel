package world

import (
	"encoding/json"
	"my_test/log"
	"os"
)

type Plant struct {
	Name        string
	GrowProcess int
}

type Farm struct {
	PlantMap map[string]Plant
	Width    int
	Height   int
	plants   map[*Plant]int
	story    *Story
}

func NewFarm(story *Story) *Farm {
	log.Info("new farm")
	f := &Farm{
		PlantMap: make(map[string]Plant),
		plants:   make(map[*Plant]int),
		story:    story,
	}
	f.loadData()
	story.timeSystem.RegisterCallback(DAY, f.Update)

	return f
}

func (f *Farm) PassBy() {

}

func (f *Farm) Explore() int {
	return 0
}

func (f *Farm) CreatePlant(name string) *Plant {
	plant := f.PlantMap[name]
	return &plant
}

func (f *Farm) Update() {
	log.Info("update farm")
	for plant := range f.plants {
		plant.GrowProcess++
		if plant.GrowProcess >= 100 {
			log.Info("plant matured %s", plant.Name)
			f.RemovePlant(plant)
		}
	}
}

func (f *Farm) AddPlant(plant *Plant) {
	f.plants[plant]++
}

func (f *Farm) RemovePlant(plant *Plant) {
	if _, exist := f.plants[plant]; exist {
		f.plants[plant]--
		if f.plants[plant] == 0 {
			delete(f.plants, plant)
		}
	}
}

func (f *Farm) Water() {
	for plant := range f.plants {
		plant.GrowProcess++
	}
}

func (f *Farm) Fertilize() {
	for plant := range f.plants {
		plant.GrowProcess++
	}
}

func (f *Farm) loadData() error {
	file := f.story.GetResources().GetPath("scene/farm.json")
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
