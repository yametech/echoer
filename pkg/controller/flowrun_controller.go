package controller

import (
	"encoding/json"
	"fmt"

	"github.com/yametech/echoer/pkg/core"
	"github.com/yametech/echoer/pkg/fsm"
	"github.com/yametech/echoer/pkg/resource"
)

func createStepObject(flowID, flowRunUUID, stepName, actionName string, returnStateMap map[string]string) resource.Step {
	// create new step action runtime data
	step := resource.Step{
		Metadata: core.Metadata{
			Name: stepName,
			Kind: resource.StepKind,
		},
		Spec: resource.StepSpec{
			FlowID:      flowID,
			FlowRunUUID: flowRunUUID,
			ActionRun: resource.ActionRun{
				ActionName:     actionName,
				Done:           false,
				ReturnStateMap: returnStateMap,
			},
			Response: resource.Response{},
		},
	}
	return step
}

func StateList(states ...string) []string { return states }

func StepStateEvent(flowName, step, state string) string {
	return fmt.Sprintf("%s_%s_%s", flowName, step, state)
}

type FlowRunController struct {
	*resource.FlowRun
	*fsm.FSM
	fsi StorageInterface
}

func CreateFlowRunController(flowRuntimeName string, fsi StorageInterface) (*FlowRunController, error) {
	frt, err := fsi.Query(flowRuntimeName)
	if err == nil && frt != nil {
		return frt, nil
	}

	// if other then create new flow flow-controller
	flowFSM := fsm.NewFSM(fsm.READY, nil, nil)
	frtSpec := resource.FlowRunSpec{
		Steps:         make([]resource.Step, 0),
		HistoryStates: make([]string, 0),
	}
	frt = &FlowRunController{
		FlowRun: &resource.FlowRun{
			Metadata: core.Metadata{
				Name: flowRuntimeName,
				Kind: resource.FlowRunKind,
			},
			Spec: frtSpec,
		},
		FSM: flowFSM,
		fsi: fsi,
	}
	frt.GenerateVersion()
	if updateErr := fsi.Update(frt.FlowRun); updateErr != nil {
		return nil, err
	}
	return frt, nil
}

func (f *FlowRunController) Start() error {
	return f.event(fsm.OpStart)
}

func (f *FlowRunController) Stop() error {
	var expectStates []string
	for _, step := range f.Spec.Steps {
		expectStates = append(expectStates, step.GetName())
	}
	f.FSM.Add(fsm.OpStop, StateList(expectStates...), fsm.STOPPED, nil)
	if err := f.event(fsm.OpStop); err != nil {
		return err
	}
	return nil
}

func (f *FlowRunController) Send(e string) error {
	return f.event(e)
}

func (f *FlowRunController) event(e string) error {
	f.FlowRun.Spec.LastEvent = e
	if err := f.FSM.Event(e); err != nil {
		return err
	}
	f.stateMemento()
	if err := f.fsi.Update(f.FlowRun); err != nil {
		return err
	}
	return nil
}

func (f *FlowRunController) Pause() error {
	var expectStates []string
	for _, step := range f.Spec.Steps {
		expectStates = append(expectStates, step.GetName())
	}
	f.FSM.Add(fsm.OpPause, StateList(expectStates...), fsm.SUSPEND, nil)

	if err := f.event(fsm.OpPause); err != nil {
		return err
	}

	return nil
}

func (f *FlowRunController) Continue() error {
	f.FSM.Add(fsm.OpContinue, StateList(fsm.SUSPEND), f.Last(), nil)
	f.stateMemento()
	if err := f.event(fsm.OpContinue); err != nil {
		return err
	}
	return nil
}

// Next Only one result of state switching can use next
// eg: A --Yes--> B
func (f *FlowRunController) Next() error {
	var next string
	for _, event := range f.AvailableTransitions() {
		needIgnore := false
		switch event {
		case fsm.OpPause, fsm.OpStop:
			needIgnore = true
		}
		if needIgnore {
			continue
		}
		next = event
		break
	}
	if next == "" {
		bs, marshalErr := json.Marshal(f)
		if marshalErr != nil {
			return fmt.Errorf("the next state is illegal (unknow error)")
		} else {
			return fmt.Errorf("the next state is illegal (%s)", bs)
		}
	}
	if err := f.event(next); err != nil {
		return err
	}

	return nil
}

func (f *FlowRunController) stateMemento() {
	_current, _last := f.Current(), f.Last()
	f.Spec.CurrentState = _current
	f.Spec.LastState = _last
	f.Spec.HistoryStates = append(f.Spec.HistoryStates, _current)
}

func contains(all []resource.Step, item resource.Step) bool {
	for _, n := range all {
		if item.GetName() == n.GetName() {
			return true
		}
	}
	return false
}

func (f *FlowRunController) addSteps(steps []resource.Step) {
	// check action exist and create flow flow-controller runtime action
	for _, step := range steps {
		if !contains(f.Spec.Steps, step) {
			f.Spec.Steps = append(f.Spec.Steps, step)
		}
	}
	return
}

func (f *FlowRunController) getStep(name string) *resource.Step {
	for _, step := range f.FlowRun.Spec.Steps {
		if step.Name == name {
			return (&step).Clone().(*resource.Step)
		}
	}
	return nil
}

func (f *FlowRunController) stepGraph(step resource.Step, first, last bool) error {
	currentStepState := step.GetName()
	if first {
		// first step init dependent fsm READY state
		f.Add(fsm.OpStart, StateList(fsm.READY), step.GetName(),
			func(event *fsm.Event) {
				if err := f.fsi.Create(&step); err != nil {
					f.Spec.LastErr = err.Error()
				}
				fmt.Printf("[INFO] flow run (%s) add first step action run (%s)\n", f.Name, step.GetName())
			})
	}

	for entryState, expectState := range step.Spec.ActionRun.ReturnStateMap {
		listenEvent := StepStateEvent(f.Name, step.GetName(), entryState)
		expectStep := expectState
		callback := func(event *fsm.Event) {
			if expectStep == fsm.DONE {
				return
			}
			// create target step based on call
			expectStateStep := f.getStep(expectStep)
			if err := f.fsi.Create(expectStateStep); err != nil {
				f.Spec.LastErr = err.Error()
			}

			fmt.Printf(
				"[INFO] flow run (%s) listen event (%s) create by step (%s)\n",
				f.Name, listenEvent, currentStepState)
		}
		f.Add(listenEvent, StateList(step.GetName()), expectState, callback)
		fmt.Printf(
			"[DEBUG] flow run (%s) listen event (%s) on state (%s) to (%s)\n",
			f.Name, listenEvent, step.GetName(), expectState)
	}

	if last {
		// last step init dependent fsm DONE state
		f.Add(fsm.OpEnd, StateList(step.GetName()), fsm.DONE, nil)
		if err := f.fsi.Update(f.FlowRun); err != nil {
			return err
		}
	}

	return nil
}
