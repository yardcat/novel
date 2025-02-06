package log

import (
	"fmt"
	"time"
)

func Info(format string, info ...any) {
	now := time.Now()
	format = fmt.Sprintf("[%s] %s\n", now.Format("2006-01-02 15:04:05"), format)
	fmt.Printf(format, info...)
}
