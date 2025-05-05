package http

import (
	"context"
	"encoding/json"
	"my_test/log"
	"net"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	server *http.Server
	ip     = "0.0.0.0"
	port   = 8899
	pushCh = make(chan string)
)

type WsEvent struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

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
	eventRouter := make(map[string]func(string) string)
	NewWebSocketRouter(eventRouter)

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Error("read:", err)
				conn, err = upgrader.Upgrade(c.Writer, c.Request, nil)
				if err != nil {
					log.Error("reconnect failed:", err)
					break
				}
				log.Info("WebSocket reconnected")
				continue
			}
			log.Info("recv: %s", message)
			if eventRouter[string(message)] != nil {
				eventRouter[string(message)](string(message))
			}
		}
	}()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case message := <-pushCh:
			log.Info("push: %s", message)
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Error("write:", err)
			}
		}
	}
}

func PushEvent(event any) error {
	eventTypeName := reflect.TypeOf(event).String()
	pushEvent := WsEvent{
		Event: eventTypeName,
		Data:  event,
	}
	bytes, err := json.Marshal(pushEvent)
	if err != nil {
		log.Error("PushMsg json marshal err %v", err)
		return err
	}
	pushCh <- string(bytes)
	return nil
}

func Stop(ctx context.Context) {
	server.Shutdown(ctx)
}
