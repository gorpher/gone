package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LogLevel defines log levels.
type LogLevel uint8

// defines our own log levels, just in case we'll change logger in future
const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Log logs at the specified level for the specified sender
func Log(level LogLevel, sender string, format string, v ...any) {
	var ev *zerolog.Event
	switch level {
	case LevelDebug:
		ev = log.Logger.Debug()
	case LevelInfo:
		ev = log.Logger.Info()
	case LevelWarn:
		ev = log.Logger.Warn()
	case LevelError:
		ev = log.Logger.Error()
	default:
		ev = log.Logger.Error()
	}
	ev.Timestamp().Str("sender", sender)
	ev.Msg(fmt.Sprintf(format, v...))
}

func Debugf(format string, v ...any) {
	log.Logger.Debug().Msgf(format, v...)
}

func Infof(format string, v ...any) {
	log.Logger.Info().Msgf(format, v...)
}

func Warnf(format string, v ...any) {
	log.Logger.Warn().Msgf(format, v...)
}

func Errorf(format string, v ...any) {
	log.Logger.Error().Msgf(format, v...)
}

// Debug logs at debug level for the specified sender
func Debug(sender, format string, v ...any) {
	Log(LevelDebug, sender, format, v...)
}

// Info logs at info level for the specified sender
func Info(sender, format string, v ...any) {
	Log(LevelInfo, sender, format, v...)
}

// Warn logs at warn level for the specified sender
func Warn(sender, format string, v ...any) {
	Log(LevelWarn, sender, format, v...)
}

// Error logs at error level for the specified sender
func Error(sender, format string, v ...any) {
	Log(LevelError, sender, format, v...)
}
