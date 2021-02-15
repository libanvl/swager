package comm

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/adrg/xdg"
	"github.com/libanvl/swager/internal/core"
)

func init() {
  gob.Register(InitBlockArgs{})
  gob.Register(SendToTagArgs{})
  gob.Register(SetTagLogArgs{})
  gob.Register(ControlArgs{})
}

func GetSwagerSocket() (string, error) {
	swaysock, present := os.LookupEnv("SWAYSOCK")
	if !present {
		return "", errors.New("SWAYSOCK not set")
	}

	info, err := os.Stat(swaysock)
	if err != nil {
		return "", err
	}

	parts := strings.Split(info.Name(), ".")
	swayuid := parts[1]
	swaypid := parts[2]

	return xdg.RuntimeFile(fmt.Sprintf("swager-ipc/%s.%s.sock", swayuid, swaypid))
}

type SwagerMethod string

const (
	InitBlock SwagerMethod = "Swager.InitBlock"
	SendToTag SwagerMethod = "Swager.SendToTag"
	SetTagLog SwagerMethod = "Swager.SetTagLog"
	Control   SwagerMethod = "Swager.Control"
)

func (sm SwagerMethod) String() string {
	return string(sm)
}

type ServerControl int8

const (
	PingServer  ServerControl = 0
	RunServer   ServerControl = 1
	ResetServer ServerControl = 2
	ExitServer  ServerControl = math.MaxInt8
)

type InitBlockArgs struct {
	Tag   string
	Block string
	Args  []string
}

type SendToTagArgs struct {
	Tag  string
	Args []string
}

type SetTagLogArgs struct {
  Tag   string
  Level core.LogLevel
}

type ControlArgs struct {
	Command ServerControl
	Args    []string
}

type Reply struct {
  Args    SwagerArgs
	Success bool
}
