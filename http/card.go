package http

import (
	"context"
	"encoding/json"
	"my_test/event"
	"my_test/log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type Card struct {
	client event.CardClient
}

func newCard(ca event.CardClient) *Card {
	return &Card{
		client: ca,
	}
}

func (w *Card) Welcome(c *gin.Context) {
	response, err := w.client.Welcome(c.Request.Context(), &event.WelcomeRequest{
		Event: c.PostForm("event"),
	})
	if err != nil {
		log.Info("card welcome err %v", err)
	}

	jsonStr, err := json.Marshal(response)
	if err != nil {
		log.Info("CardChooseEvent json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))

}

func (w *Card) SendCards(c *gin.Context) {
	strIdx := strings.Split(c.PostForm("cards"), ",")
	target := cast.ToInt32(c.PostForm("target"))
	cards := make([]int32, len(strIdx))
	for i, v := range strIdx {
		cards[i] = cast.ToInt32(v)
	}
	response, err := w.client.SendCard(context.Background(), &event.SendCardRequest{
		Cards:  cards,
		Target: target,
	})
	if err != nil {
		log.Info("card send err %v", err)
	}

	jsonStr, err := json.Marshal(response)
	if err != nil {
		log.Info("CardTurnStart json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))

}

func (w *Card) DiscardCards(c *gin.Context) {
	strIdx := strings.Split(c.PostForm("cards"), ",")
	cards := make([]int32, len(strIdx))
	for i, v := range strIdx {
		cards[i] = cast.ToInt32(v)
	}
	response, err := w.client.DiscardCard(c.Request.Context(), &event.DiscardCardRequest{
		Cards: cards,
	})
	if err != nil {
		log.Info("card discard err %v", err)
	}

	jsonStr, err := json.Marshal(response)
	if err != nil {
		log.Info("CardTurnStart json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))

}

func (w *Card) EndTurn(c *gin.Context) {
	response, err := w.client.EndTurn(c.Request.Context(), &event.EndTurnRequest{})
	if err != nil {
		log.Info("card end turn err %v", err)
	}

	jsonStr, err := json.Marshal(response)
	if err != nil {
		log.Info("CardTurnEnd json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))

}

func (w *Card) NextFloor(c *gin.Context) {
	response, err := w.client.NextFloor(c.Request.Context(), &event.NextFloorRequest{})
	if err != nil {
		log.Info("card next floor err %v", err)
	}

	jsonStr, err := json.Marshal(response)
	if err != nil {
		log.Info("CardNextFloor json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))
}

func (w *Card) EnterRoom(c *gin.Context) {
	response, err := w.client.EnterRoom(c.Request.Context(), &event.EnterRoomRequest{})
	if err != nil {
		log.Info("card next floor err %v", err)
	}

	jsonStr, err := json.Marshal(response)
	if err != nil {
		log.Info("CardEnterRoom json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))
}

func (w *Card) ChooseBonus(c *gin.Context) {
	ev := event.ChooseBonusRequest{}
	err := json.Unmarshal([]byte(c.PostForm("bonus")), &ev.Bonus)
	if err != nil {
		log.Info("card choose bonus err %v", err)
	}
	response, err := w.client.ChooseBonus(c.Request.Context(), &ev)
	if err != nil {
		log.Info("card choose bonus err %v", err)
	}

	jsonStr, err := json.Marshal(response)
	if err != nil {
		log.Info("card choose bonus err %v", err)
	}
	c.JSON(200, jsonStr)
}
