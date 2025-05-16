package http

import (
	"my_test/pb"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var ()

func NewGinRouter(conn *grpc.ClientConn) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	setCors(r)

	p := newPlayer()
	playerRouterGroup := r.Group("/player")
	{
		playerRouterGroup.POST("/get_player_info", p.GetPlayerInfo)
		playerRouterGroup.POST("/get_bag", p.GetBag)
		playerRouterGroup.POST("/collect", p.Collect)
	}

	worldClient := pb.NewWorldClient(conn)
	w := newWorld(worldClient)
	worldRouterGroup := r.Group("/world")
	{
		worldRouterGroup.POST("/get_ui_info", w.GetUiInfo)
		worldRouterGroup.POST("/card_start", w.CardStart)
	}

	cardClient := pb.NewCardClient(conn)
	c := newCard(cardClient)
	cardRouterGroup := r.Group("/card")
	{
		cardRouterGroup.POST("/welcome", c.Welcome)
		cardRouterGroup.POST("/can_use", c.CanUse)
		cardRouterGroup.POST("/send_cards", c.SendCards)
		cardRouterGroup.POST("/discard_cards", c.DiscardCards)
		cardRouterGroup.POST("/end_turn", c.EndTurn)
		cardRouterGroup.POST("/choose_bonus", c.ChooseBonus)
		cardRouterGroup.POST("/next_floor", c.NextFloor)
		cardRouterGroup.POST("/use_potion", c.UsePotion)
		cardRouterGroup.POST("/discard_potion", c.DiscardPotion)
		cardRouterGroup.POST("/buy", c.Buy)
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
