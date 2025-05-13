package libraryLogging

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/getsentry/sentry-go"

	"os"
	"runtime"
	"strings"
	"time"
)

var (
	Info  func(string, ...any)
	Warn  func(string, ...any)
	Error func(string, ...any)
	Panic func(string, ...any)
	Fatal func(string, ...any)
	Debug func(string, ...any)

	Sync = func() {}

	sentryEnabled bool
)

func InitSentry(sentryDSN string) {
	if sentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: sentryDSN,
		})
		if err != nil {
			Error("Sentry initialization failed: %v", err)
			return
		}
		Sync = func() {
			sentry.Flush(5 * time.Second) // wait up to 2 seconds
		}
		go func() {
			for {
				time.Sleep(time.Second * 5)
				Sync()
			}
		}()
	}
}

func init() {
	Info = logAndCapture("info", false)
	Warn = logAndCapture("warning", true)
	Error = logAndCapture("error", true)
	Debug = logAndCapture("debug", false)
	Panic = logAndCapture("panic", true)
	Fatal = logAndCapture("fatal", true)

	sentryDSN, _ := os.LookupEnv("LOGGING_SENTRY_DSN")
	InitSentry(sentryDSN)
}

func logAndCapture(level string, sendToSentry bool) func(string, ...any) {
	var sentryLevel sentry.Level
	var printFunc = func(v ...any) {
		fmt.Println(v[0])
	}
	var fileLine = func() string {
		_, file, line, _ := runtime.Caller(2)
		return fmt.Sprintf("%s:%d", file, line)
	}
	switch level {
	case "info":
		sentryLevel = sentry.LevelInfo
	case "warning":
		sentryLevel = sentry.LevelWarning
	case "error":
		sentryLevel = sentry.LevelDebug
	case "debug":
		sentryLevel = sentry.LevelDebug
	case "fatal":
	case "panic":
		sentryLevel = sentry.LevelFatal
		printFunc = func(v ...any) {
			Sync()
			fmt.Println(v[0])
		}
	}

	return func(msg string, args ...any) {
		timeStamp := time.Now()
		var formatted string
		if sentryLevel == sentry.LevelDebug {
			fileAndLine := fileLine()
			formatted = fmt.Sprintf("%d %s %s [%s] %s", timeStamp.UnixMicro(), timeStamp.Format("2006-01-02 15:04:05.000"), fileAndLine, strings.ToUpper(level), fmt.Sprintf(msg, args...))
		} else {
			formatted = fmt.Sprintf("%s [%s] %s", timeStamp.Format("2006-01-02 15:04:05.000"), strings.ToUpper(level), fmt.Sprintf(msg, args...))
		}
		if sendToSentry && sentryEnabled {
			stacktrace2 := sentry.NewStacktrace(3)
			exception := sentry.Exception{
				Type:       level,
				Value:      formatted,
				Stacktrace: stacktrace2,
			}

			event := sentry.NewEvent()
			event.Level = sentryLevel
			event.Message = formatted
			event.Exception = []sentry.Exception{exception}
			sentry.CaptureEvent(event)
		}
		printFunc(formatted)
	}
}

func Dumper(args ...any) {
	spew.Dump(args...)
}
