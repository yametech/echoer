package command

import (
	"encoding/json"
	"fmt"

	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/storage"
)

type Get struct {
	storage.IStorage
}

func (g *Get) Name() string {
	return `GET`
}

func (g *Get) Execute(args ...string) Reply {
	if reply := checkArgsExpected(args, 2); reply != nil {
		return reply
	}
	resType := args[0]
	if storage.GetResourceCoder(resType) == nil {
		return &ErrorReply{Message: fmt.Sprintf("this type (%s) is not supported", resType)}
	}
	result := make(map[string]interface{})
	if err := g.Get(common.DefaultNamespace, resType, args[1], &result); err != nil {
		return &ErrorReply{Message: fmt.Sprintf("resource (%s) (%s) not exist or get error (%s)", resType, args[1], err)}
	}
	bs, err := json.Marshal(result)
	if err != nil {
		return &ErrorReply{Message: fmt.Sprintf("resource (%s) (%s) unmarshal byte error (%s)", resType, args[1], err)}
	}
	return &RawReply{Message: bs}
}

func (g *Get) Help() string {
	return `GET resource_type name`
}
