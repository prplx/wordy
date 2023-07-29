package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
	// axiomClient *axiom.Client
}

func New(out io.Writer, minLevel Level) *Logger {
	// go min version is 1.19 whereas max go version on Railway is 1.18
	// axiomClient, err := axiom.NewClient(
	// axiom.SetPersonalTokenConfig(os.Getenv("AXIOM_TOKEN"), os.Getenv("AXIOM_ORG_ID")),
	// )
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	return &Logger{out: out, minLevel: minLevel}
}

func (l *Logger) Info(message string, properties ...map[string]string) {
	l.print(LevelInfo, message, properties...)
}

func (l *Logger) Error(err error, properties ...map[string]string) {
	l.print(LevelError, err.Error(), properties...)
}

func (l *Logger) Fatal(err error, properties ...map[string]string) {
	l.print(LevelFatal, err.Error(), properties...)
	os.Exit(1)
}

// func (l *Logger) Ingest(events []axiom.Event) error {
// 	ctx := context.Background()
// 	_, err := l.axiomClient.IngestEvents(ctx, os.Getenv("AXIOM_DATASET"), events)
// 	return err
// }

func (l *Logger) print(level Level, message string, properties ...map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	var props map[string]string
	if len(properties) > 0 {
		props = properties[0]
	}

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: props,
	}

	// if level >= LevelError {
	// 	aux.Trace = string(debug.Stack())
	// }

	var line []byte

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// err = l.Ingest([]axiom.Event{{
	// 	ingest.TimestampField: aux.Time,
	// 	"message":             aux.Message,
	// 	"level":               aux.Level,
	// 	"properties":          aux.Properties,
	// 	"environment":         os.Getenv("APP_ENV"),
	// }})

	if err != nil {
		l.Error(err)
	}

	return l.out.Write(append(line, '\n'))
}

func (l *Logger) Write(message []byte) (int, error) {
	return l.print(LevelError, string(message), nil)
}
