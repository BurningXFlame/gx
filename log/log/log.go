/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

// Logging Facade with leveled logging, tagged logging.
package log

import (
	"errors"
	"fmt"
	"sync/atomic"
)

// A concrete logger should implement interface Logger.
type Logger interface {
	// Print a log message
	Printf(format string, v ...any)
	// Close the logger. Flush buffer, close files, etc.
	Close() error
}

// Log Level
type Level uint8

const (
	LevelError Level = iota
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

type conf struct {
	logger Logger
	level  Level
}

var (
	conf0   = conf{}
	theConf atomic.Value
)

func init() {
	theConf.Store(conf0)
}

var (
	errNilLogger    = errors.New("nil logger")
	errInvalidLevel = errors.New("invalid level")
)

// Register logger as the global logger.
// Can be called multiple times. In this case, the old logger will be closed, before the new one takes effect.
func Set(logger Logger, level Level) error {
	if logger == nil {
		return errNilLogger
	}

	if level < LevelError || level > LevelTrace {
		return errInvalidLevel
	}

	_ = Close()

	theConf.Store(conf{
		logger: logger,
		level:  level,
	})

	return nil
}

func Error(format string, v ...any) {
	printf(LevelError, format, v...)
}

func Warn(format string, v ...any) {
	printf(LevelWarn, format, v...)
}

func Info(format string, v ...any) {
	printf(LevelInfo, format, v...)
}

func Debug(format string, v ...any) {
	printf(LevelDebug, format, v...)
}

func Trace(format string, v ...any) {
	printf(LevelTrace, format, v...)
}

func printf(level Level, format string, v ...any) {
	if level < LevelError || level > LevelTrace {
		return
	}

	the := theConf.Load().(conf)

	if the.level < level {
		return
	}

	if the.logger == nil {
		return
	}

	the.logger.Printf(levelPrefixes[level]+format, v...)
}

var levelPrefixes = [5]string{"ERROR ", "WARN  ", "INFO  ", "DEBUG ", "TRACE "}

type TagLogger interface {
	Error(format string, v ...any)
	Warn(format string, v ...any)
	Info(format string, v ...any)
	Debug(format string, v ...any)
	Trace(format string, v ...any)

	WithTag(tag string) TagLogger
}

// Create a TagLogger, which prints "[tag]" before every log message.
// Usually used for module-specific logging, request-specific logging, etc.
// WithTag may be chained together. e.g. WithTag("tag").WithTag("tag2") creates a TagLogger, which prints "[tag] [tag2]" before every log message.
func WithTag(tag string) TagLogger {
	if len(tag) == 0 {
		return &tagLogger{}
	}

	return &tagLogger{
		tag: fmt.Sprintf("[%v] ", tag),
	}
}

type tagLogger struct {
	tag string
}

func (l *tagLogger) Error(format string, v ...any) {
	Error(l.tag+format, v...)
}

func (l *tagLogger) Warn(format string, v ...any) {
	Warn(l.tag+format, v...)
}

func (l *tagLogger) Info(format string, v ...any) {
	Info(l.tag+format, v...)
}

func (l *tagLogger) Debug(format string, v ...any) {
	Debug(l.tag+format, v...)
}

func (l *tagLogger) Trace(format string, v ...any) {
	Trace(l.tag+format, v...)
}

func (l *tagLogger) WithTag(tag string) TagLogger {
	return &tagLogger{
		tag: l.tag + fmt.Sprintf("[%v] ", tag),
	}
}

// Close the logger. Flush buffer, close files, etc.
// Must be called before process exit.
func Close() error {
	the := theConf.Swap(conf0).(conf)

	if the.logger == nil {
		return nil
	}

	return the.logger.Close()
}
