package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

func Info(format string, info ...any) {
	now := time.Now()
	_, file, line, _ := runtime.Caller(1)
	format = fmt.Sprintf("[%s] [%s:%d] %s\n", now.Format("2006-01-02 15:04:05"), filepath.Base(file), line, format)
	fmt.Printf(format, info...)
}

func Warn(format string, info ...any) {
	now := time.Now()
	_, file, line, _ := runtime.Caller(1)
	format = fmt.Sprintf("[%s] [%s:%d] WARN %s\n", now.Format("2006-01-02 15:04:05"), filepath.Base(file), line, format)
	fmt.Printf(format, info...)
}

func Error(format string, info ...any) {
	now := time.Now()
	_, file, line, _ := runtime.Caller(1)
	format = fmt.Sprintf("[%s] [%s:%d] ERROR %s\n", now.Format("2006-01-02 15:04:05"), filepath.Base(file), line, format)
	fmt.Printf(format, info...)
}
