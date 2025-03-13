package push

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
