package http

import "github.com/gin-gonic/gin"

func newGinRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	u := newUser()
	userRouterGroup := r.Group("/user")
	{
		userRouterGroup.POST("/user_register", u.Register)
	}
	return r
}
