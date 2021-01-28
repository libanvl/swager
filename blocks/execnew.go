package blocks

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/libanvl/swager/internal/core"
)


/*
ExecNew switches to the first uncreated workspace
then sends a command to sway

Registration

The block is registered with the name 'execnew'

Configuration

ExecNew requires two parameters for init configuration, min and max.
The newly created workspace will have the first available number in the range (min, max]

  swayctrl init <tagname> execnew 2 99

Send

Send the command to execute on the newly created workspace. This is intended
to be bound in sway config to a keyboard shortcut:

  $mod+shift+t exec swayctrl send <tagname> "exec alacritty"
*/
type ExecNew struct {
	client core.Client
	opts   *core.Options
	min    int
	max    int
}

func (e *ExecNew) Init(client core.Client, sub core.Sub, opts *core.Options) error {
	e.client = client
	e.opts = opts
	return nil
}

func (e *ExecNew) Configure(args []string) error {
	if len(args) != 2 {
    return errors.New("2 arguments are required: <min> <max>")
	}

	min, err := strconv.Atoi(args[0])
	max, err2 := strconv.Atoi(args[1])
	if err != nil || err2 != nil {
		return errors.New("Arguments must be ints")
	}

  if e.opts.Debug {
    e.opts.Log <- fmt.Sprintf("min: %d max: %d", min, max)
  }

	e.min = min
	e.max = max

	return nil
}

func (e *ExecNew) Run() {
}

func (e *ExecNew) Close() {
}

func (e *ExecNew) Receive(args []string) error {
	ws, err := e.client.Workspaces()
	if err != nil {
		return err
	}

	curr := e.min
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
  if e.opts.Debug {
    e.opts.Log <- fmt.Sprintf("running command on workspace: %d, '%s'", next, cmd)
  }

	res, err := e.client.Command(fmt.Sprintf("workspace number %d, exec %s", curr+1, cmd))
	if err != nil {
		return err
	}

	for _, r := range res {
		if !r.Success {
			return errors.New(r.Error)
		}
	}

	return nil
}
