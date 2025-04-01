package world

import (
	"encoding/json"
	"my_test/log"
	"os"
)

type Celler struct {
	story *Story
}

func NewCeller(story *Story) *Celler {
	f := &Celler{
		story: story,
	}
	f.loadData()
	return f
}

func (f *Celler) PassBy() {

}

func (f *Celler) Explore() int {
	return 0
}

func (f *Celler) loadData() error {
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
