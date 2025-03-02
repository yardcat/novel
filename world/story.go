package world

import (
	"context"
	"encoding/json"
	"my_test/log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"
)

type Story struct {
	taskCh        chan EventTask
	ticker        *time.Ticker
	done          chan bool
	players       []*Player
	timeEvents    []TimeEventTask
	daysData      []DayData
	ItemSystem    *ItemSystem
	resources     *Resources
	eventHandlers map[string]any
}

type TimeEventTask struct {
	Time time.Duration
	EventTask
}

type EventTask struct {
	EventName string
	Event     string
}

type DayEvent struct {
	Time   time.Duration `json:"time"`
	Desc   string        `json:"description"`
	Action string        `json:"action"`
}

type DayData struct {
	Events []DayEvent
}

var (
	StoryInstance *Story
)

func GetStory() *Story {
	return StoryInstance
}

func NewStory() *Story {
	StoryInstance = &Story{resources: NewResources("island")}
	return StoryInstance
}

func (s *Story) Init() {
	s.taskCh = make(chan EventTask)
	s.done = make(chan bool)
	s.loadData()
	s.ItemSystem = NewItemSystem(s.resources)
	player := NewPlayer(s)
	s.players = append(s.players, player)
	s.RegisterEventHandler()
}

func (s *Story) RegisterEventHandler() {
	s.eventHandlers = make(map[string]any)
	NewEnvSystem().RegisterEventHander(s.eventHandlers)
	for _, player := range s.players {
		player.RegisterEventHander(s.eventHandlers)
	}
}

func (s *Story) PostEvent(name string, event string) {
	task := EventTask{
		EventName: name,
		Event:     event,
	}
	s.taskCh <- task
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
			s.handleTask(task)
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
				s.taskCh <- event.EventTask
			case <-s.done:
				return
			}
		}
	}
}

func (s *Story) handleTask(event EventTask) {
	action := event.EventName
	handler := s.eventHandlers[action]
	if handler == nil {
		log.Info("no handler for action %s", action)
		return
	}

	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func || handlerType.NumIn() != 1 {
		log.Info("handler %s is invalid", action)
		return
	}

	paramType := handlerType.In(0)
	paramValue := reflect.New(paramType).Interface()

	if err := json.Unmarshal([]byte(event.Event), paramValue); err != nil {
		log.Info("Failed to unmarshal event %s: %v", action, err)
		return
	}

	convertedValue := reflect.ValueOf(paramValue).Elem()
	handlerValue := reflect.ValueOf(handler)
	handlerValue.Call([]reflect.Value{convertedValue})
	log.Info("handle event done %s", action)
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

// TODO : use real id , not array index
func (s *Story) GetPlayerInfo(id string) string {
	idx, err := strconv.Atoi(id)
	if err != nil {
		log.Info("GetPlayerInfo invalid id %s", id)
	}
	player := s.players[idx]
	return player.ToJson()
}

func (s *Story) GetBag() string {
	player := s.players[0]
	return player.Bag.ToJson()
}

func (s *Story) loadData() error {
	if err := s.loadDays(); err != nil {
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
			}

			// Add time event to task queue
			eventStr, _ := json.Marshal(eventMap["params"])
			s.timeEvents = append(s.timeEvents, TimeEventTask{event.Time,
				EventTask{event.Action, string(eventStr)}})
			// record events
			events = append(events, event)
		}
		s.daysData[day] = DayData{Events: events}
		log.Info("Loaded day events done")
	}
	return nil
}

func (s *Story) HandleDayEvent(action string, params map[string]string) {
	handler := s.eventHandlers[action]
	if handler == nil {
		log.Info("no handler for action %s", action)
		return
	}
	handler.(func(map[string]string))(params)
	log.Info("handle day event done %s", action)
}
