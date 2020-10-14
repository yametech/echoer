package controller

import (
	"github.com/yametech/echoer/pkg/resource"
	"testing"

	"github.com/yametech/echoer/pkg/fsm"
)

func TestFlowRunController(t *testing.T) {
	fakeFS := &fakeFlowStorage{}
	fr, err := CreateFlowRunController("fake_test", fakeFS)
	if err != nil {
		t.Fatal(err)
	}
	step := createStepObject("fake_test", "", "fake_step_1", "fake_action_1", nil)
	fr.addSteps([]resource.Step{step})

	if err := fr.stepGraph(step, true, true); err != nil {
		t.Fatal(err)
	}

	if err := fr.Start(); err != nil {
		t.Fatal(err)
	}

	if fr.FlowRun.Spec.LastState != fsm.READY {
		t.Fatal("non expect state")
	}

	if fr.FlowRun.Spec.CurrentState != "fake_step_1" {
		t.Fatal("non expect state")
	}

	if fr.FlowRun.Spec.LastEvent != fsm.OpStart {
		t.Fatal("non expect state")
	}

	if err := fr.Pause(); err != nil {
		t.Fatal(err)
	}

	if fr.FlowRun.Spec.CurrentState != fsm.SUSPEND {
		t.Fatal("non expect state")
	}

	if err := fr.Continue(); err != nil {
		t.Fatal(err)
	}

	if fr.FlowRun.Spec.LastState != fsm.SUSPEND {
		t.Fatal("non expect state")
	}

	if fr.FlowRun.Spec.CurrentState != "fake_step_1" {
		t.Fatal("non expect state")
	}

	if fr.FlowRun.Spec.LastEvent != fsm.OpContinue {
		t.Fatal("non expect state")
	}

	if err := fr.Send(fsm.OpEnd); err != nil {
		t.Fatal(err)
	}

	if fr.FlowRun.Spec.LastState != "fake_step_1" {
		t.Fatal("non expect state")
	}

	if fr.FlowRun.Spec.CurrentState != fsm.DONE {
		t.Fatal("non expect state")
	}

	if fr.FlowRun.Spec.LastEvent != fsm.OpEnd {
		t.Fatal("non expect state")
	}

	_ = fr

}

/*
* The following unit test process
*
*  (☢️)READY --start--> STEP1 ---Yes--> STEP2 --end--> DONE(🚀)
*				|  ^
*               |  |------------------- |
*              pause --> SUSPEND --continue
 */
func TestFlowRunController2(t *testing.T) {
	fakeFS := &fakeFlowStorage{}
	fr, err := CreateFlowRunController("fake_test", fakeFS)
	if err != nil {
		t.Fatal(err)
	}

	step1 := createStepObject("fake_test", "", "fake_step_1", "fake_action_1", map[string]string{"Yes": "fake_step_2"})
	step2 := createStepObject("fake_test", "", "fake_step_2", "fake_action_2", nil)
	steps := []resource.Step{step1, step2}
	fr.addSteps(steps)

	if err := fr.stepGraph(step1, true, false); err != nil {
		t.Fatal(err)
	}
	if err := fr.stepGraph(step2, false, true); err != nil {
		t.Fatal(err)
	}

	// -------new state
	if err := fr.Start(); err != nil {
		t.Fatal(err)
	}

	if fr.FlowRun.Spec.LastState != fsm.READY {
		t.Fatal("non expect state")
	}

	if fr.FlowRun.Spec.CurrentState != "fake_step_1" {
		t.Fatal("non expect state")
	}

	if fr.FlowRun.Spec.LastEvent != fsm.OpStart {
		t.Fatal("non expect state")
	}

	// -------new state
	if err := fr.Pause(); err != nil {
		t.Fatal(err)
	}
	if fr.FlowRun.Spec.CurrentState != fsm.SUSPEND {
		t.Fatal("non expect state")
	}

	// -------new state
	if err := fr.Continue(); err != nil {
		t.Fatal(err)
	}

	if fr.FlowRun.Spec.LastState != fsm.SUSPEND {
		t.Fatal("non expect state")
	}

	if fr.FlowRun.Spec.CurrentState != "fake_step_1" {
		t.Fatal("non expect state")
	}

	if fr.FlowRun.Spec.LastEvent != fsm.OpContinue {
		t.Fatal("non expect state")
	}

	x := fr.FSM.AvailableTransitions()
	_ = x

	// -------new state
	event := StepStateEvent("fake_test", "fake_step_1", "Yes")
	if err := fr.Send(event); err != nil {
		t.Fatal(err)
	}

	if fr.FlowRun.Spec.LastState != "fake_step_1" {
		t.Fatal("non expect state")
	}

	if fr.FlowRun.Spec.CurrentState != "fake_step_2" {
		t.Fatal("non expect state")
	}

	if fr.FlowRun.Spec.LastEvent != event {
		t.Fatal("non expect state")
	}

	// -------new state
	if err := fr.Send(fsm.OpEnd); err != nil {
		t.Fatal(err)
	}

	if fr.FlowRun.Spec.LastState != "fake_step_2" {
		t.Fatal("non expect state")
	}

	if fr.FlowRun.Spec.CurrentState != fsm.DONE {
		t.Fatal("non expect state")
	}

	if fr.FlowRun.Spec.LastEvent != fsm.OpEnd {
		t.Fatal("non expect state")
	}

	_ = fr

}

/*
* The following unit test process
*
*  (☢️)READY --start--> STEP1 ---(NEXT trigger)Yes--> STEP2 --end--> DONE(🚀)
*				|  ^
*               |  |------------------- |
*              pause --> SUSPEND --continue
 */
func TestFlowRunController2Next(t *testing.T) {
	fakeFS := &fakeFlowStorage{}
	fr, err := CreateFlowRunController("fake_test", fakeFS)
	if err != nil {
		t.Fatal(err)
	}

	step1 := createStepObject("fake_test", "", "fake_step_1", "fake_action_1", map[string]string{"Yes": "fake_step_2"})
	step2 := createStepObject("fake_test", "", "fake_step_2", "fake_action_2", nil)
	steps := []resource.Step{step1, step2}
	fr.addSteps(steps)

	if err := fr.stepGraph(step1, true, false); err != nil {
		t.Fatal(err)
	}
	if err := fr.stepGraph(step2, false, true); err != nil {
		t.Fatal(err)
	}

	// -------new state
	if err := fr.Start(); err != nil {
		t.Fatal(err)
	}

	if fr.FlowRun.Spec.LastState != fsm.READY ||
		fr.FlowRun.Spec.CurrentState != "fake_step_1" ||
		fr.FlowRun.Spec.LastEvent != fsm.OpStart {
		t.Fatal("non expect state")
	}

	// -------new state
	if err := fr.Pause(); err != nil {
		t.Fatal(err)
	}
	if fr.FlowRun.Spec.CurrentState != fsm.SUSPEND {
		t.Fatal("non expect state")
	}

	// -------new state
	if err := fr.Continue(); err != nil {
		t.Fatal(err)
	}

	if fr.FlowRun.Spec.LastState != fsm.SUSPEND ||
		fr.FlowRun.Spec.LastEvent != fsm.OpContinue ||
		fr.FlowRun.Spec.CurrentState != "fake_step_1" {
		t.Fatal("non expect state")
	}

	// -------next
	if err := fr.Next(); err != nil {
		t.Fatal(err)
	}

	// -------new state
	if err := fr.Send(fsm.OpEnd); err != nil {
		t.Fatal(err)
	}

	if fr.FlowRun.Spec.LastState != "fake_step_2" {
		t.Fatal("non expect state")
	}

	if fr.FlowRun.Spec.CurrentState != fsm.DONE {
		t.Fatal("non expect state")
	}

	if fr.FlowRun.Spec.LastEvent != fsm.OpEnd {
		t.Fatal("non expect state")
	}

	_ = fr

}
