package main

import (
	"context"
	"my_test/http"
	"my_test/push"
	"my_test/scene"
	"my_test/world"
	"os"
	"os/signal"
	"syscall"
)

type Message struct {
	cmd  string
	text string
}

type Config struct {
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	loadConfig()
	scene.LoadDataFromCSV()

	world := world.NewStory()
	world.Init()
	go world.Start(ctx)

	http.StartServer(ctx, cancel)
	push.SetPusher(http.PushEvent)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)

	defer cancel()
	select {
	case <-sigCh:
		http.Stop(ctx)
	}
	close((sigCh))
}

func loadConfig() {
}
