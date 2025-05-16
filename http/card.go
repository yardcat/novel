package http

import (
	"context"
	"encoding/json"
	"my_test/log"
	"my_test/pb"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type Card struct {
	client pb.CardClient
}

func newCard(ca pb.CardClient) *Card {
	return &Card{
		client: ca,
	}
}

func (w *Card) Welcome(c *gin.Context) {
	response, err := w.client.Welcome(c.Request.Context(), &pb.WelcomeRequest{
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

func (w *Card) CanUse(c *gin.Context) {
	response, err := w.client.CanUseCard(context.Background(), &pb.CanUseRequest{
		Card: cast.ToInt32(c.PostForm("card")),
	})
	if err != nil {
		log.Info("can use err %v", err)
	}

	jsonStr, err := json.Marshal(response)
	if err != nil {
		log.Info("CardTurnStart json marshal err %v", err)
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
	response, err := w.client.SendCard(context.Background(), &pb.SendCardRequest{
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
	response, err := w.client.DiscardCard(c.Request.Context(), &pb.DiscardCardRequest{
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
	response, err := w.client.EndTurn(c.Request.Context(), &pb.EndTurnRequest{})
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
	response, err := w.client.NextFloor(c.Request.Context(), &pb.NextFloorRequest{
		Room: cast.ToInt32(c.PostForm("room")),
	})
	if err != nil {
		log.Info("card next floor err %v", err)
	}

	jsonStr, err := json.Marshal(response)
	if err != nil {
		log.Info("CardNextFloor json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))
}

func (w *Card) ChooseBonus(c *gin.Context) {
	ev := pb.ChooseBonusRequest{}
	bonus := c.PostForm("bonus")
	err := json.Unmarshal([]byte(bonus), &ev.Bonus)
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

func (w *Card) UsePotion(c *gin.Context) {
	ev := pb.UsePotionRequest{
		Name: c.PostForm("name"),
	}
	response, err := w.client.UsePotion(c.Request.Context(), &ev)
	if err != nil {
		log.Info("card use potion err %v", err)
	}

	jsonStr, err := json.Marshal(response)
	if err != nil {
		log.Info("card use potion err %v", err)
	}
	c.JSON(200, jsonStr)
}

func (w *Card) DiscardPotion(c *gin.Context) {
	ev := pb.UsePotionRequest{
		Name: c.PostForm("name"),
	}
	response, err := w.client.UsePotion(c.Request.Context(), &ev)
	if err != nil {
		log.Info("card use potion err %v", err)
	}

	jsonStr, err := json.Marshal(response)
	if err != nil {
		log.Info("card use potion err %v", err)
	}
	c.JSON(200, jsonStr)
}

func (w *Card) Buy(c *gin.Context) {
	ev := pb.BuyRequest{
		Type: cast.ToInt32(c.PostForm("type")),
		Name: c.PostForm("name"),
	}
	response, err := w.client.Buy(c.Request.Context(), &ev)
	if err != nil {
		log.Info("card buy err %v", err)
	}

	jsonStr, err := json.Marshal(response)
	if err != nil {
		log.Info("card buy err %v", err)
	}
	c.JSON(200, jsonStr)
}
