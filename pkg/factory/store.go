package factory

import (
	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/resource"
	"github.com/yametech/echoer/pkg/storage"
)

var _ IStore = &StoreImpl{}

type IStore interface {
	GetAction(string) (*resource.Action, error)
	GetFlowRun(string) (*resource.FlowRun, error)
	CreateFlowRun(fr *resource.FlowRun) error
	CreateFlow(fl *resource.Flow) error
	CreateAction(ac *resource.Action) error
}

type StoreImpl struct {
	storage.IStorage
}

func (s *StoreImpl) CreateFlow(fl *resource.Flow) error {
	_, _, err := s.Apply(common.DefaultNamespace, common.FlowCollection, fl.GetName(), fl)
	if err != nil {
		return err
	}
	return nil
}

func (s *StoreImpl) CreateAction(ac *resource.Action) error {
	_, err := s.Create(common.DefaultNamespace, common.ActionCollection, ac)
	if err != nil {
		return err
	}
	return nil
}

func (s *StoreImpl) GetAction(s2 string) (*resource.Action, error) {
	action := &resource.Action{}
	err := s.Get(common.DefaultNamespace, common.ActionCollection, s2, action)
	if err != nil {
		return nil, err
	}
	return action, nil
}

func (s *StoreImpl) GetFlowRun(s2 string) (*resource.FlowRun, error) {
	flowRun := &resource.FlowRun{}
	err := s.Get(common.DefaultNamespace, common.FlowRunCollection, s2, flowRun)
	if err != nil {
		return nil, err
	}
	return flowRun, nil
}

func (s *StoreImpl) CreateFlowRun(fr *resource.FlowRun) error {
	_, _, err := s.Apply(common.DefaultNamespace, common.FlowRunCollection, fr.GetName(), fr)
	if err != nil {
		return err
	}
	return nil
}

func NewStoreImpl(storage storage.IStorage) *StoreImpl {
	return &StoreImpl{
		storage,
	}
}
