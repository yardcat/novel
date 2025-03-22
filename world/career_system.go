package world

import (
	"encoding/json"
	"my_test/career"
	"my_test/log"
	"my_test/util"
	"os"
	"path/filepath"
)

type CareerSystem struct {
	CareerProto [career.CareerTypeCount]*career.Career
	NameMap     map[string]int
	story       *Story
}

func (c *CareerSystem) GetCareer(name string) *career.Career {
	return c.CareerProto[c.NameMap[name]]
}

func NewCareerSystem(s *Story) *CareerSystem {
	c := &CareerSystem{
		NameMap: map[string]int{
			"doctor":     career.Doctor,
			"teacher":    career.Teacher,
			"programmer": career.Programmer},
		story: s,
	}
	c.loadData()
	return c
}

func (c *CareerSystem) loadData() error {
	careerPath := c.story.resources.GetPath("career")
	files, err := os.ReadDir(careerPath)
	if err != nil {
		log.Error("read career dir %v", err)
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			data, err := os.ReadFile(filepath.Join(careerPath, file.Name()))
			if err != nil {
				log.Error("load career file %s err %v", file.Name(), err)
				return err
			}

			idx := c.NameMap[util.GetFileNameWithoutExt(file.Name())]
			proto := &career.Career{}
			err = json.Unmarshal(data, &proto)
			if err != nil {
				log.Error("Unmarshal career file %s err %v", file.Name(), err)
				return err
			}
			c.CareerProto[idx] = proto
		}
	}
	return nil
}
