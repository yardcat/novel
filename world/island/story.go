package island

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Story struct {
	tasks      chan interface{}
	ticker     *time.Ticker
	done       chan bool
	players    []*Player
	timeEvents []*TimeEvent
}

type TimeEvent struct {
	time   time.Time
	action func()
}

type TimeEventTask struct {
	action func()
}

type StuffData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Energy      int    `json:"energy"`
}

type DayData struct {
	Message struct {
		Msg  string `json:"msg"`
		Type string `json:"type"`
		Time string `json:"time"`
	} `json:"message"`
	Choices []struct {
		Description string `json:"description"`
		Action      string `json:"action"`
	} `json:"choices"`
	End bool `json:"end"`
}

func (s *Story) Init() {
	s.tasks = make(chan interface{})
	s.done = make(chan bool)
	s.loadData()
	player := NewPlayer()
	s.players = append(s.players, player)
}

func (s *Story) Start() {
	s.ticker = time.NewTicker(1 * time.Minute)

	go s.runTimeline()

	// Start main loop
	for {
		select {
		case <-s.done:
			s.ticker.Stop()
			return
		case <-s.ticker.C:
			if err := s.update(); err != nil {
				continue
			}
		case task := <-s.tasks:
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

	// 按时间排序检查和执行事件
	for _, event := range s.timeEvents {
		// 等待直到事件时间到达
		waitTime := time.Until(event.time)
		if waitTime > 0 {
			select {
			case <-time.After(waitTime):
				// 将事件动作包装成任务发送到任务通道
				s.tasks <- TimeEventTask{action: event.action}
			case <-s.done:
				// 如果收到停止信号，终止事件循环
				return
			}
		}
	}
}

func (s *Story) handleTask(task any) error {
	switch t := task.(type) {
	case TimeEventTask:
		t.action()
	}
	return nil
}

func (s *Story) loadData() error {
	// Load stuff data
	stuffBytes, err := os.ReadFile("data/stuff/build.json")
	if err != nil {
		return err
	}

	stuffData := make(map[string]StuffData)
	if err := json.Unmarshal(stuffBytes, &stuffData); err != nil {
		return err
	}

	// Load day data
	dayFiles, err := filepath.Glob("data/days/day*.json")
	if err != nil {
		return err
	}

	dayData := make(map[string]DayData)
	for _, file := range dayFiles {
		dayBytes, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		var day DayData
		if err := json.Unmarshal(dayBytes, &day); err != nil {
			return err
		}

		dayData[filepath.Base(file)] = day
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
	for _, p in range s.players {
		if p.Id == id {
			return p.GetInfoAsJSON()
		}
	}
	return nil
}
