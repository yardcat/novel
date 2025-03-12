package http

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewGinRouter() *gin.Engine {
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

	w := newWorld()
	worldRouterGroup := r.Group("/world")
	{
		worldRouterGroup.POST("/get_ui_info", w.GetUiInfo)
	}

	return r
}

func NewWebSocketRouter(eventRouter map[string]any) {
	eventRouter["get_player_info"] = func(c *gin.Context) {}

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
