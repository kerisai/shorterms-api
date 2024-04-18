package config

import (
	"errors"
	stdlog "log"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidLogLevel = errors.New("invalid log level")

	HttpLogger zerolog.Logger
)

func logLevel(c Config) (level zerolog.Level, err error) {
	switch c.LogLevel {
	case "panic":
		return zerolog.PanicLevel, nil
	case "fatal":
		return zerolog.FatalLevel, nil
	case "error":
		return zerolog.ErrorLevel, nil
	case "warn":
		return zerolog.WarnLevel, nil
	case "info":
		return zerolog.InfoLevel, nil
	case "debug":
		return zerolog.DebugLevel, nil
	case "trace":
		return zerolog.TraceLevel, nil
	default:
		return -1, ErrInvalidLogLevel
	}
}

func configureLogger(c Config) {
	level, err := logLevel(c)
	if err != nil {
		stdlog.Fatal("failed to configure logger: ", err)
	}
	
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.MessageFieldName = "msg"

	mw := zerolog.MultiLevelWriter(os.Stdout)
	logger := zerolog.New(mw).With().Timestamp().Caller().Stack().Logger()

	log.Logger = logger

	HttpLogger = zerolog.New(mw).With().Timestamp().Logger()
}