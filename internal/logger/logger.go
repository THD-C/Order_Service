package logger

import (
	"github.com/rs/zerolog"
	"os"
	"sync"
	"time"
)

var (
	globalLogger zerolog.Logger
	once         sync.Once
)

func Init() {
	once.Do(
		func() {
			zerolog.TimeFieldFormat = time.StampMilli
			globalLogger = zerolog.New(
				zerolog.ConsoleWriter{
					Out:        os.Stdout,
					TimeFormat: "2006-01-02 15:04:05.000",
				},
			).With().
				Timestamp().
				Caller().
				Logger()
		},
	)
}

func GetLogger() zerolog.Logger {
	return globalLogger
}
