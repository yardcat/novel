package world

import (
	"encoding/json"
	"my_test/log"
	"os"
)

type Forest struct {
	story *Story
}

func NewForest(story *Story) *Forest {
	f := &Forest{
		story: story,
	}
	f.loadData()
	return f
}

func (f *Forest) PassBy() {

}

func (f *Forest) Explore() int {
	return 0
}

func (f *Forest) loadData() error {
	file := f.story.GetResources().GetPath("scene/forest.json")
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
