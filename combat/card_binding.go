package combat

import (
	"my_test/log"
	"reflect"
)

func (t *Tower) registerCardBindings() {
	t.cardBindingMap = make(map[string]reflect.Type)
	t.cardBindingMap["fuhua"] = reflect.TypeOf(fuhua{})
}

type CardBinding interface {
	Use(*Tower, *Card)
	GetCost() int
	GetDamage() int
	CanPlay() bool
}

type fuhua struct {
}

func (f *fuhua) Use(t *Tower, card *Card) {
	t.regiserTimingCallback(TIMING_USE_CARD, func(c *Card) {
		log.Info("fuhua take effect on %s", c.Name)
	})
}

func (f *fuhua) GetCost() int {
	return 10
}

func (f *fuhua) GetDamage() int {
	return 0
}

func (f *fuhua) Modify() {

}

func (f *fuhua) CanPlay() bool {
	return false
}
