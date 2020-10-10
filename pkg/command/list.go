package command

import (
	"encoding/json"
	"fmt"

	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/storage"
)

type List struct {
	storage.IStorage
}

func (l *List) Name() string {
	return `LIST`
}

func (l *List) Execute(args ...string) Reply {
	if reply := checkArgsExpected(args, 1); reply != nil {
		return reply
	}
	resType := args[0]
	if storage.GetResourceCoder(resType) == nil {
		return &ErrorReply{Message: fmt.Sprintf("list this type (%s) is not supported", resType)}
	}
	result, err := l.List(common.DefaultNamespace, resType, "")
	if err != nil {
		return &ErrorReply{Message: fmt.Sprintf("list resource (%s) not exist or get error (%s)", resType, err)}
	}
	bs, err := json.Marshal(result)
	if err != nil {
		return &ErrorReply{Message: fmt.Sprintf("list resource (%s) unmarshal byte error (%s)", resType, err)}
	}
	return &RawReply{Message: bs}
}

func (l *List) Help() string {
	return `list resource_type`
}
