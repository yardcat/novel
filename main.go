package main

import (
	"context"
	"my_test/http"
	"my_test/scene"
	"my_test/world"
	"my_test/world/island"
)

type Message struct {
	cmd  string
	text string
}

type Config struct {
}

func main() {
	loadConfig()
	scene.LoadDataFromCSV()

	ctx := context.Background()
	http.StartServer(ctx)

	world := &island.Story{}
	world.Init()
	world.Start()
}

func GetCurrentWorld() world.World {
	return nil

}

func loadConfig() {
}
