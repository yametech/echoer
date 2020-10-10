package resource

import (
	"github.com/yametech/echoer/pkg/core"
	"github.com/yametech/echoer/pkg/storage"
	"github.com/yametech/echoer/pkg/storage/gtm"
)

var _ core.IObject = &Flow{}

const FlowKind core.Kind = "flow"

type FlowStep struct {
	ActionName string            `json:"action_name" bson:"action_name"`
	Returns    map[string]string `json:"returns" bson:"returns"`
}

type FlowSpec struct {
	Steps []FlowStep `json:"steps" bson:"steps"`
}

type Flow struct {
	core.Metadata `json:"metadata" bson:"metadata"`
	Spec          FlowSpec `json:"spec" bson:"spec"`
}

func (f *Flow) Clone() core.IObject {
	result := &Step{}
	core.Clone(f, result)
	return result
}

var _ storage.Coder = &Flow{}

// Flow impl Coder
func (*Flow) Decode(op *gtm.Op) (core.IObject, error) {
	flow := &Flow{}
	if err := core.ObjectToResource(op.Data, flow); err != nil {
		return nil, err
	}
	return flow, nil
}

func init() {
	storage.AddResourceCoder("flow", &Flow{})
}
