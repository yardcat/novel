package util

import (
	"fmt"
	"testing"
)

func TestFormatString(t *testing.T) {
	ret := FormatString("Hello {name}, you are {age} years old.", map[string]any{
		"name": "Alice",
		"age":  30,
	})
	fmt.Println(ret)
}
