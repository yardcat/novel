package world

import (
	"my_test/event"
	"strconv"
)

const (
	Rain = iota
	Snow
	Sun
)

type EnvSystem struct {
	Temperature int
	Weather     int
}

func NewEnvSystem() *EnvSystem {
	return &EnvSystem{
		Temperature: 0,
		Weather:     0,
	}
}

func (s *EnvSystem) RegisterEventHander(maps map[string]any) {
	maps["ChangeEnv"] = s.OnChangeEnv
}

func (s *EnvSystem) OnChangeEnv(event event.ChangeEnvEvent) {
	switch event.Type {
	case "temperature":
		value, _ := strconv.Atoi(event.Value)
		s.Temperature += value
	case "weather":
		s.Weather = string2Weather(event.Value)
	}
}

func string2Weather(s string) int {
	switch s {
	case "rain":
		return Rain
	case "snow":
		return Snow
	case "sun":
		return Sun
	default:
		return 0
	}
}
