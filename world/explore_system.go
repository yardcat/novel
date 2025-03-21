package world

import (
	"encoding/json"
	"my_test/log"
	"my_test/util"
	"os"
)

const (
	EXPLORE_MINE = 1
	EXPLORE_LAKE = iota
)

type Map struct {
	Name   string
	Width  int
	Height int
}

type Cord struct {
	X int
	Y int
}

type Explorable interface {
	Explore()
}

type Grid struct {
	Name       string
	Discovered bool
	Type       int
	Explorable bool
}

type ExploreResult struct {
}

type ExploreSystem struct {
	client   ExploreClient
	story    *Story
	mp       Map
	mapData  []Grid
	homeCord int
	MineMap  map[string]Mine
}

type ExploreClient interface {
	OnGridDiscovered(cord Cord)
}

func NewExploreSystem(story *Story) *ExploreSystem {
	e := &ExploreSystem{
		story: story,
	}
	e.loadMap()
	e.mapData = make([]Grid, e.mp.Width*e.mp.Height)
	e.fillMap()
	e.loadMine()
	return e
}

func (e *ExploreSystem) Explore(path []int) {
	cords := e.getGridFromPath(path)
	for _, cord := range cords {
		idx := e.cord2Index(cord)
		grid := &e.mapData[idx]
		grid.Discovered = true
		e.client.OnGridDiscovered(cord)
		e.exploreGrid(grid.Type)
	}
}

func (e *ExploreSystem) exploreGrid(typ int) {
	switch typ {
	case EXPLORE_MINE:
	}
}

func (e *ExploreSystem) getGridFromPath(path []int) []Cord {
	ret := []Cord{}
	start := e.index2Cord(e.homeCord)
	for i := 0; i < len(path); i++ {
		cord := e.index2Cord(path[i])
		dx := cord.X - start.X
		dy := cord.Y - start.Y
		if dx == 0 {
			for j := 0; j < util.Abs(dy); j++ {
				ret = append(ret, Cord{cord.X, cord.Y + util.Abs(dy)/dy*j})
			}
		} else if dy == 0 {
			for j := 0; j < util.Abs(dx); j++ {
				ret = append(ret, Cord{cord.X + util.Abs(dx)/dx*j, cord.Y})
			}
		}
		start = cord
	}
	return ret
}

func (e *ExploreSystem) index2Cord(index int) Cord {
	return Cord{
		X: index % e.mp.Width,
		Y: index / e.mp.Height,
	}
}

func (e *ExploreSystem) cord2Index(cord Cord) int {
	return cord.Y*e.mp.Height + cord.X
}

func (e *ExploreSystem) loadMap() error {
	mapFile := e.story.GetResources().GetPath("map/island.json")
	jsonData, err := os.ReadFile(mapFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, &e.mp)
	if err != nil {
		log.Error("load map err %v", err)
		return err
	}
	return nil
}

func (e *ExploreSystem) fillMap() {
	for i := 0; i < e.mp.Width; i++ {
		for j := 0; j < e.mp.Height; j++ {
			e.mapData[i*e.mp.Height+j] = Grid{
				Name:       "earth",
				Discovered: false,
			}
		}
	}
}

func (e *ExploreSystem) loadMine() error {
	mineFile := e.story.GetResources().GetPath("scene/mine.json")
	jsonData, err := os.ReadFile(mineFile)
	if err != nil {
		return err
	}

	var mines map[string]Mine
	if err := json.Unmarshal(jsonData, &mines); err != nil {
		log.Error("load mine err %v", err)
		return err
	}

	e.MineMap = mines
	return nil
}
