package resource

import (
	"fmt"

	"github.com/yametech/echoer/pkg/core"
	"github.com/yametech/echoer/pkg/storage"
	"github.com/yametech/echoer/pkg/storage/gtm"
)

var _ core.IObject = &Flow{}

const FlowRunKind core.Kind = "flowrun"

type FlowRunSpec struct {
	Steps         []Step   `json:"steps" bson:"steps"`
	HistoryStates []string `json:"history_states" bson:"history_states"`
	LastState     string   `json:"last_state" bson:"last_state"`
	CurrentState  string   `json:"current_state" bson:"current_state"`
	LastEvent     string   `json:"last_event" bson:"last_event"`
	LastErr       string   `json:"last_err" bson:"last_err"`
}

func (f FlowRunSpec) GetStepByName(name string) (*Step, error) {
	var obj core.IObject
	for index, item := range f.Steps {
		if item.GetName() != name {
			continue
		}
		obj = (&f.Steps[index]).Clone()
	}
	if obj == nil {
		return nil, fmt.Errorf("get step (%s) not found", name)
	}
	return obj.(*Step), nil
}

type FlowRun struct {
	core.Metadata `json:"metadata" bson:"metadata"`
	Spec          FlowRunSpec `json:"spec" bson:"spec"`
}

func (f *FlowRun) Clone() core.IObject {
	result := &FlowRun{}
	core.Clone(f, result)
	return result
}

var _ storage.Coder = &FlowRun{}

// FlowRun impl Coder
func (*FlowRun) Decode(op *gtm.Op) (core.IObject, error) {
	flowRun := &FlowRun{}
	if err := core.ObjectToResource(op.Data, flowRun); err != nil {
		return nil, err
	}
	return flowRun, nil
}

func init() {
	storage.AddResourceCoder(string(FlowRunKind), &FlowRun{})
}
