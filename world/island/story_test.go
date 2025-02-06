package island

import (
	"testing"
	"time"
)

func TestStory_loadDays(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	str, _ := time.ParseInLocation("15:04:05", "01:06:06", loc)
	print(str.String())
}
