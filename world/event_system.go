package world

type EventListener interface {
	OnEvent(eventType int, data map[string]interface{})
}
