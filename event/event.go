package event

const (
	Kill = iota
	Die
)

type Event struct {
	Type   int
	Bundle map[string]any
}

type EventBus struct {
	Handlers map[int][]func(event *Event)
}

var eventBus *EventBus = &EventBus{
	Handlers: make(map[int][]func(event *Event)),
}

func GetEventBus() *EventBus {
	return eventBus
}

func (e *EventBus) AddEventListener(event_type int, handler func(event *Event)) {
	e.Handlers[event_type] = append(e.Handlers[event_type], handler)
}

func (e *EventBus) RemoveEventListener(event *Event, handler func(event *Event)) {
	e.Handlers[event.Type] = nil
}

func (e *EventBus) OnEvent(typ int, bundle map[string]any) {
	event := &Event{Type: typ, Bundle: bundle}
	e.dispatch(event)
}

func (e *EventBus) dispatch(event *Event) {
	for _, handler := range e.Handlers[event.Type] {
		handler(event)
	}
}
