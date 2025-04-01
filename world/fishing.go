package world

import (
	"encoding/json"
	"my_test/log"
	"my_test/util"
	"os"
)

type Fish struct {
	Name string
	Rate int
}

type Fishing struct {
	fishMap map[string]Fish `json:"fishMap"`
	fishes  []string
	story   *Story
}

func NewFish(story *Story) *Fishing {
	f := &Fishing{
		fishMap: make(map[string]Fish),
		story:   story,
	}
	f.loadData()
	return f
}

func (f *Fishing) PassBy() {

}

func (f *Fishing) Explore() int {
	return 0
}

func (f *Fishing) Fish() {
	keys := make([]string, 0, len(f.fishMap))
	for k := range f.fishMap {
		keys = append(keys, k)
	}
	idx := util.GetRandomInt(len(keys))
	rate := util.GetRandomInt(100)
	if rate < f.fishMap[keys[idx]].Rate {
		f.fishes = append(f.fishes, keys[idx])
	} else {
		log.Info("fish nothing")
	}
}

func (f *Fishing) loadData() error {
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
