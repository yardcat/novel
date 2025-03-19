package world

import (
	"encoding/json"
	"os"
)

type PetSystem struct {
	PetMap map[string]Pet
	story  *Story
}

func NewPetSystem(story *Story) *PetSystem {
	p := &PetSystem{
		PetMap: make(map[string]Pet),
		story:  story,
	}
	p.loadData()
	return p
}

func (p *PetSystem) GetPet(name string) Pet {
	return p.PetMap[name]
}

func (p *PetSystem) loadData() error {
	if err := p.loadPets(); err != nil {
		return err
	}
	return nil
}

func (c *PetSystem) loadPets() error {
	jsonData, err := os.ReadFile(c.story.GetResources().GetPath("pet/pet.json"))
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &c.PetMap)
	if err != nil {
		return err
	}
	return nil
}
