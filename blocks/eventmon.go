package blocks

import (
	"errors"
	"fmt"
	"sync"

	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/ipc"
)

type EventMon struct {
	core.BasicBlock
	sub     core.Sub
	opts    *core.Options
	cookies map[ipc.EventPayloadType]ipc.Cookie
	logmx   sync.Mutex
	log     core.Logger
}

func init() {
	var _ core.BlockInitializer = (*EventMon)(nil)
	var _ core.Receiver = (*EventMon)(nil)
}

func (em *EventMon) Init(client core.Client, sub core.Sub, opts *core.Options, log core.Logger, args ...string) error {
	em.sub = sub
	em.opts = opts
	em.log = log

	em.cookies = make(map[ipc.EventPayloadType]ipc.Cookie)
	return nil
}

func (em *EventMon) SetLogLevel(level core.LogLevel) {

}

func (em *EventMon) Receive(args []string) error {
	if len(args) < 1 {
		return errors.New("EventMon requires one argument: <workspace|window|tick>")
	}

	evt := args[1]

	switch evt {
	case "workspace":
		_, err := em.sub.WorkspaceChanges(em.WorkspaceChanged)
		if err != nil {
			return err
		}
		break
	case "window":
		_, err := em.sub.WindowChanges(em.WindowChanged)
		if err != nil {
			return err
		}
		break
	case "tick":
		_, err := em.sub.Ticks(em.Ticked)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unsupported event: %v", evt)
	}

	return nil
}

func (em *EventMon) WorkspaceChanged(evt ipc.WorkspaceChange) {
	em.logmx.Lock()
	defer em.logmx.Unlock()
	em.log.Defaultf("%#v\n", evt)
}

func (em *EventMon) WindowChanged(evt ipc.WindowChange) {
	em.logmx.Lock()
	defer em.logmx.Unlock()
	em.log.Defaultf("%#v\n", evt)
}

func (em *EventMon) Ticked(evt ipc.Tick) {
	em.logmx.Lock()
	defer em.logmx.Unlock()
	em.log.Defaultf("%#v\n", evt)
}
