package world

import "testing"

func TestPlayer(t *testing.T) {
	p := NewPlayer(nil)
	print(p.ToJson())
}
