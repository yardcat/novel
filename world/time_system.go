package world

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

const (
	TICK   = 1000 * time.Millisecond
	MINUTE = 0
	HOUR   = iota
	DAY
	MONTH
	YEAR
	SECOND
)

type VirtualTime struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	SECOND int
}

type TimeSystem struct {
	callbacks map[int]func()
	story     *Story
	tiker     *time.Ticker
	vtime     VirtualTime
	timeRate  int `json:"TimeRate"`
	dayHour   int `json:"DayHour"`
	monthDay  int `json:"MonthDay"`
	yearMonth int `json:"YearMonth"`
}

func NewTimeSystem(s *Story) *TimeSystem {
	t := &TimeSystem{
		story: s,
		tiker: time.NewTicker(TICK),
	}
	t.loadData()
	t.vtime = VirtualTime{0, 0, 0, 0, 0, 0}
	return t
}

func (t *TimeSystem) RegisterCallback(typ int, callback func()) {
}

func (t *TimeSystem) Tick() <-chan time.Time {
	t.UpdateTime()
	return t.tiker.C
}

func (t *TimeSystem) Stop() {
	t.tiker.Stop()
}

func (t *TimeSystem) UpdateTime() {
}

func (t *TimeSystem) loadData() {
	file, err := os.Open("time/config.json")
	if err != nil {
		log.Fatalf("无法打开配置文件: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(t)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
}
