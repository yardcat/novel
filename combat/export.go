package combat

import (
	"errors"
	"fmt"
	"my_test/log"
	"reflect"
)

type Export struct {
	Life          *int
	Strength      *int
	Defense       *int
	Energy        *int
	InitEnergy    *int
	CardCount     *int
	InitCardCount *int
}

func (e Export) Get(name string) (any, error) {
	val := reflect.ValueOf(e)
	field := val.FieldByName(name)

	if !field.IsValid() {
		return nil, errors.New(fmt.Sprintf("can't find field '%s'", name))
	}

	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return nil, errors.New(fmt.Sprintf("field '%s is nil", name))
		}
		return field.Elem().Interface(), nil
	}

	return nil, errors.New(fmt.Sprintf("field '%s is not pointer", name))
}

func (e *Export) Modify(key string, value int) bool {
	val := reflect.ValueOf(e).Elem()
	field := val.FieldByName(key)

	if !field.IsValid() {
		log.Error("Field '%s' not found in Export", key)
		return false
	}

	if !field.CanSet() {
		log.Error("Cannot set field '%s'. It is unexported or not addressable", key)
		return false
	}

	if field.Kind() == reflect.Ptr {
		if field.Type().Elem().Kind() != reflect.Int {
			log.Error("Field '%s' is not a pointer to an int", key)
			return false
		}
		if field.IsNil() {
			newValue := value
			field.Set(reflect.ValueOf(&newValue))
		} else {
			currentValue := field.Elem().Int()
			newValue := currentValue + int64(value)
			field.Elem().SetInt(newValue)
		}
	} else {
		log.Error("Field '%s' is not an *int", key)
		return false
	}
	return true
}
