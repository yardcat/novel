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
<<<<<<< HEAD
	taskCh        chan Task
=======
	taskCh        chan *EventTask
>>>>>>> 015ea37 (add reply task)
	ticker        *time.Ticker
	done          chan bool
	players       []*Player
	timeEvents    []TimeEventTask
	daysData      []DayData
	ItemSystem    *ItemSystem
	resources     *Resources
	eventHandlers map[string]any
}

type Task interface {
	GetName() string
	GetEvent() string
}

type TimeEventTask struct {
	Time time.Duration
	EventTask
}

type EventTask struct {
	EventName string
	Event     string
	ReplyCh   chan string
}

func (e *EventTask) GetName() string { return e.EventName }

func (e *EventTask) GetEvent() string { return e.Event }

type ReplyEventTask struct {
	Reply any
	EventTask
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
<<<<<<< HEAD
	s.taskCh = make(chan Task)
=======
	s.taskCh = make(chan *EventTask)
>>>>>>> 015ea37 (add reply task)
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
	s.taskCh <- &task
}

<<<<<<< HEAD
func (s *Story) PostReplyEvent(name string, event string, reply any) {
	task := ReplyEventTask{
		Reply: reply,
		EventTask: EventTask{
			EventName: name,
			Event:     event,
		},
	}
	s.taskCh <- &task
=======
func (s *Story) PostReplyEvent(name string, event string, callback func(string)) {
	replyCh := make(chan string)
	task := EventTask{
		EventName: name,
		Event:     event,
		ReplyCh:   replyCh,
	}
	s.taskCh <- &task
	reply := <-replyCh
	callback(reply)
>>>>>>> 015ea37 (add reply task)
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
				s.taskCh <- &event.EventTask
			case <-s.done:
				return
			}
		}
	}
}

<<<<<<< HEAD
func (s *Story) handleTask(event Task) {
	action := event.GetName()
=======
func (s *Story) handleTask(event *EventTask) {
	action := event.EventName
>>>>>>> 015ea37 (add reply task)
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

	if err := json.Unmarshal([]byte(event.GetEvent()), paramValue); err != nil {
		log.Info("Failed to unmarshal event %s: %v", action, err)
		return
	}

	convertedValue := reflect.ValueOf(paramValue).Elem()
	handlerValue := reflect.ValueOf(handler)
	results := handlerValue.Call([]reflect.Value{convertedValue})

	if event.ReplyCh != nil {
		if len(results) > 0 {
			replyValue := results[0].Interface()
			replyBytes, err := json.Marshal(replyValue)
			if err != nil {
				log.Info("Failed to marshal reply for event %s: %v", action, err)
				return
			}
			event.ReplyCh <- string(replyBytes)
		} else {
			log.Info("No reply from handler for event %s", action)
		}
	}

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
		var rawEvents map[string]any
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
				EventTask{EventName: event.Action, Event: string(eventStr)}})
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

func (s *Story) GetCollectable() []string {
	return s.ItemSystem.GetCollectable()
}
