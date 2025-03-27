package world

import (
	"my_test/log"
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
}

func NewRanch(name string) *Ranch {
	log.Info("new farm %s", name)
	f := &Ranch{
		AnimalMap: make(map[string]Animal),
		animals:   make(map[*Animal]int),
	}
	f.loadData()
	return f
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

func (f *Ranch) loadData() {
}
