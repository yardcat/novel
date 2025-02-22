package http

import (
	"my_test/log"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id string
}

func newUser() *User {
	return &User{
		Id: "123",
	}
}

func (u *User) GetUserInfo(c *gin.Context) {
	log.Info("get user info %s", c.PostForm("key1"))

}
