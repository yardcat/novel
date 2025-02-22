package world

import "strconv"

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

func (s *EnvSystem) OnChangeEnv(typ string, value string) {
	switch typ {
	case "temperature":
		v, _ := strconv.Atoi(value)
		s.Temperature += v
	case "weather":
		s.Weather = string2Weather(value)
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
