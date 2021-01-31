package core

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/libanvl/swager/pkg/ipc"
)

func init() {
	// core.Client must be a subset of ipc.Client
	// core.Subscription must be a subset of ipc.Subscription
	var _ Client = (*ipc.Client)(nil)
	var _ Sub = (*ipc.Subscription)(nil)
	var _ BlockLogMessage = defaultLogMessage{}
	var _ flag.Value = (*BlockLogLevel)(nil)
}

// Client exports a limited set of methods for use by core.Block instances.
type Client interface {
	Command(cmd string) ([]ipc.Command, error)
	CommandRaw(cmd string) (string, error)
	Workspaces() ([]ipc.Workspace, error)
	WorkspacesRaw() (string, error)
	Tree() (*ipc.Node, error)
	TreeRaw() (string, error)
	Version() (*ipc.Version, error)
	VersionRaw() (string, error)
}

// Sub exports a limited set of methods for use by core.Block instances.
type Sub interface {
	WorkspaceChanges() <-chan *ipc.WorkspaceChange
	WindowChanges() <-chan *ipc.WindowChange
	BindingChanges() <-chan *ipc.BindingChange
	ShutdownChanges() <-chan *ipc.ShutdownChange
	Ticks() <-chan *ipc.Tick
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=BlockLogLevel
type BlockLogLevel int8

const (
	DefaultLog BlockLogLevel = 0
	InfoLog    BlockLogLevel = 10
	DebugLog   BlockLogLevel = 30
)

func (l *BlockLogLevel) Set(s string) error {
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

type BlockLogMessage interface {
	String() string
	Level() BlockLogLevel
}

type defaultLogMessage struct {
	prefix  string
	message string
	level   BlockLogLevel
}

func (dlm defaultLogMessage) String() string {
	return fmt.Sprintf("[%s] %s", dlm.prefix, dlm.message)
}

func (dlm defaultLogMessage) Level() BlockLogLevel {
	return dlm.level
}

type BlockLogChannel chan<- BlockLogMessage

func (blc BlockLogChannel) Default(prefix string, msg string) {
	blc <- defaultLogMessage{prefix, msg, DefaultLog}
}

func (blc BlockLogChannel) Defaultf(prefix string, format string, args ...interface{}) {
	blc.Default(prefix, fmt.Sprintf(format, args...))
}

func (blc BlockLogChannel) Info(prefix string, msg string) {
	blc <- defaultLogMessage{prefix, msg, InfoLog}
}

func (blc BlockLogChannel) Infof(prefix string, format string, args ...interface{}) {
	blc.Info(prefix, fmt.Sprintf(format, args...))
}

func (blc BlockLogChannel) Debug(prefix string, msg string) {
	blc <- defaultLogMessage{prefix, msg, DebugLog}
}

func (blc BlockLogChannel) Debugf(prefix string, format string, args ...interface{}) {
	blc.Debug(prefix, fmt.Sprintf(format, args...))
}

type ServerControlRequest int8

const (
	ReloadRequest ServerControlRequest = 0
	ExitRequest   ServerControlRequest = 1
)

type ServerControlChannel chan<- ServerControlRequest

func (scc ServerControlChannel) RequestReload() {
	scc <- ReloadRequest
}

func (scc ServerControlChannel) RequestExit() {
	scc <- ExitRequest
}

// Options are shared options for use by core.Block instances.
// Debug indicates that debug logging was requested when starting the daemon.
// Use the Log channel to send log data back to the daemon.
type Options struct {
	Log    BlockLogChannel
	Server ServerControlChannel
}
