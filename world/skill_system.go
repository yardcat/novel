package world

import (
	"encoding/json"
	"my_test/log"
	"os"
)

type SkillSystem struct {
	story *Story
}

func NewSkillSystem(story *Story) *SkillSystem {
	s := &SkillSystem{
		story: story,
	}
	s.loadData()
	return s
}

func (s *SkillSystem) loadData() error {
	file := s.story.GetResources().GetPath("skill/skill.json")
	jsonData, err := os.ReadFile(file)
	if err != nil {
		log.Error("load config file err: %v", err)
		return err
	}

	err = json.Unmarshal(jsonData, &s)
	if err != nil {
		log.Error("unmarshal config file err: %v", err)
		return err
	}
	return nil
}
