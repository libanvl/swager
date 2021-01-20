package main

import (
  "bytes"
  "fmt"
  "log"
  "os/exec"
  "sync"

  "github.com/libanvl/swager/internal/blocks"
  "github.com/libanvl/swager/internal/core"
  "go.i3wm.org/i3/v4"
)

func main() {
  i3.SocketPathHook = func() (string, error) {
    out, err := exec.Command("sway", "--get-socketpath").CombinedOutput()
    if err != nil {
      return "", fmt.Errorf("getting sway socketpath: %v (output: %s)", err, out)
    }

    return string(out), nil
  }

  i3.IsRunningHook = func() bool {
    out, err := exec.Command("pgrep", "-c", "sway\\$").CombinedOutput()
    if err != nil {
      log.Printf("sway running: %v (output %s)", err, out)
    }

    return bytes.Compare(out, []byte("1")) == 0
  }

  blocks.RegisterBlocks()
  var wg sync.WaitGroup

  for key, factory := range core.Blocks {
    log.Printf("Found Block: %v", key)
    block := factory()
    if err := block.Init(core.WinMgrSway); err != nil {
      log.Panicf("Failed to Init block: %s", key)
    }
    if err := block.Configure(nil); err != nil {
      log.Panicf("Failed to Configure block: %s", key)
    }

    eventBlock, ok := block.(core.ChangeEventBlock)
    if ok {
      for _, et := range eventBlock.Event() {
        if et == i3.WindowEventType {
          wg.Add(1)
          go ProcessWindowEvents(&wg, eventBlock)
        }
      }
    }
  }

  log.Print("waiting...")
  wg.Wait()
}

func ProcessWindowEvents(wg *sync.WaitGroup, b core.ChangeEventBlock) {
  r := i3.Subscribe(i3.WindowEventType)
  for r.Next() {
    evt, ok := r.Event().(*i3.WindowEvent)
    if ok {
      log.Print("Got window event")
      if !b.MatchChange(evt.Change) {
        log.Print("Does not match change")
        continue
      }
      go b.OnEvent(evt)
    }
  }

  r.Close()
  b.Close()
  wg.Done()
}
