package http

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func newGinRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	setCors(r)
	u := newUser()
	userRouterGroup := r.Group("/user")
	{
		userRouterGroup.POST("/get_user_info", u.GetUserInfo)
	}
	return r
}

func setCors(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

}
