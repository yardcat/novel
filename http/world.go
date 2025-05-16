package http

import (
	"encoding/json"
	"my_test/log"
	"my_test/pb"
	"my_test/world"

	"github.com/gin-gonic/gin"
)

type World struct {
	uiConfig *UiConfig
	client   pb.WorldClient
}

type UiConfig struct {
	Collectable []string
}

func newWorld(cl pb.WorldClient) *World {
	return &World{
		client: cl,
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
	response, err := w.client.StartCard(c.Request.Context(), &pb.StartCardRequest{
		Difficuty: c.PostForm("difficuty"),
	})
	if err != nil {
		log.Info("StartCard call err %v", err)
	}

	jsonStr, err := json.Marshal(response)
	if err != nil {
		log.Info("CardStart json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))
}
