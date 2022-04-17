package blocks

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/libanvl/swager/internal/core"
)

type ExecNew struct {
	client   core.Client
	opts     *core.Options
	min      int
	max      int
	log      core.Logger
	loglevel core.LogLevel
}

func init() {
	var _ core.BlockInitializer = (*ExecNew)(nil)
	var _ core.Receiver = (*ExecNew)(nil)
}

func (e *ExecNew) Init(client core.Client, sub core.Sub, opts *core.Options, log core.Logger, args ...string) error {
	e.client = client
	e.opts = opts
	e.log = log

	if len(args) != 2 {
		return errors.New("2 arguments are required: <min> <max>")
	}

	min, err := strconv.Atoi(args[0])
	max, err2 := strconv.Atoi(args[1])
	if err != nil || err2 != nil {
		return errors.New("Arguments must be ints")
	}

	e.log.Infof("min: %d max: %d", min, max)

	e.min = min
	e.max = max
	return nil
}

func (e *ExecNew) SetLogLevel(level core.LogLevel) {
	e.loglevel = level
}

func (e *ExecNew) Receive(args []string) error {
	ws, err := e.client.Workspaces()
	if err != nil {
		return err
	}

	if e.loglevel.Debug() {
		e.log.Debugf("got workspaces. count: %d", len(ws))
	}

	curr := e.min - 1
	for _, w := range ws {
		if w.Num > e.max {
			continue
		}

		if w.Num > curr {
			curr = w.Num
		}
	}

	next := curr + 1
	cmd := strings.Join(args, " ")

	e.log.Debugf("running command on workspace: %d, '%s'", next, cmd)

	res, err := e.client.Command(fmt.Sprintf("workspace number %d, exec %s", next, cmd))
	if err != nil {
		return err
	}

	for _, r := range res {
		e.log.Debugf("result: %#v", r)
		if !r.Success {
			return errors.New(r.Error)
		}
	}

	return nil
}
