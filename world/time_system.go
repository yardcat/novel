package world

import (
	"sync"
)

const (
	HOUR = iota
	DAY
)

var (
	instance *TimeSystem
	once     sync.Once
)

type TimeSystem struct {
}

func GetTimeSystem() *TimeSystem {
	once.Do(func() {
		instance = &TimeSystem{}
	})
	return instance
}
