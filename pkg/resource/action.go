package resource

import (
	"encoding/json"
	"fmt"

	"github.com/yametech/echoer/pkg/core"
	"github.com/yametech/echoer/pkg/storage"
	"github.com/yametech/echoer/pkg/storage/gtm"
)

const ActionKind core.Kind = "action"

var _ core.IObject = &Action{}

type ServeType uint8

const (
	HTTP ServeType = iota
	GRPC
	HTTPS
)

type ParamType uint8

const (
	STR ParamType = iota
	INT
)

type ParamNameType string

type ActionParams = map[ParamNameType]ParamType

func CheckActionParams(runParams map[string]interface{}, actParams ActionParams) error {
	if len(runParams) < len(actParams) {
		actual, _ := json.Marshal(runParams)
		expected, _ := json.Marshal(actParams)
		return fmt.Errorf("parameters not enough actual (%s) and expected (%s)", actual, expected)
	}

	for k, v := range actParams {
		t, exist := runParams[string(k)]
		if !exist {
			return fmt.Errorf("run params (%s) not define", k)
		}

		err := fmt.Errorf("run params (%s) is Illegal type", k)
		switch t.(type) {
		case string:
			if v != STR {
				return err
			}
		case int64, int32, int16, int8, int, float32, float64, uint64, uint32, uint16, uint8, uint:
			if v != INT {
				return err
			}
		default:
			return err
		}

	}
	return nil
}

type ActionSpec struct {
	System string `json:"system" bson:"system"`
	// ServeType if client flow-controller is http, just support POST method
	ServeType `json:"serve_type" bson:"serve_type"`
	// Endpoints load balance client
	Endpoints []string `json:"endpoints" bson:"endpoints"`
	// if type is https CAPEM
	CaPEM string `json:"ca_pem" bson:"ca_pem"`
	// Params user define server params
	Params ActionParams `json:"params" bson:"params"`
	// ReturnStates the action will return the following expected state to the process flow-controller
	ReturnStates []string `json:"return_states" bson:"return_states"`
	// ReferenceCount if the reference count greater then 0 then it can't not delete
	ReferenceCount uint64 `json:"reference_count" bson:"reference_count"`
}

type Action struct {
	// Metadata default IObject Metadata
	core.Metadata `json:"metadata"`
	// Spec action spec
	Spec ActionSpec `json:"spec"`
}

func (a *Action) Clone() core.IObject {
	result := &Action{}
	core.Clone(a, result)
	return result
}

var _ storage.Coder = &Action{}

// Action impl Coder
func (*Action) Decode(op *gtm.Op) (core.IObject, error) {
	action := &Action{}
	if err := core.ObjectToResource(op.Data, action); err != nil {
		return nil, err
	}
	return action, nil
}

func init() {
	storage.AddResourceCoder(string(ActionKind), &Action{})
}
