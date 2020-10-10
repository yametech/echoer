package command

import (
	"fmt"

	"github.com/yametech/echoer/pkg/factory"
	"github.com/yametech/echoer/pkg/fsm/fss"
	"github.com/yametech/echoer/pkg/storage"
)

type FlowCmd struct {
	data []byte
	storage.IStorage
}

func (f *FlowCmd) Name() string {
	return `FLOW`
}

func (f *FlowCmd) Execute(args ...string) Reply {
	if reply := checkArgsExpected(args, 0); reply != nil {
		return reply
	}
	fsl := string(f.data)
	stmt, err := fss.NewFlowFSLParser().Parse(fsl)
	if err != nil {
		return &ErrorReply{
			Message: fmt.Sprintf("parse flow (%s) error: (%s)", fsl, err),
		}
	}
	if err := factory.NewTranslation(factory.NewStoreImpl(f.IStorage)).ToFlow(stmt); err != nil {
		return &ErrorReply{
			Message: fmt.Sprintf("factory translation flow (%s) error: (%s)", fsl, err),
		}
	}
	return &OkReply{}
}

func (f *FlowCmd) Help() string {
	return `
	FLOW flow_name|flow_identifier  
		STEP step_name|step_identifier => RETURN_EXPRESSION {
			ACTION = action_name  ;
		};
	FLOW_END
	`
}
