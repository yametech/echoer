package command

import (
	"encoding/json"
	"fmt"
	"github.com/yametech/echoer/pkg/core"

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
	results, err := l.List(common.DefaultNamespace, resType, "")
	if err != nil {
		return &ErrorReply{Message: fmt.Sprintf("list resource (%s) not exist or get error (%s)", resType, err)}
	}

	format := NewFormat()
	format.Header("name", "type", "uuid", "version")
	for _, result := range results {
		var _res = result
		bs, err := json.Marshal(_res)
		if err != nil {
			return &ErrorReply{Message: fmt.Sprintf("list resource (%s) unmarshal byte error (%s)", resType, err)}
		}
		typePointer := storage.GetResourceCoder(resType)
		obj := typePointer.(core.IObject).Clone()
		if err := core.JSONRawToResource(bs, obj); err != nil {
			return &ErrorReply{Message: fmt.Sprintf("list resource (%s) unmarshal byte error (%s)", resType, err)}
		}
		format.Row(obj.GetName(), resType, obj.GetUUID(), fmt.Sprintf("%d", obj.GetResourceVersion()))
	}

	return &RawReply{Message: format.Out()}
}

func (l *List) Help() string {
	return `list resource_type`
}
