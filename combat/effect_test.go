package combat

import (
	"errors"
	"fmt"
	"my_test/log"
	"reflect"
	"testing"

	"github.com/Knetic/govaluate"
)

type EffectTest struct {
	A *int
}

func (e *EffectTest) Foo() int {
	return *e.A
}

func (e EffectTest) Get(name string) (interface{}, error) {
	val := reflect.ValueOf(e)
	field := val.FieldByName(name)

	if !field.IsValid() {
		return nil, errors.New(fmt.Sprintf("can't find field '%s' in Export", name))
	}

	if field.Kind() == reflect.Ptr {
		return field.Elem().Interface(), nil
	}

	return field.Interface(), nil
}

func TestTower_EffectOn(t *testing.T) {
	c := 10
	e := EffectTest{
		A: &c,
	}
	exp, err := govaluate.NewEvaluableExpression("A + 1")
	parameters := make(map[string]interface{})
	parameters["c"] = &e
	if err != nil {
		log.Error("evaluate condition err %v", err)
	}
	ret, err := exp.Eval(e)
	if err != nil {
		log.Error("evaluate condition err %v", err)
	}
	ok := ret.(bool)
	print(ok)
}
