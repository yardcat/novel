package util

import (
	"fmt"
	"strconv"
)

const (
	Int = iota
	Float
	Percent
)

type Value struct {
	Number string
	Type   int
}

func NewValue(t int, n string) Value {
	return Value{Type: t, Number: n}
}

func (value Value) Int() int {
	if value.Type != Int {
		err := fmt.Sprintf("%s not int", value.Number)
		panic(err)
	}
	ret, _ := strconv.Atoi(value.Number)
	return ret
}

func (value Value) Float() float64 {
	if value.Type != Float {
		err := fmt.Sprintf("%s not float", value.Number)
		panic(err)
	}
	ret, _ := strconv.ParseFloat(value.Number, 64)
	return ret
}
