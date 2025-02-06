package http

import (
	"context"
	"net"
	"net/http"
	"strconv"

	"github.com/openimsdk/tools/utils/network"
)

var (
	server *http.Server
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

func Stop(ctx context.Context) {
	server.Shutdown(ctx)
}
