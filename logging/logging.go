package libraryLogging

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/getsentry/sentry-go"
	"log"
	"os"
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

	Sync func() error

	sentryEnabled bool
)

func init() {
	sentryDSN, ok := os.LookupEnv("LOGGING_SENTRY_DSN")
	log.SetFlags(log.LstdFlags | log.Llongfile | log.Lmicroseconds)

	if ok && sentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: sentryDSN,
		})
		if err != nil {
			log.Printf("Sentry initialization failed: %v", err)
		} else {
			sentryEnabled = true
			//log.Println("Sentry initialized")
		}
		go func() {
			for {
				time.Sleep(time.Second * 5)
				Sync()
			}
		}()
	}

	Info = logAndCapture("info", false)
	Warn = logAndCapture("warning", true)
	Error = logAndCapture("error", true)
	Debug = logAndCapture("debug", false)
	Panic = logAndCapture("panic", true)
	Fatal = logAndCapture("fatal", true)

	Sync = func() error {
		if sentryEnabled {
			sentry.Flush(5 * time.Second) // wait up to 2 seconds
		}
		return nil
	}
}

func logAndCapture(level string, sendToSentry bool) func(string, ...any) {
	var sentryLevel sentry.Level
	var printFunc = func(v ...any) {
		log.Print(v[0])
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
			log.Fatal(v[0])
		}
	}

	return func(msg string, args ...any) {
		formatted := formatLog(level, msg, args...)
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

func formatLog(level, msg string, args ...any) string {
	return fmt.Sprintf("[%s] %s", strings.ToUpper(level), fmt.Sprintf(msg, args...))
}

func Dumper(args ...any) {
	spew.Dump(args...)
}
