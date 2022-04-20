package core

import (
	"errors"
	"fmt"
	"strings"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=LogLevel
type LogLevel int8

const (
	DefaultLog LogLevel = 0
	InfoLog    LogLevel = 10
	DebugLog   LogLevel = 30
)

func (l LogLevel) Debug() bool {
	return l >= DebugLog
}

func (l LogLevel) Info() bool {
	return l >= InfoLog
}

func (l LogLevel) Default() bool {
	return l >= DefaultLog
}

func (l *LogLevel) Set(s string) error {
	switch strings.ToLower(s) {
	case "debug":
		*l = DebugLog
		break
	case "info":
		*l = InfoLog
		break
	case "default":
		*l = DefaultLog
		break
	default:
		return errors.New("valid log levels are: default, info, debug")
	}

	return nil
}

type LogMessage interface {
	String() string
	Level() LogLevel
}

type PrefixLogMessage struct {
	prefix  string
	message string
	level   LogLevel
}

func (lm PrefixLogMessage) String() string {
	return fmt.Sprintf("[%s] %s", lm.prefix, lm.message)
}

func (lm PrefixLogMessage) Level() LogLevel {
	return lm.level
}

type LogChannel chan<- LogMessage

func (lc LogChannel) Send(level LogLevel, prefix string, msg string) {
	go func() {
		lc <- PrefixLogMessage{prefix, msg, level}
	}()
}

func (lc LogChannel) Sendf(level LogLevel, prefix string, format string, args ...interface{}) {
	lc.Send(level, prefix, fmt.Sprintf(format, args...))
}

type prefixLogger struct {
	logch  LogChannel
	prefix string
}

type Logger interface {
	Default(msg string)
	Defaultf(format string, args ...any)
	Info(msg string)
	Infof(format string, args ...any)
	Debug(msg string)
	Debugf(format string, args ...any)
}

func NewPrefixLogger(prefix string, logch LogChannel) Logger {
	return &prefixLogger{logch: logch, prefix: prefix}
}

func (l prefixLogger) Default(msg string) {
	l.logch.Send(DefaultLog, l.prefix, msg)
}

func (l prefixLogger) Defaultf(format string, args ...any) {
	l.logch.Sendf(DefaultLog, l.prefix, format, args...)
}

func (l prefixLogger) Info(msg string) {
	l.logch.Send(InfoLog, l.prefix, msg)
}

func (l prefixLogger) Infof(format string, args ...any) {
	l.logch.Sendf(InfoLog, l.prefix, format, args...)
}

func (l prefixLogger) Debug(msg string) {
	l.logch.Send(DebugLog, l.prefix, msg)
}

func (l prefixLogger) Debugf(format string, args ...any) {
	l.logch.Sendf(DebugLog, l.prefix, format, args...)
}
