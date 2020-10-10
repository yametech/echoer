package command

import (
	"fmt"

	"github.com/yametech/echoer/pkg/factory"
	"github.com/yametech/echoer/pkg/fsm/fss"
	"github.com/yametech/echoer/pkg/storage"
)

type FlowRunCmd struct {
	data []byte
	storage.IStorage
}

func (f *FlowRunCmd) Name() string {
	return `FLOW_RUN`
}

func (f *FlowRunCmd) Execute(args ...string) Reply {
	if reply := checkArgsExpected(args, 0); reply != nil {
		return reply
	}
	fsl := string(f.data)
	stmt, err := fss.NewFlowRunFSLParser().Parse(fsl)
	if err != nil {
		return &ErrorReply{
			Message: fmt.Sprintf("parse flow run (%s) error: (%s)", fsl, err),
		}
	}
	if err := factory.NewTranslation(factory.NewStoreImpl(f.IStorage)).ToFlowRun(stmt); err != nil {
		return &ErrorReply{
			Message: fmt.Sprintf("factory translation flow run (%s) error: (%s)", fsl, err),
		}
	}

	return &OkReply{}
}

func (f *FlowRunCmd) Help() string {
	return `
	FLOW_RUN flow_run_name|flow_run_identifier
		STEP step_name|identifier => RETURN_EXPRESSION {
			ACTION = action_name ;
			ARGS = ARGS_EXPRESSION ;
		};
	FLOW_RUN_END
	`
}
