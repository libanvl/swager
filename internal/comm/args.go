package comm

import (
	"errors"
	"strings"

	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/pkg/stoker"
)

type SwagerArgs interface {
}

func ToInitBlockArgs(tl stoker.TokenList) (SwagerArgs, error) {
	if len(tl) < 2 {
		return nil, errors.New("--init requires a tagname and blocktype")
	}

	args := &InitBlockArgs{Tag: tl[0], Block: tl[1]}
	if len(tl) > 2 {
		args.Args = tl[2:]
	}

	return args, nil
}

func ToSendToTagArgs(tl stoker.TokenList) (SwagerArgs, error) {
	if len(tl) < 2 {
		return nil, errors.New("--send requires a tagname and arguments")
	}

	args := &SendToTagArgs{Tag: tl[0], Args: tl[1:]}
	return args, nil
}

func ToServerArgs(tl stoker.TokenList) (SwagerArgs, error) {
	if len(tl) < 1 {
		return nil, errors.New("--server requires a subcommand")
	}

	switch tl[0] {
	case "exit":
		// swagerctl server exit
		return &ControlArgs{Command: ExitServer}, nil
	case "ping":
		// swagerctl server ping
		return &ControlArgs{Command: PingServer}, nil
	case "listen":
		return &ControlArgs{Command: RunServer}, nil
	case "reset":
		return &ControlArgs{Command: ResetServer}, nil
	}

	return nil, errors.New("unknown method")
}

func ToSetTagLogArgs(tl stoker.TokenList) (SwagerArgs, error) {
	if len(tl) < 2 {
		return nil, errors.New("--log requires a tag and a log level")
	}

	switch strings.ToLower(tl[1]) {
	case "default":
		return &SetTagLogArgs{Tag: tl[0], Level: core.DefaultLog}, nil
	case "info":
		return &SetTagLogArgs{Tag: tl[0], Level: core.InfoLog}, nil
	case "debug":
		return &SetTagLogArgs{Tag: tl[0], Level: core.DebugLog}, nil
	}

	return nil, errors.New("could not parse log level")
}
