package world

import (
	"encoding/json"
	"my_test/log"
	"os"
	"path/filepath"
)

type NpcSystem struct {
	NpcMap map[string]Npc
	story  *Story
}

func NewNpcSystem(story *Story) *NpcSystem {
	p := &NpcSystem{
		NpcMap: make(map[string]Npc),
		story:  story,
	}
	p.loadData()
	return p
}

func (p *NpcSystem) GetNpc(name string) Npc {
	return p.NpcMap[name]
}

func (p *NpcSystem) loadData() error {
	if err := p.loadNpcs(); err != nil {
		return err
	}
	return nil
}

func (c *NpcSystem) loadNpcs() error {
	npcPath := c.story.GetResources().GetPath("npc")
	files, err := os.ReadDir(npcPath)
	if err != nil {
		return err
	}

	npcs := make(map[string]*Npc)
	for _, file := range files {
		jsonData, err := os.ReadFile(filepath.Join(npcPath, file.Name()))
		if err != nil {
			log.Error("load npc err %v", err)
			continue
		}
		err = json.Unmarshal(jsonData, &npcs)
		if err != nil {
			log.Error("load npc err %v", err)
			continue
		}
		for k, v := range npcs {
			c.NpcMap[k] = *v
		}
	}

	return nil
}
