package combat

import (
	"errors"
	"fmt"
	"my_test/log"
	"reflect"
	"testing"

	"github.com/Knetic/govaluate"
)

type Ea struct {
	v1 int
	v2 string
}

type EffectTest struct {
	A  *int
	ea *Ea
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

func TestTower_EffectNest(t *testing.T) {
	c := 10
	e := EffectTest{
		A: &c,
		ea: &Ea{
			v1: 1,
			v2: "test",
		},
	}
	exp, err := govaluate.NewEvaluableExpression("c.A")
	parameters := make(map[string]interface{})
	parameters["c"] = &e
	if err != nil {
		log.Error("evaluate condition err %v", err)
	}
	ret, err := exp.Evaluate(e)
	if err != nil {
		log.Error("evaluate condition err %v", err)
	}
	ok := ret.(bool)
	print(ok)
}
