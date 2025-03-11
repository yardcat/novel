package util

import (
	"math/rand"
	"time"
)

var (
	seed   = time.Now().UnixNano()
	source = rand.NewSource(seed)
	r      = rand.New(source)
)

func GetRandomInt(max int) int {
	return r.Intn(max)
}
