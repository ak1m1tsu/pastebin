package log

import (
	"io"
	"log"

	"github.com/rs/zerolog"
)

const (
	fiedsKey       = "data"
	skipFrameCount = 1
)

type Level int

const (
	Trace Level = iota - 1
	Debug
	Info
	Warn
	Error
	Fatal
)

func Stol(level string) Level {
	switch level {
	case "trace":
		return Trace
	case "debug":
		return Debug
	case "info":
		return Info
	case "warn":
		return Warn
	case "error":
		return Error
	case "fatal":
		return Fatal
	default:
		return Info
	}
}

type F struct {
	Key   string
	Value any
}

func (f F) MarshalZerologObject(e *zerolog.Event) {
	e.Any(f.Key, f.Value)
}

type FF []F

func (fields FF) MarshalZerologArray(a *zerolog.Array) {
	for _, field := range fields {
		a.Object(field)
	}
}

type Logger struct {
	zerolog *zerolog.Logger
	level   Level
}

func New(w io.Writer, level Level) *Logger {
	if level < Trace || level > Fatal {
		level = Info
	}

	zerologger := zerolog.New(w).
		Level(zerolog.Level(level)).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()

	return &Logger{
		zerolog: &zerologger,
		level:   level,
	}
}

// Log logs a message.
func (l *Logger) Log(msg string, fields FF) {
	l.zerolog.Log().Array(fiedsKey, fields).Msg(msg)
}

// Trace logs a message with log level Trace.
func (l *Logger) Trace(msg string, err error, fields FF) {
	l.zerolog.Trace().Array(fiedsKey, fields).Err(err).Msg(msg)
}

// Debug logs a message with log level Debug.
func (l *Logger) Debug(msg string, fields FF) {
	l.zerolog.Debug().Array(fiedsKey, fields).Msg(msg)
}

// Info logs a message with log level Info.
func (l *Logger) Info(msg string, fields FF) {
	l.zerolog.Info().Array(fiedsKey, fields).Msg(msg)
}

// Warn logs a message with log level Warn.
func (l *Logger) Warn(msg string, fields FF) {
	l.zerolog.Warn().Array(fiedsKey, fields).Msg(msg)
}

// Error logs a message with log level Error.
func (l *Logger) Error(msg string, err error, fields FF) {
	l.zerolog.Error().Array(fiedsKey, fields).Err(err).Msg(msg)
}

// Fatal logs a message with log level Fatal.
func (l *Logger) Fatal(msg string, err error, fields FF) {
	l.zerolog.Fatal().Array(fiedsKey, fields).Err(err).Msg(msg)
}

// Panic logs a message with log level Panic.
func (l *Logger) Panic(msg string, fields FF) {
	l.zerolog.Panic().Array(fiedsKey, fields).Msg(msg)
}

// GetLevel returns the logger log level.
func (l *Logger) GetLevel() Level {
	return l.level
}

func (l *Logger) Logger() *log.Logger {
	return nil
}
