package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(level string) zerolog.Logger {
	l := zerolog.New(os.Stdout).With().Timestamp().Logger()
	switch level {
	case "debug":
		l = l.Level(zerolog.DebugLevel)
	case "info":
		l = l.Level(zerolog.InfoLevel)
	case "warn":
		l = l.Level(zerolog.WarnLevel)
	case "error":
		l = l.Level(zerolog.ErrorLevel)
	default:
		l = l.Level(zerolog.InfoLevel)
	}
	return l
}
