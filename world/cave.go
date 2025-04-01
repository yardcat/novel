package world

import (
	"encoding/json"
	"my_test/log"
	"os"
)

type Cave struct {
	story *Story
}

func NewCave(story *Story) *Cave {
	f := &Cave{
		story: story,
	}
	f.loadData()
	return f
}

func (f *Cave) PassBy() {

}

func (f *Cave) Explore() int {
	return 0
}

func (f *Cave) loadData() error {
	file := f.story.GetResources().GetPath("scene/fishing.json")
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
