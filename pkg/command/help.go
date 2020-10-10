package command

import (
	"fmt"
	"strings"
)

const (
	NotEnoughArgs = `expected %d argument but not enough.`
)

func checkArgsExpected(args []string, expected int) Reply {
	if len(args) != expected {
		return &ErrorReply{Message: fmt.Errorf(NotEnoughArgs, expected)}
	}
	return nil
}

type Help struct {
	Message string
}

func (h *Help) Name() string {
	return `help`
}

func (h *Help) Execute(args ...string) Reply {
	if len(args) == 0 {
		return &RawReply{Message: []byte(h.Help())}
	}
	reply := &RawReply{}
	var cmd Command
	switch strings.ToLower(args[0]) {
	case "del":
		cmd = &Del{}
	case "flow":
		cmd = &FlowCmd{}
	case "flow_run":
		cmd = &FlowRunCmd{}
	case "get":
		cmd = &Get{}
	default:
		cmd = &Help{}
	}
	reply.Message = []byte(cmd.Help())
	return reply
}

func (h *Help) Help() string {
	return `
USAGE:
	HELP cmd
`
}
