package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"time"

	"github.com/adrg/xdg"
	"github.com/libanvl/swager/internal/comm"
)

func main() {
  daemon_path := flag.String("daemon-path", "swagerd", "The path to the swager daemon executable")
  flag.Parse()

  cmd := exec.Command("pgrep", "--uid", fmt.Sprint(os.Getuid()), "--exact", *daemon_path)
  out, err := cmd.StdoutPipe()
  if err != nil {
    log.Fatal("cmd pipe:", err)
  }

  var pid string
  if err := cmd.Start(); err != nil {
    pid = startd(*daemon_path)
  } else {
    scanner := bufio.NewScanner(out)
    for scanner.Scan() {
      p := scanner.Text()
      addr, err := xdg.RuntimeFile(fmt.Sprintf("swager/d-%s.sock", p))
      if err != nil {
        log.Fatal("runtime file:", err)
      }

      conn, err := net.DialTimeout("unix", addr, time.Millisecond * 500)
      if err != nil {
        log.Fatal("dial:", err)
      }

      client := rpc.NewClient(conn)
      reply := new(comm.Reply)
      if err := client.Call(comm.Control.String(), &comm.ControlArgs{Command: comm.PingServer}, reply); err != nil {
        log.Print("ping:", err)
        continue
      } else {
        pid = p
        break
      }
    }

    if err := cmd.Wait(); err != nil {
      log.Fatal("cmd:", err)
    }
  }
}

func startd(path string) int {
  
}
