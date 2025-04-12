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
	event := &world.StartCardEvent{}
	replay := w.story.ChallengeTower(event)
	jsonStr, err := json.Marshal(replay)
	if err != nil {
		log.Info("CardStart json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))
}

func (w *World) CardTurnStart(c *gin.Context) {
	event := &world.CardTurnStartEvent{}
	replay := w.story.ChallengeTower(event)
	jsonStr, err := json.Marshal(replay)
	if err != nil {
		log.Info("CardTurnStart json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))

}

func (w *World) CardTurnEnd(c *gin.Context) {
	start := &world.CardTurnEndEvent{}
	replay := w.story.ChallengeTower(start)
	jsonStr, err := json.Marshal(replay)
	if err != nil {
		log.Info("CardTurnEnd json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))

}
