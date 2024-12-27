package logger

import (
	"github.com/rs/zerolog"
	"os"
	"sync"
)

var (
	globalLogger zerolog.Logger
	once         sync.Once
)

func Init() {
	once.Do(
		func() {
			globalLogger = zerolog.New(
				zerolog.ConsoleWriter{
					Out:        os.Stdout,
					TimeFormat: "2006-01-02 15:04:05",
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
