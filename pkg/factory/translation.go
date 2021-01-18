package factory

import (
	"fmt"
	"github.com/yametech/echoer/pkg/core"
	"github.com/yametech/echoer/pkg/fsm"
	"github.com/yametech/echoer/pkg/fsm/fss"
	"github.com/yametech/echoer/pkg/resource"
	"github.com/yametech/echoer/pkg/storage"
)

type Translation struct {
	IStore
}

func NewTranslation(store IStore) *Translation {
	return &Translation{store}
}

func (t *Translation) ToFlowRun(stmt *fss.FlowRunStmt) error {
	_, err := t.GetFlowRun(stmt.Flow)
	if err != nil {
		if err == storage.NotFound {
			goto CREATE
		}
		return err
	}
	return fmt.Errorf("flow run (%s) already exist", stmt.Flow)
CREATE:
	var fr = &resource.FlowRun{
		Metadata: core.Metadata{
			Name: stmt.Flow,
			Kind: resource.FlowRunKind,
		},
		Spec: resource.FlowRunSpec{
			Steps:         make([]resource.Step, 0),
			HistoryStates: []string{fsm.READY},
		},
	}
	fr.GenerateVersion()

	var returnSteps []string
	steps := make(map[string]interface{})

	// check step action if exist
	for _, step := range stmt.Steps {
		steps[step.Name] = ""
		action, err := t.GetAction(step.Action.Name)
		if err != nil {
			return fmt.Errorf("not without getting action (%s) definition", step.Action.Name)
		}

		actionParams := make(map[string]interface{})
		// check whether the parameter type and variable name are correct
		for _, arg := range step.Action.Args {
			actParamType, exist := action.Spec.Params[resource.ParamNameType(arg.Name)]
			if !exist {
				return fmt.Errorf("step (%s) args (%s) not defined", step.Name, arg.Name)
			}
			switch arg.ParamType {
			case fss.StringType:
				if actParamType != resource.STR {
					return fmt.Errorf("step (%s) args (%s) illegal type", step.Name, arg.Name)
				}
			case fss.NumberType:
				if actParamType != resource.INT {
					return fmt.Errorf("step (%s) args (%s) illegal type", step.Name, arg.Name)
				}
			default:
				return fmt.Errorf("step (%s) args (%s) illegal type", step.Name, arg.Name)
			}
			actionParams[arg.Name] = arg.Value
		}

		returnStateMap := make(map[string]string)

		// check whether the returns are correct
		for _, _return := range step.Returns {
			if _return.Next != "done" {
				returnSteps = append(returnSteps, _return.Next)
			}
			if !stringInSlice(_return.State, action.Spec.ReturnStates) {
				return fmt.Errorf("step (%s) return state (%s) illegal type", step.Name, _return.State)
			}
			returnStateMap[_return.State] = _return.Next
		}
		flowRunStep := resource.Step{
			Metadata: core.Metadata{
				Name: step.Name,
				Kind: resource.StepKind,
			},
			Spec: resource.StepSpec{
				FlowID:      stmt.Flow,
				FlowRunUUID: fr.GetUUID(),
				ActionRun: resource.ActionRun{
					ActionName:     step.Action.Name,
					ActionParams:   actionParams,
					ReturnStateMap: returnStateMap,
					Done:           false,
				},
			},
		}
		flowRunStep.GenerateVersion()
		fr.Spec.Steps = append(fr.Spec.Steps, flowRunStep)
	}

	for _, name := range returnSteps {
		if _, ok := steps[name]; !ok {
			return fmt.Errorf("not without getting step (%s) definition", name)
		}
	}
	err = t.CreateFlowRun(fr)
	if err != nil {
		return err
	}

	return nil
}

func (t *Translation) ToFlow(stmt *fss.FlowStmt) error {
	fl := &resource.Flow{
		Metadata: core.Metadata{
			Name: stmt.Flow,
			Kind: resource.FlowKind,
		},
	}
	// check step action if exist
	for _, step := range stmt.Steps {
		action, err := t.GetAction(step.Action.Name)
		if err != nil {
			return fmt.Errorf("not without getting action (%s) definition", step.Action.Name)
		}

		returnStateMap := make(map[string]string)
		// check whether the returns are correct
		for _, _return := range step.Returns {
			if !stringInSlice(_return.State, action.Spec.ReturnStates) {
				return fmt.Errorf("return state (%s) illegal type", _return.State)
			}
			returnStateMap[_return.State] = _return.Next
		}

		flowStep := resource.FlowStep{
			ActionName: action.GetName(),
			Returns:    returnStateMap,
		}
		fl.Spec.Steps = append(fl.Spec.Steps, flowStep)
	}
	err := t.CreateFlow(fl)
	if err != nil {
		return err
	}
	return nil
}

func (t *Translation) ToAction(stmt *fss.ActionStmt) error {
	_, err := t.GetAction(stmt.ActionStatement.Name)
	if err != nil {
		if err == storage.NotFound {
			goto CREATE
		}
		return err
	}
	return fmt.Errorf("action (%s) already exist", stmt.ActionStatement.Name)
CREATE:
	action := &resource.Action{
		Metadata: core.Metadata{
			Name: stmt.ActionStatement.Name,
			Kind: resource.ActionKind,
		},
		Spec: resource.ActionSpec{
			Params:       make(resource.ActionParams),
			Endpoints:    make([]string, 0),
			ReturnStates: make([]string, 0),
		},
	}
	for _, addr := range stmt.ActionStatement.Addr {
		action.Spec.Endpoints = append(action.Spec.Endpoints, addr)
	}

	for _, _return := range stmt.ActionStatement.Returns {
		action.Spec.ReturnStates = append(action.Spec.ReturnStates, _return.State)
	}

	for _, _arg := range stmt.ActionStatement.Args {
		switch _arg.ParamType {
		case fss.StringType:
			action.Spec.Params[resource.ParamNameType(_arg.Name)] = resource.STR
		case fss.NumberType:
			action.Spec.Params[resource.ParamNameType(_arg.Name)] = resource.INT
		default:
			return fmt.Errorf("invalid parameter definition (%s) type (%d)\n", _arg.Name, _arg.ParamType)
		}
	}
	if err := t.CreateAction(action); err != nil {
		return err
	}
	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
