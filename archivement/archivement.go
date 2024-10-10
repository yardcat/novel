package archivement

import (
	"fmt"
	"my_test/event"
)

type Archivement struct {
}

func GetArchivement() *Archivement {
	return &Archivement{}
}

func (a *Archivement) Init() {
	event.GetEventBus().AddEventListener(event.Die, a.onDie)
}

func (a *Archivement) onDie(ev *event.Event) {
	fmt.Println("archivement on die")
}
