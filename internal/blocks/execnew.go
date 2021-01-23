package blocks

import (
  "errors"
  "fmt"
  "strconv"
  "strings"

  "github.com/libanvl/swager/internal/core"
)

type ExecNew struct {
  client core.Client
  opts *core.Options
  min int
  max int
}


func (e *ExecNew) Init(client core.Client, sub core.Sub, opts *core.Options) error {
  e.client = client
  e.opts = opts
  return nil
}

func (e *ExecNew) Configure(args []string) error {
  if len(args) != 2 {
    return errors.New("2 arguments are required.")
  }

  min, err := strconv.Atoi(args[0])
  max, err2 := strconv.Atoi(args[1])
  if err != nil || err2 != nil {
    return errors.New("Arguments must be ints")
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

  res, err := e.client.Command(fmt.Sprintf("workspace number %v, exec %s", curr + 1, strings.Join(args, " ")))
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

