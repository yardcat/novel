package http

import (
	"context"
	"my_test/log"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	server *http.Server
	ip     = "0.0.0.0"
	port   = 8899
)

func StartServer(ctx context.Context, cancel context.CancelFunc) {
	address := net.JoinHostPort(ip, strconv.Itoa(port))
	router := NewGinRouter()
	router.GET("/ws", webSocketHandler)
	server = &http.Server{Addr: address, Handler: router}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error("net error")
			cancel()
		}
	}()
}

func webSocketHandler(c *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error("upgrade:", err)
		return
	}
	defer conn.Close()
	eventRouter := make(map[string]any)
	NewWebSocketRouter(eventRouter)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Error("read:", err)
			break
		}
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Error("write:", err)
			break
		}
	}
}

func Stop(ctx context.Context) {
	server.Shutdown(ctx)
}
