package controller

import (
	"fmt"
	"time"

	"github.com/yametech/echoer/pkg/action"
	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/core"
	"github.com/yametech/echoer/pkg/resource"
	"github.com/yametech/echoer/pkg/storage"
)

var _ Controller = &ActionController{}

type ActionController struct {
	stop chan struct{}
	storage.IStorage
	//act    action.Interface
	tqStop chan struct{}
	tq     *Queue
}

func NewActionController(stage storage.IStorage) *ActionController {
	server := &ActionController{
		stop:     make(chan struct{}),
		IStorage: stage,
		//act:      act,
		tqStop: make(chan struct{}),
		tq:     &Queue{},
	}
	return server
}

func (a *ActionController) Stop() error {
	a.tqStop <- struct{}{}
	a.stop <- struct{}{}
	return nil
}

func (a *ActionController) Run() error { return a.recv() }

func (a *ActionController) recv() error {
	stepObjs, err := a.List(common.DefaultNamespace, common.Step, "")
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
			if err := a.realAction(stepObj); err != nil {
				fmt.Printf("[ERROR] reconcile step (%s) error %s\n", stepObj.GetName(), err)
			}
		}
		a.Watch2(common.DefaultNamespace, common.Step, version, stepWatchChan)
	}()

	for {
		select {
		case <-a.stop:
			stepWatchChan.CloseStop() <- struct{}{}
			return nil

		case item, ok := <-stepWatchChan.ResultChan():
			if !ok {
				return nil
			}
			if item.GetName() == "" {
				continue
			}
			stepObj := &resource.Step{}
			if err := core.UnmarshalInterfaceToResource(&item, stepObj); err != nil {
				fmt.Printf("[ERROR] receive step UnmarshalInterfaceToResource error %s\n", err)
				continue
			}
			fmt.Printf("[INFO] receive step (%s) flowID (%s) \n", item.GetName(), stepObj.Spec.FlowID)
			if err := a.realAction(stepObj); err != nil {
				fmt.Printf("[ERROR] receive flow run (%s) step (%s) reconcile error (%s)\n", stepObj.Spec.FlowID, stepObj.GetName(), err)
			}
		}
	}
}

func (a *ActionController) realAction(obj *resource.Step) error {
	if obj.GetKind() != resource.StepKind {
		return nil
	}
	if obj.Spec.Done {
		fmt.Printf("[INFO] real action reconcile step (%s) flowrun (%s) done\n", obj.GetName(), obj.Spec.FlowID)
		return nil
	}

	fmt.Printf("[INFO] real action reconcile step (%s) flowrun (%s) action (%s) \n", obj.GetName(), obj.Spec.FlowID, obj.Spec.ActionName)

	_action := &resource.Action{}
	if err := a.Get(common.DefaultNamespace, common.ActionCollection, obj.Spec.ActionName, _action); err != nil {
		return err
	}

	if err := resource.CheckActionParams(obj.Spec.ActionParams, _action.Spec.Params); err != nil {
		return err
	}

	_flowRun := &resource.FlowRun{}
	if err := a.Get(common.DefaultNamespace, common.FlowRunCollection, obj.Spec.FlowID, _flowRun); err != nil {
		return err
	}

	obj.Spec.ActionParams[common.FlowId] = obj.Spec.FlowID
	obj.Spec.ActionParams[common.StepName] = obj.GetName()
	obj.Spec.ActionParams[common.AckStates] = _action.Spec.ReturnStates
	obj.Spec.ActionParams[common.UUID] = obj.UUID
	obj.Spec.ActionParams[common.GlobalVariables] = _flowRun.Spec.GlobalVariables
	obj.Spec.ActionParams[common.CaKey] = _action.Spec.CaKey
	obj.Spec.ActionParams[common.CaPEM] = _action.Spec.CaPEM

	switch _action.Spec.ServeType {
	case resource.HTTP:
		go func() {
			err := action.NewHookClient().
				HttpInterface().
				Params(obj.Spec.ActionParams).
				Post(_action.Spec.Endpoints).
				Do()

			retryCount := time.Duration(obj.Spec.RetryCount)
			if err != nil {
				fmt.Printf(
					"[ERROR] flowrun (%s) step (%s) execute action (%s) error: %s\n",
					obj.Spec.FlowID,
					obj.GetName(),
					obj.Spec.ActionName,
					err,
				)
				a.tq.Schedule(
					&DelayStepAction{obj, a.IStorage},
					retryCount*time.Second,
				)
			}
		}()
	case resource.HTTPS:
		go func() {
			err := action.NewHookClient().
				HttpsInterface().
				Params(obj.Spec.ActionParams).
				Post(_action.Spec.Endpoints).
				Do()

			retryCount := time.Duration(obj.Spec.RetryCount)
			if err != nil {
				fmt.Printf(
					"[ERROR] flowrun (%s) step (%s) execute action (%s) error: %s\n",
					obj.Spec.FlowID,
					obj.GetName(),
					obj.Spec.ActionName,
					err,
				)
				a.tq.Schedule(
					&DelayStepAction{obj, a.IStorage},
					retryCount*time.Second,
				)
			}
		}()
	case resource.GRPC:
		// TODO current unsupported grpc
	}

	return nil
}
