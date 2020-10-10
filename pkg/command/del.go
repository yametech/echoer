package command

import (
	"fmt"

	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/storage"
)

type Del struct {
	storage.IStorage
}

func (d *Del) Name() string {
	return `DEL`
}

func (d *Del) Execute(args ...string) Reply {
	if reply := checkArgsExpected(args, 3); reply != nil {
		return reply
	}
	resType := args[0]
	if storage.GetResourceCoder(resType) == nil {
		return &ErrorReply{Message: fmt.Sprintf("this type (%s) is not supported", resType)}
	}
	result := make(map[string]interface{})
	if err := d.Get(common.DefaultNamespace, resType, args[1], &result); err != nil {
		return &ErrorReply{Message: fmt.Sprintf("resource (%s) (%s) not exist or get error (%s)", resType, args[1], err)}
	}
	if err := d.Delete(common.DefaultNamespace, resType, args[1], args[2]); err != nil {
		return &ErrorReply{Message: fmt.Sprintf("delete resource (%s) (%s) error (%s)", resType, args[1], err)}
	}
	return &OkReply{Message: []byte("Ok")}
}

func (d *Del) Help() string {
	return `DEL resource_type name uuid`
}
