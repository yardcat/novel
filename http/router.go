package http

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	p *Player
	w *World
)

func NewGinRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	setCors(r)

	p = newPlayer()
	playerRouterGroup := r.Group("/player")
	{
		playerRouterGroup.POST("/get_player_info", p.GetPlayerInfo)
		playerRouterGroup.POST("/get_bag", p.GetBag)
		playerRouterGroup.POST("/collect", p.Collect)
	}

	w = newWorld()
	worldRouterGroup := r.Group("/world")
	{
		worldRouterGroup.POST("/get_ui_info", w.GetUiInfo)
		worldRouterGroup.POST("/card_start", w.CardStart)
		worldRouterGroup.POST("/card_welcome", w.CardWelcome)
		worldRouterGroup.POST("/send_cards", w.CardSendCards)
		worldRouterGroup.POST("/discard_cards", w.CardDiscardCards)
		worldRouterGroup.POST("/end_turn", w.CardEndTurn)
		worldRouterGroup.POST("/card_next_floor", w.CardEndTurn)
		worldRouterGroup.POST("/card_enter_room", w.CardEndTurn)
	}

	c = newCard()
	cardRouterGroup := r.Group("/card")
	{
		cardRouterGroup.POST("/get_ui_info", w.GetUiInfo)
		cardRouterGroup.POST("/card_start", w.CardStart)
		cardRouterGroup.POST("/card_welcome", w.CardWelcome)
		cardRouterGroup.POST("/send_cards", w.CardSendCards)
		cardRouterGroup.POST("/discard_cards", w.CardDiscardCards)
		cardRouterGroup.POST("/end_turn", w.CardEndTurn)
		cardRouterGroup.POST("/card_next_floor", w.CardEndTurn)
		cardRouterGroup.POST("/card_enter_room", w.CardEndTurn)
	}

	return r
}

func NewWebSocketRouter(eventRouter map[string]func(string) string) {

}

func setCors(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "Accept-Encoding"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

}
