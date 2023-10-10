package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type Interface interface {
	Debug(msg interface{}, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg interface{}, args ...interface{})
	Fatal(msg interface{}, args ...interface{})
}

var _ Interface = (*Logger)(nil)

type Logger struct {
	logger *zerolog.Logger
}

func New(level string) *Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	skipFrameCount := 3
	zerologger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()

	return &Logger{
		logger: &zerologger,
	}
}

// Debug implements Logger.
func (l *Logger) Debug(msg interface{}, args ...interface{}) {
	l.msg("debug", msg, args...)
}

// Error implements Logger.
func (l *Logger) Error(msg interface{}, args ...interface{}) {
	if l.logger.GetLevel() == zerolog.DebugLevel {
		l.Debug(msg, args...)
	}

	l.msg("error", msg, args...)
}

// Fatal implements Logger.
func (l *Logger) Fatal(msg interface{}, args ...interface{}) {
	l.msg("fatal", msg, args...)
	os.Exit(1)
}

// Info implements Logger.
func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(msg, args...)
}

// Warn implements Logger.
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(msg, args...)
}

func (l *Logger) log(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Info().Msg(msg)
	} else {
		l.logger.Info().Msgf(msg, args...)
	}
}

func (l *Logger) msg(level string, msg interface{}, args ...interface{}) {
	switch m := msg.(type) {
	case error:
		l.log(m.Error(), args...)
	case string:
		l.log(m, args...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", level, msg, m))
	}
}
