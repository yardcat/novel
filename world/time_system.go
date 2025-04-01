package world

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

const (
	TICK   = 1000 * time.Millisecond
	SECOND = 0
	MINUTE = iota
	HOUR
	DAY
	MONTH
	YEAR
	TIME_UNIT_COUNT
)

type VirtualTime struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
}

type TimeSystem struct {
	callbacks   [TIME_UNIT_COUNT][]func()
	story       *Story
	ticker      *time.Ticker
	startTime   time.Time
	virtualTime VirtualTime
	timeRate    int
	daySecond   int
	monthSecond int
	yearSecond  int
}

func NewTimeSystem(s *Story) *TimeSystem {
	t := &TimeSystem{
		story:       s,
		ticker:      time.NewTicker(TICK),
		startTime:   time.Now(),
		virtualTime: VirtualTime{},
	}
	for i := SECOND; i < TIME_UNIT_COUNT; i++ {
		t.callbacks[i] = make([]func(), 0)
	}
	t.loadData()
	return t
}

func (t *TimeSystem) RegisterCallback(typ int, callback func()) {
}

func (t *TimeSystem) Tick() <-chan time.Time {
	t.updateVirtualTime()
	return t.ticker.C
}

func (t *TimeSystem) Stop() {
	t.ticker.Stop()
}

func (t *TimeSystem) GetVirtualTime(ts int) VirtualTime {
	seconds := int(ts)

	year := seconds / t.yearSecond
	seconds %= t.yearSecond
	month := seconds / t.monthSecond
	seconds %= t.monthSecond
	day := seconds / t.daySecond
	seconds %= t.daySecond
	hour := seconds / 3600
	seconds %= 3600
	minute := seconds / 60
	second := seconds % 60

	return VirtualTime{
		Year:   int(year),
		Month:  int(month),
		Day:    int(day),
		Hour:   int(hour),
		Minute: int(minute),
		Second: int(second),
	}
}

func (t *TimeSystem) GetRealTime(vt time.Duration) time.Duration {
	return vt / time.Duration(t.timeRate)
}

func (t *TimeSystem) updateVirtualTime() {
	stamp := int(t.startTime.Unix()-time.Now().Unix()) * t.timeRate
	current := t.GetVirtualTime(stamp)
	if current.Minute != t.virtualTime.Minute {
		t.triggerCallback(MINUTE)
	}
	if current.Hour != t.virtualTime.Hour {
		t.triggerCallback(HOUR)
	}
	if current.Day != t.virtualTime.Day {
		t.triggerCallback(DAY)
	}
	if current.Month != t.virtualTime.Month {
		t.triggerCallback(MONTH)
	}
	if current.Year != t.virtualTime.Year {
		t.triggerCallback(YEAR)
	}
	t.virtualTime = current
}

func (t *TimeSystem) triggerCallback(typ int) {
	for _, callback := range t.callbacks[typ] {
		callback()
	}
}

func (t *TimeSystem) loadData() error {
	file := t.story.GetResources().GetPath("time/config.json")
	jsonData, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("load config file err: %v", err)
		return err
	}

	s := struct {
		YearMonth int `json:"YearMonth"`
		MonthDay  int `json:"MonthDay"`
		DayHour   int `json:"DayHour"`
		TimeRate  int `json:"TimeRate"`
	}{}
	err = json.Unmarshal(jsonData, &s)
	if err != nil {
		log.Fatalf("unmarshal config file err: %v", err)
		return err
	}

	t.timeRate = s.TimeRate
	t.daySecond = s.DayHour * 3600
	t.monthSecond = t.daySecond * s.MonthDay
	t.yearSecond = t.monthSecond * s.YearMonth
	return nil
}
