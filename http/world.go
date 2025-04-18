package http

import (
	"encoding/json"
	"my_test/event"
	"my_test/log"
	"my_test/world"
	"strconv"
	"strings"

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
	difficuty := c.PostForm("difficuty")
	event := &event.CardStartEvent{Difficulty: difficuty}
	replay := w.story.ChallengeTower(event)
	jsonStr, err := json.Marshal(replay)
	if err != nil {
		log.Info("CardStart json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))
}

func (w *World) CardChooseEvent(c *gin.Context) {
	event := &event.CardChooseStartEvent{
		Event: c.PostForm("event"),
	}
	replay := w.story.CardChooseEvent(event)
	jsonStr, err := json.Marshal(replay)
	if err != nil {
		log.Info("CardChooseEvent json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))

}

func (w *World) CardSendCards(c *gin.Context) {
	cardsParam := c.PostForm("cards")
	strIdx := strings.Split(cardsParam, ",")
	ev := &event.CardSendCards{
		Cards: make([]int, len(strIdx)),
	}
	for i, v := range strIdx {
		ev.Cards[i], _ = strconv.Atoi(v)
	}
	replay := w.story.SendCards(ev)
	jsonStr, err := json.Marshal(replay)
	if err != nil {
		log.Info("CardTurnStart json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))

}

func (w *World) CardEndTurn(c *gin.Context) {
	start := &event.CardTurnEndEvent{}
	replay := w.story.EndTurn(start)
	jsonStr, err := json.Marshal(replay)
	if err != nil {
		log.Info("CardTurnEnd json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))

}
