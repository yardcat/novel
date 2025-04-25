package push

import (
	"fmt"
	"my_test/event"
)

type pushFunc func(event any) error

var pusher pushFunc

func SetPusher(p pushFunc) {
	pusher = p
}

func PushEvent(event any) error {
	if pusher != nil {
		return pusher(event)
	}
	return nil
}

func PushAction(format string, args ...any) {
	if pusher != nil {
		action := fmt.Sprintf(format, args...)
		PushEvent(event.ActionUpdateEvent{Action: action})
	}
}
