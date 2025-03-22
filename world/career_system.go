package world

import (
	"encoding/json"
	"io/ioutil"
	"my_test/career"
	"path/filepath"
)

type CareerSystem struct {
	CareerProto [career.CareerTypeCount]career.Career
	NameMap     map[string]int
}

func NewCareerSystem() *CareerSystem {
	c := &CareerSystem{
		NameMap: map[string]int{
			"doctor":     career.Doctor,
			"teacher":    career.Teacher,
			"programmer": career.Programmer},
	}
	c.loadData()
	return c
}

func (c *CareerSystem) loadData() {
	files, err := ioutil.ReadDir("career")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			data, err := ioutil.ReadFile(filepath.Join("career", file.Name()))
			if err != nil {
				panic(err)
			}

			var careerData career.Career
			err = json.Unmarshal(data, &careerData)
			if err != nil {
				panic(err)
			}

			c.CareerProto[careerData.Type] = careerData
		}
	}
}
