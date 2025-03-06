package http

import (
	"context"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/openimsdk/tools/utils/network"
)

var (
	server   *http.Server
	upgrader = websocket.Upgrader{} // use default options
)

func StartServer(ctx context.Context, cancel context.CancelFunc) {
	ip := "0.0.0.0"
	port := 8899
	address := net.JoinHostPort(network.GetListenIP(ip), strconv.Itoa(port))
	router := newGinRouter()
	server = &http.Server{Addr: address, Handler: router}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			print("net error")
			cancel()
		}
	}()
}

func StartWebSocketServer(ctx context.Context, cancel context.CancelFunc) {
	ip := "0.0.0.0"
	port := 8899
	address := net.JoinHostPort(network.GetListenIP(ip), strconv.Itoa(port)+"/ws")

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			print("upgrade:", err)
			return
		}
		defer conn.Close()

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				print("read:", err)
				break
			}
			print("recv: ", message)
			err = conn.WriteMessage(messageType, message)
			if err != nil {
				print("write:", err)
				break
			}
		}
	})

	go func() {
		err := http.ListenAndServe(address, nil)
		if err != nil && err != http.ErrServerClosed {
			print("websocket net error")
			cancel()
		}
	}()
}

func Stop(ctx context.Context) {
	server.Shutdown(ctx)
}
