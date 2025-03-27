package world

import (
	"my_test/log"
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
}

func NewFarm(story *Story) *Farm {
	log.Info("new farm")
	f := &Farm{
		PlantMap: make(map[string]Plant),
		plants:   make(map[*Plant]int),
	}
	f.loadData()

	return f
}

func (f *Farm) CreatePlant(name string) *Plant {
	plant := f.PlantMap[name]
	return &plant
}

func (f *Farm) Update() {
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

func (f *Farm) loadData() {
}
