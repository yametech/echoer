package controller

import (
	"fmt"
	"time"

	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/core"
	"github.com/yametech/echoer/pkg/fsm"
	"github.com/yametech/echoer/pkg/resource"
	"github.com/yametech/echoer/pkg/storage"
)

var _ Controller = &FlowController{}

type FlowController struct {
	stop chan struct{}
	storage.IStorage
}

func NewFlowController(stage storage.IStorage) *FlowController {
	server := &FlowController{
		stop:     make(chan struct{}),
		IStorage: stage,
	}
	return server
}

func (s *FlowController) Stop() error {
	s.stop <- struct{}{}
	return nil
}

func (s *FlowController) recv() error {
	flowRunObjs, err := s.List(common.DefaultNamespace, common.FlowRunCollection, "")
	if err != nil {
		return err
	}

	flowFlowCoder := storage.GetResourceCoder(string(resource.FlowRunKind))
	if flowFlowCoder == nil {
		return fmt.Errorf("(%s) %s", resource.FlowRunKind, "coder not exist")
	}
	flowRunWatchChan := storage.NewWatch(flowFlowCoder)

	go func() {
		version := int64(0)
		for _, item := range flowRunObjs {
			flowRunObj := &resource.FlowRun{}
			if err := core.UnmarshalInterfaceToResource(&item, flowRunObj); err != nil {
				fmt.Printf("[ERROR] reconcile error %s\n", err)
				continue
			}
			if flowRunObj.GetResourceVersion() > version {
				version = flowRunObj.GetResourceVersion()
			}
			if err := s.reconcileFlowRun(flowRunObj); err != nil {
				fmt.Printf("[ERROR] reconcile flow run (%s) error %s\n", flowRunObj.GetName(), err)
			}
		}
		s.Watch2(common.DefaultNamespace, string(resource.FlowRunKind), version, flowRunWatchChan)

	}()

	stepObjs, err := s.List(common.DefaultNamespace, common.Step, "")
	if err != nil {
		return err
	}
	stepCoder := storage.GetResourceCoder(string(resource.StepKind))
	if stepCoder == nil {
		return fmt.Errorf("(%s) %s", resource.StepKind, "coder not exist")
	}
	stepWatchChan := storage.NewWatch(stepCoder)

	go func() {
		version := int64(0)
		for _, item := range stepObjs {
			stepObj := &resource.Step{}
			if err := core.UnmarshalInterfaceToResource(&item, stepObj); err != nil {
				fmt.Printf("[ERROR] reconcile error %s\n", err)
			}
			if stepObj.GetResourceVersion() > version {
				version = stepObj.GetResourceVersion()
			}
			if err := s.reconcileStep(stepObj); err != nil {
				fmt.Printf("[ERROR] reconcile step (%s) error %s\n", stepObj.GetName(), err)
			}
		}
		s.Watch2(common.DefaultNamespace, common.Step, version, stepWatchChan)
	}()

	for {
		select {
		case <-s.stop:
			flowRunWatchChan.CloseStop() <- struct{}{}
			stepWatchChan.CloseStop() <- struct{}{}
			return nil

		case item, ok := <-flowRunWatchChan.ResultChan():
			if !ok {
				return nil
			}
			fmt.Printf("[INFO] receive flow run (%s) \n", item.GetName())
			flowRunObj := &resource.FlowRun{}
			if err := core.UnmarshalInterfaceToResource(&item, flowRunObj); err != nil {
				fmt.Printf("[ERROR] receive flow run UnmarshalInterfaceToResource error %s\n", err)
				continue
			}
			if err := s.reconcileFlowRun(flowRunObj); err != nil {
				fmt.Printf("[ERROR] receive flow run reconcile error %s\n", err)
			}

		case item, ok := <-stepWatchChan.ResultChan():
			if !ok {
				return nil
			}
			fmt.Printf("[INFO] receive step (%s) \n", item.GetName())
			stepObj := &resource.Step{}
			if err := core.UnmarshalInterfaceToResource(&item, stepObj); err != nil {
				fmt.Printf("[ERROR] receive step UnmarshalInterfaceToResource error %s\n", err)
				continue
			}
			if err := s.reconcileStep(stepObj); err != nil {
				fmt.Printf("[ERROR] receive step reconcile error %s\n", err)
			}
		}
	}
}

func (s *FlowController) reconcileStep(obj *resource.Step) error {
	if obj.GetKind() != resource.StepKind || obj.GetName() == "" {
		return nil
	}
	fmt.Printf("[INFO] start reconcile flow (%s) step (%s) \n", obj.Spec.FlowID, obj.GetName())

	if !obj.Spec.Done {
		fmt.Printf("[INFO] reconcile flow (%s) step (%s) waiting action response \n", obj.Spec.FlowID, obj.GetName())
		return nil
	}
	if obj.Spec.Response.State == "" {
		return fmt.Errorf("flow (%s) step (%s) not response state", obj.Spec.FlowID, obj.GetName())
	}

	flowRun := &resource.FlowRun{}
	if err := s.Get(common.DefaultNamespace, common.FlowRunCollection, obj.Spec.FlowID, flowRun); err != nil {
		return err
	}
	for index, flowRunStep := range flowRun.Spec.Steps {
		if flowRunStep.GetName() != obj.GetName() {
			continue
		}
		flowRun.Spec.Steps[index].Spec.Response.State = obj.Spec.State
		flowRun.Spec.Steps[index].Spec.ActionRun.Done = obj.Spec.Done
		flowRun.Spec.Steps[index].Spec.Data = obj.Spec.Data
	}
	if _, _, err := s.Apply(common.DefaultNamespace, common.FlowRunCollection, flowRun.GetName(), flowRun); err != nil {
		return err
	}

	return nil
}

func (s *FlowController) reconcileFlowRun(obj *resource.FlowRun) error {
	if obj.GetKind() != resource.FlowRunKind || obj.GetName() == "" {
		return nil
	}
	if obj.Spec.CurrentState == fsm.STOPPED || obj.Spec.CurrentState == fsm.DONE {
		return nil
	}
	fmt.Printf("[INFO] reconcile flow run (%s) start \n", obj.GetName())
	startTime := time.Now()
	defer func() {
		fmt.Printf("[INFO] reconcile flow run (%s) end, time (%s)\n", obj.GetName(), time.Now().Sub(startTime))
	}()

	frt, err := CreateFlowRunController(obj.GetName(), &StorageImpl{s.IStorage})
	if err != nil {
		return err
	}
	currentState := frt.Current()
	switch currentState {
	case fsm.STOPPED, fsm.SUSPEND, fsm.DONE:
		return nil
	case fsm.READY:
		if err := frt.Start(); err != nil {
			return err
		}
	default:
		step, err := obj.Spec.GetStepByName(currentState)
		if err != nil {
			return err
		}
		if !step.Spec.Done {
			return nil
		}
		if step.Spec.Response.State == "" {
			return fmt.Errorf("flow (%s) step (%s) not return state", obj.GetName(), step.GetName())
		}
		event := StepStateEvent(frt.GetName(), step.GetName(), step.Spec.Response.State)
		if err := frt.Send(event); err != nil {
			return err
		}
	}

	return nil
}

func (s *FlowController) Run() error { return s.recv() }
