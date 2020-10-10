package controller

import (
	"fmt"
	"github.com/yametech/echoer/pkg/core"

	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/fsm"
	"github.com/yametech/echoer/pkg/resource"
	"github.com/yametech/echoer/pkg/storage"
)

type StorageInterface interface {
	Query(flowRun string) (*FlowRunController, error)
	Update(*resource.FlowRun) error
	Create(step *resource.Step) error
	Record(event resource.Event) error
}

var NotFoundErr = fmt.Errorf("%s", "resource not found")

var _ StorageInterface = &fakeFlowStorage{}

type fakeFlowStorage struct {
	data interface{}
}

func (f fakeFlowStorage) Create(step *resource.Step) error {
	fmt.Printf("create step %v\n", step)
	return nil
}

func (f fakeFlowStorage) Query(s string) (*FlowRunController, error) {
	return nil, fmt.Errorf("%s", "mock_error")
}

func (f fakeFlowStorage) Update(run *resource.FlowRun) error {
	m, err := core.ObjectToMap(run)
	if err != nil {
		return err
	}
	f.data = m
	return nil
}

func (f fakeFlowStorage) Record(event resource.Event) error {
	panic("implement me")
}

var _ StorageInterface = &StorageImpl{}

type StorageImpl struct {
	storage.IStorage
}

func (s StorageImpl) Query(s2 string) (*FlowRunController, error) {
	flowRun := &resource.FlowRun{}
	if err := s.Get(common.DefaultNamespace, common.FlowRunCollection, s2, flowRun); err != nil {
		return nil, err
	}
	if flowRun.Name == "" {
		return nil, NotFoundErr
	}

	currentState := flowRun.Spec.CurrentState
	if currentState == "" {
		currentState = fsm.READY
	}
	frt := &FlowRunController{
		flowRun, fsm.NewFSM(currentState, nil, nil), s,
	}

	for index, step := range frt.Spec.Steps {
		first := false
		last := false
		if index == 0 { //add first
			first = true
		}
		if hasDestAbortState(step.Spec.ReturnStateMap) {
			last = true
		} else if len(frt.Spec.Steps) == index+1 { //add last
			last = true
		}

		stepCopy := step
		if err := frt.stepGraph(stepCopy, first, last); err != nil {
			return nil, err
		}
	}
	return frt, nil
}

func hasDestAbortState(returnStateMap map[string]string) bool {
	for _, v := range returnStateMap {
		if v == fsm.DONE || v == fsm.STOPPED {
			return true
		}
	}
	return false
}

func (s StorageImpl) Update(run *resource.FlowRun) error {
	_, _, err := s.Apply(common.DefaultNamespace, common.FlowRunCollection, run.GetName(), run)
	if err != nil {
		return err
	}
	return nil
}

func (s StorageImpl) Create(step *resource.Step) error {
	_, _, err := s.Apply(common.DefaultNamespace, common.Step, step.GetName(), step)
	if err != nil {
		return err
	}
	return nil
}

func (s StorageImpl) Record(event resource.Event) error {
	panic("implement me")
}
