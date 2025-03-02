package http

import (
	"my_test/log"
	"my_test/world"

	"github.com/gin-gonic/gin"
)

type Player struct {
	Id string
}

func newPlayer() *Player {
	return &Player{
		Id: "0",
	}
}

func (u *Player) GetPlayerInfo(c *gin.Context) {
	log.Info("get player info")
	info := world.GetStory().GetPlayerInfo(u.Id)
	c.JSON(200, info)
}

func (u *Player) GetBag(c *gin.Context) {
	log.Info("get bag")
	bag := world.GetStory().GetBag()
	c.JSON(200, bag)
}

func (u *Player) Collect(c *gin.Context) {
	log.Info("collect %s", c.PostForm("items"))
	items := c.PostForm("items")
	world.GetStory().PostEvent("Collect", items)
}
