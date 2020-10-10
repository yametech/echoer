package resource

import (
	"github.com/yametech/echoer/pkg/core"
	"github.com/yametech/echoer/pkg/storage"
	"github.com/yametech/echoer/pkg/storage/gtm"
)

type StepState = string

type StepSpec struct {
	FlowID     string `json:"flow_id" bson:"flow_id"`
	ActionRun  `json:"action_run" bson:"action_run"`
	Response   `json:"response" bson:"response"`
	RetryCount int32 `json:"retry_count" bson:"retry_count"`
}

type Response struct {
	State string `json:"state"`
}

type ActionRun struct {
	// reference define action
	ActionName string `json:"action_name" bson:"action_name"`
	// Params
	ActionParams map[string]interface{} `json:"action_params" bson:"action_params"`
	// parse from DSL .. eg: return Yes -> NextStep map[string]string{"Yes":"NextStep"}
	ReturnStateMap map[string]string `json:"return_state_map" bson:"return_state_map"`
	// action is done
	Done bool `json:"done" bson:"done"`
}

var _ core.IObject = &Step{}

var StepKind core.Kind = "step"

type Step struct {
	// default metadata for IObject
	core.Metadata `json:"metadata" bson:"metadata"`
	Spec          StepSpec `json:"spec" bson:"spec"`
}

func (s *Step) Clone() core.IObject {
	result := &Step{}
	core.Clone(s, result)
	return result
}

// Step impl Coder
func (*Step) Decode(op *gtm.Op) (core.IObject, error) {
	step := &Step{}
	if err := core.ObjectToResource(op.Data, step); err != nil {
		return nil, err
	}
	return step, nil
}

func init() {
	storage.AddResourceCoder(string(StepKind), &Step{})
}
