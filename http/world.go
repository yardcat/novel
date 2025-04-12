package http

import (
	"encoding/json"
	"my_test/log"
	"my_test/world"

	"github.com/gin-gonic/gin"
)

type World struct {
	story    *world.Story
	uiConfig *UiConfig
}

type UiConfig struct {
	Collectable []string
}

func newWorld() *World {
	return &World{
		story: world.GetStory(),
	}
}

func (w *World) GetUiConfig() *UiConfig {
	if w.uiConfig == nil {
		w.uiConfig = &UiConfig{
			Collectable: world.GetStory().GetCollectable(),
		}
	}
	return w.uiConfig
}

func (w *World) GetUiInfo(c *gin.Context) {
	jsonStr, err := json.Marshal(w.GetUiConfig())
	if err != nil {
		log.Info("GetUiInfo json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))
}

func (w *World) CardStart(c *gin.Context) {
	w.story.ChallengeTower()
	jsonStr, err := json.Marshal()
	if err != nil {
		log.Info("GetUiInfo json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))
}
