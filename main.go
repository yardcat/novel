package main

import (
	"my_test/scene"
	"my_test/user"
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
	group := []*user.Player{
		user.NewPlayer(1, "player1"),
	}
	scene := scene.CreateScene(5000001, group)
	scene.DoCombat()
}

func loadConfig() {
}

