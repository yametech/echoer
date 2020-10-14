package controller

import (
	"fmt"
	"time"

	"github.com/laik/timerqueue"
	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/fsm"
	"github.com/yametech/echoer/pkg/resource"
	"github.com/yametech/echoer/pkg/storage"
)

type Controller interface {
	Run() error
	Stop() error
}

var _ timerqueue.Timer = &DelayStepAction{}

type DelayStepAction struct {
	step *resource.Step
	storage.IStorage
}

func (dsa *DelayStepAction) OnTimer(t time.Time) {
	dsa.step.Spec.RetryCount += 1
	// check the flow run state
	// if flow stopped the stop requeue
	flowRun := &resource.FlowRun{}
	if err := dsa.Get(common.DefaultNamespace, common.FlowRunCollection, dsa.step.Spec.FlowID, flowRun); err != nil {
		fmt.Printf("[INFO] retry flow run (%s) step (%s) action execute error (%s)", dsa.step.Spec.FlowID, dsa.step.GetName(), err)
		return
	}
	if flowRun.Spec.CurrentState == fsm.STOPPED {
		return
	}
	if flowRun.GetUUID() != dsa.step.Spec.FlowRunUUID {
		fmt.Printf("[INFO] delay step action requeue ignore the (%s.%s.%s)", dsa.step.Spec.FlowID, dsa.step.GetName(), dsa.step.Spec.ActionName)
		return
	}
	_, isUpdate, err := dsa.Apply(common.DefaultNamespace, common.Step, dsa.step.GetName(), dsa.step)
	if err != nil || !isUpdate {
		fmt.Printf("[ERROR] update step error (%s)\n", err)
		return
	}

	fmt.Printf("[INFO] requeue step action (%s) time:(%s) \n", dsa.step.GetName(), t.String())
}
