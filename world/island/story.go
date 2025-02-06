package island

import (
	"context"
	"encoding/json"
	"my_test/log"
	"my_test/world"
	"os"
	"path/filepath"
	"time"
)

type Story struct {
	taskCh     chan interface{}
	ticker     *time.Ticker
	done       chan bool
	players    []*Player
	timeEvents []TimeEventTask
	daysData   []DayData
	stuffData  map[string]StuffData
	resources  world.Resources
}

type TimeEventTask struct {
	Time     time.Duration
	Callback func()
}

type StuffData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Energy      int    `json:"energy"`
}

type DayEvent struct {
	Time   time.Duration     `json:"time"`
	Desc   string            `json:"description"`
	Action string            `json:"action"`
	Params map[string]string `json:"params"`
}

type DayData struct {
	Events []DayEvent
}

func (s *Story) Init() {
	s.taskCh = make(chan interface{})
	s.done = make(chan bool)
	s.resources.Init("island")
	s.loadData()
	player := NewPlayer()
	s.players = append(s.players, player)
}

func (s *Story) Start(ctx context.Context) {
	log.Info("start story")
	s.ticker = time.NewTicker(1 * time.Minute)

	go s.runTimeline()

	for {
		select {
		case <-ctx.Done():
		case <-s.done:
			s.ticker.Stop()
			return
		case <-s.ticker.C:
			if err := s.update(); err != nil {
				continue
			}
		case task := <-s.taskCh:
			if err := s.handleTask(task); err != nil {
				continue
			}
		}
	}
}

func (s *Story) runTimeline() {
	if len(s.timeEvents) == 0 {
		return
	}
	for _, event := range s.timeEvents {
		waitTime := event.Time
		if waitTime > 0 {
			select {
			case <-time.After(waitTime):
				s.taskCh <- event
			case <-s.done:
				return
			}
		}
	}
}

func (s *Story) handleTask(task any) error {
	switch t := task.(type) {
	case TimeEventTask:
		t.Callback()
	}
	return nil
}

func (s *Story) update() error {
	for _, player := range s.players {
		player.Update()
	}
	return nil
}

func (s *Story) Stop() {
	s.done <- true
}

func (s *Story) GetUserInfo(id string) string {
	return ""
}

func (s *Story) loadData() error {
	if err := s.loadStuff(); err != nil {
		return err
	}
	if err := s.loadDays(); err != nil {
		return err
	}
	return nil
}

func (s *Story) loadStuff() error {
	stuffBytes, err := os.ReadFile(s.resources.GetPath("stuff/build.json"))
	if err != nil {
		return err
	}

	s.stuffData = make(map[string]StuffData)
	if err := json.Unmarshal(stuffBytes, &s.stuffData); err != nil {
		return err
	}
	return nil
}

func (s *Story) loadDays() error {
	dayFiles, err := filepath.Glob(s.resources.GetPath("days/day*.json"))
	if err != nil {
		return err
	}

	s.daysData = make([]DayData, len(dayFiles))
	timeStart, _ := time.Parse("15:04:05", "00:00:00")
	for day, file := range dayFiles {
		dayBytes, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		var rawEvents map[string]interface{}
		if err := json.Unmarshal(dayBytes, &rawEvents); err != nil {
			return err
		}

		var events []DayEvent
		for k, v := range rawEvents {
			eventMap := v.(map[string]any)
			timePoint, _ := time.Parse("15:04:05", k)
			description, ok := eventMap["description"].(string)
			if !ok {
				description = ""
			}
			event := DayEvent{
				Time:   timePoint.Sub(timeStart),
				Desc:   description,
				Action: eventMap["action"].(string),
				Params: make(map[string]string),
			}
			if params, ok := eventMap["params"].(map[string]any); ok {
				for k, v := range params {
					event.Params[k] = v.(string)
				}
			}

			// Add time event to task queue
			s.timeEvents = append(s.timeEvents, TimeEventTask{event.Time,
				func() { s.HandleDayEvent(event.Action, event.Params) }})
			// record events
			events = append(events, event)
		}
		s.daysData[day] = DayData{Events: events}
		log.Info("Loaded day events done")
	}
	return nil
}

func (s *Story) HandleDayEvent(action string, params map[string]string) {
	switch action {
	case "SendMessage":
		log.Info(params["value"])
	case "Bonus":
		item := s.stuffData[params["type"]].Item
		s.players[0].Bag.Add(item)

		log.Info(params["value"])
	case "ChangeEnv":
		log.Info(params["type"])
		log.Info(params["value"])
	case "ChangeStatus":
		log.Info(params["type"])
		log.Info(params["value"])
	}
	log.Info("handle day event done")
}
