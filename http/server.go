package http

import (
	"context"
	"net"
	"net/http"
	"strconv"

	"github.com/openimsdk/tools/utils/network"
)

func StartServer(ctx context.Context) {
	ip := "0.0.0.0"
	port := 9550
	address := net.JoinHostPort(network.GetListenIP(ip), strconv.Itoa(port))
	router := newGinRouter()
	server := http.Server{Addr: address, Handler: router}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			print("net error")
		}
	}()
}
