package http

import "github.com/gin-gonic/gin"

type User struct {
	Id string
}

func newUser() *User {
	return &User{
		Id: "123",
	}
}

func (u *User) GetUserInfo(c *gin.Context) {

}
