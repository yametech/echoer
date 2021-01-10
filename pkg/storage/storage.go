package storage

import (
	"fmt"

	"github.com/yametech/echoer/pkg/core"
	"github.com/yametech/echoer/pkg/storage/gtm"
)

type ErrorType error

var (
	NotFound ErrorType = fmt.Errorf("notFound")
)

var coderList = make(map[string]Coder)

func AddResourceCoder(res string, coder Coder) {
	coderList[res] = coder
}

func GetResourceCoder(res string) Coder {
	coder, exist := coderList[res]
	if !exist {
		return nil
	}
	return coder
}

type Coder interface {
	Decode(*gtm.Op) (core.IObject, error)
}

type WatchInterface interface {
	ResultChan() <-chan core.IObject
	Handle(*gtm.Op) error
	ErrorStop() chan error
	CloseStop() chan struct{}
}

type Watch struct {
	r     chan core.IObject
	err   chan error
	c     chan struct{}
	coder Coder
}

func NewWatch(coder Coder) *Watch {
	return &Watch{
		r:     make(chan core.IObject, 1),
		err:   make(chan error),
		c:     make(chan struct{}),
		coder: coder,
	}
}

// Delegate Handle
func (w *Watch) Handle(op *gtm.Op) error {
	obj, err := w.coder.Decode(op)
	if err != nil {
		return err
	}
	w.r <- obj
	return nil
}

// ResultChan
func (w *Watch) ResultChan() <-chan core.IObject {
	return w.r
}

// ErrorStop
func (w *Watch) CloseStop() chan struct{} {
	return w.c
}

// ErrorStop
func (w *Watch) ErrorStop() chan error {
	return w.err
}

type WatchChan struct {
	ResultChan chan map[string]interface{}
	ErrorChan  chan error
	CloseChan  chan struct{}
}

func NewWatchChan() *WatchChan {
	return &WatchChan{
		ResultChan: make(chan map[string]interface{}, 1),
		ErrorChan:  make(chan error, 1),
		CloseChan:  make(chan struct{}, 1),
	}
}

func (w *WatchChan) Close() {
	w.CloseChan <- struct{}{}
	close(w.ResultChan)
	close(w.CloseChan)
	close(w.ErrorChan)
}

type IStorage interface {
	List(namespace, resource, labels string) ([]interface{}, error)
	Get(namespace, resource, name string, result interface{}) error
	GetByUUID(namespace, resource, uuid string, result interface{}) error
	GetByFilter(namespace, resource string, result interface{}, filter map[string]interface{}) error
	Create(namespace, resource string, object core.IObject) (core.IObject, error)
	Watch(namespace, resource string, resourceVersion int64, watchChan *WatchChan) error
	Watch2(namespace, resource string, resourceVersion int64, watch WatchInterface)
	Apply(namespace, resource, name string, object core.IObject) (core.IObject, bool, error)
	Delete(namespace, resource, name string) error
}
