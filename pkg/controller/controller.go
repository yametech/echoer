package controller

import (
	"fmt"
	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/fsm"
	"github.com/yametech/echoer/pkg/resource"
	"github.com/yametech/echoer/pkg/storage"
	"time"
)

type Controller interface {
	Run() error
	Stop() error
}

// Timer is an interface that types implement to schedule and receive OnTimer
// callbacks.
type Timer interface {
	OnTimer(t <-chan struct{})
}

type Queue struct{}

func (q *Queue) Schedule(t Timer, duration time.Duration) {
	<-(time.NewTicker(duration)).C
	sig := make(chan struct{})
	defer func() { close(sig) }()
	go t.OnTimer(sig)
	sig <- struct{}{}
}

var _ Timer = &DelayStepAction{}

type DelayStepAction struct {
	step *resource.Step
	storage.IStorage
}

func (dsa *DelayStepAction) OnTimer(t <-chan struct{}) {
	<-t
	needStop := false
	// check retry count large then step retryCount value then stop requeue
	if dsa.step.Spec.RetryCount >= 3 {
		needStop = true
	}
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

	if needStop {
		dsa.step.Spec.Done = true
		_, isUpdate, err := dsa.Apply(common.DefaultNamespace, common.Step, dsa.step.GetName(), dsa.step)
		if err != nil || !isUpdate {
			fmt.Printf("[ERROR] force update stop flowrun (%s) step (%s) error (%s)\n", flowRun.GetName(), dsa.step.GetName(), err)
			return
		}

		flowRun.Spec.CurrentState = fsm.STOPPED
		_, isUpdate, err = dsa.Apply(common.DefaultNamespace, common.FlowRunCollection, flowRun.GetName(), flowRun)
		if err != nil || !isUpdate {
			fmt.Printf("[ERROR] force update stop flowrun (%s) error (%s)\n", flowRun.GetName(), err)
		}

		fmt.Printf("[WARN] force update stop flowrun (%s) step (%s) because exceed retry count\n", flowRun.GetName(), dsa.step.GetName())

		return
	}

	dsa.step.Spec.RetryCount += 1
	_, isUpdate, err := dsa.Apply(common.DefaultNamespace, common.Step, dsa.step.GetName(), dsa.step)
	if err != nil || !isUpdate {
		fmt.Printf("[ERROR] update step error (%s)\n", err)
		return
	}

	fmt.Printf("[INFO] requeue step action (%s) \n", dsa.step.GetName())
}
