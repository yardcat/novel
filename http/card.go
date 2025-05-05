package http

import (
	"encoding/json"
	"my_test/event"
	"my_test/log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Card struct {
}

func newCard() *Card {
	return &Card{}
}

func (w *Card) CardStart(c *gin.Context) {
	difficuty := c.PostForm("difficuty")
	event := &event.CardStartEvent{Difficulty: difficuty}
	reply := w.story.ChallengeTower(event)
	jsonStr, err := json.Marshal(reply)
	if err != nil {
		log.Info("CardStart json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))
}

func (w *Card) CardWelcome(c *gin.Context) {
	event := &event.CardWelcomeEvent{
		Event: c.PostForm("event"),
	}
	reply := w.story.CardWelcome(event)
	jsonStr, err := json.Marshal(reply)
	if err != nil {
		log.Info("CardChooseEvent json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))

}

func (w *Card) CardSendCards(c *gin.Context) {
	reply := &event.CardSendCardsReply{
		Status: "no_card",
	}
	cardsParam := c.PostForm("cards")
	if len(cardsParam) != 0 {
		strIdx := strings.Split(cardsParam, ",")
		targetIdx, _ := strconv.Atoi(c.PostForm("target"))
		ev := &event.CardSendCards{
			Cards:  make([]int, len(strIdx)),
			Target: targetIdx,
		}
		for i, v := range strIdx {
			ev.Cards[i], _ = strconv.Atoi(v)
		}
		reply = w.story.SendCards(ev)
	}
	jsonStr, err := json.Marshal(reply)
	if err != nil {
		log.Info("CardTurnStart json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))

}

func (w *Card) CardDiscardCards(c *gin.Context) {
	cardsParam := c.PostForm("cards")
	strIdx := strings.Split(cardsParam, ",")
	ev := &event.CardDiscardCards{
		Cards: make([]int, len(strIdx)),
	}
	for i, v := range strIdx {
		ev.Cards[i], _ = strconv.Atoi(v)
	}
	reply := w.story.DiscardCards(ev)
	jsonStr, err := json.Marshal(reply)
	if err != nil {
		log.Info("CardTurnStart json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))

}

func (w *Card) CardEndTurn(c *gin.Context) {
	start := &event.CardTurnEndEvent{}
	reply := w.story.EndTurn(start)
	jsonStr, err := json.Marshal(reply)
	if err != nil {
		log.Info("CardTurnEnd json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))

}

func (w *Card) CardNextFloor(c *gin.Context) {
	start := &event.CardNextFloorEvent{}
	reply := w.story.CardNextFloor(start)
	jsonStr, err := json.Marshal(reply)
	if err != nil {
		log.Info("CardNextFloor json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))
}

func (w *Card) CardEnterRoom(c *gin.Context) {
	start := &event.CardEnterRoomEvent{}
	reply := w.story.CardEnterRoom(start)
	jsonStr, err := json.Marshal(reply)
	if err != nil {
		log.Info("CardEnterRoom json marshal err %v", err)
	}
	c.JSON(200, string(jsonStr))
}
