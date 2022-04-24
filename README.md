# swager

Send commands using Sway IPC protocol in response to Sway events or user input.

TBD...

## Example Usage

Call this script from Sway config using exec_always:

```sh
#! /bin/sh

printf '========[ %s ]========\n' "$(date)" >> /tmp/swagerd.log
swagerd -log info >> /tmp/swagerd.log 2>&1 &
trap "swagerctl --server exit" HUP INT TERM EXIT

printf '========[ %s ]========\n' "$(date)" >> /tmp/swagerctl.log
{
  swagerctl --server reset
  swagerctl --init mon swaymon
  swagerctl --init al autolay -masterstack 1 2 3 4 -autotiler 5 6 7 8
  swagerctl --init newws execnew 3 7
  swagerctl --init spawn initspawn
  swagerctl --server listen
  swagerctl --send spawn 8 "exec termc"
} >> /tmp/swagerctl.log 2>&1

wait
```

