package command

import (
	"fmt"

	"github.com/yametech/echoer/pkg/factory"
	"github.com/yametech/echoer/pkg/fsm/fss"
	"github.com/yametech/echoer/pkg/storage"
)

type ActionCmd struct {
	data []byte
	storage.IStorage
}

func (a *ActionCmd) Name() string {
	return `ACTION`
}

func (a *ActionCmd) Execute(args ...string) Reply {
	if reply := checkArgsExpected(args, 0); reply != nil {
		return reply
	}
	fsl := string(a.data)
	stmt, err := fss.NewActionFSLParser().Parse(fsl)
	if err != nil {
		return &ErrorReply{
			Message: fmt.Sprintf("parse flow (%s) error: (%s)", fsl, err),
		}
	}
	if err := factory.NewTranslation(factory.NewStoreImpl(a.IStorage)).ToAction(stmt); err != nil {
		return &ErrorReply{
			Message: fmt.Sprintf("factory translation flow (%s) error: (%s)", fsl, err),
		}
	}
	return &OkReply{}
}

func (a *ActionCmd) Help() string {
	return `
	ACTION name 
		ADDR = url ;
		METHOD = HTTP|GRPC ;
		ARGS = ARGS_EXPRESSION;
		RETURN = RETURN_STATE_EXPRESSION;
	ACTION_END
	`
}
