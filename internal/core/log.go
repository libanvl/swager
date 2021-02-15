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

func (lc LogChannel) Print(prefix string, msg string) {
	go func() {
		lc <- PrefixLogMessage{prefix, msg, DefaultLog}
	}()
}

func (lc LogChannel) Printf(prefix string, format string, args ...interface{}) {
	lc.Print(prefix, fmt.Sprintf(format, args...))
}
