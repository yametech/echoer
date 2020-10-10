package service

import (
	"github.com/yametech/echoer/pkg/common"
	"github.com/yametech/echoer/pkg/core"
	"github.com/yametech/echoer/pkg/resource"
	pb "github.com/yametech/echoer/proto"
)

type eventService struct {
	*Service
}

func (e *eventService) RecordEvent(evntype pb.EventType, object core.IObject, msg string) error {
	data, err := core.ObjectToMap(object)
	if err != nil {
		return err
	}
	event := &resource.Event{
		Metadata: core.Metadata{
			Kind: "event",
		},
		EventType: evntype,
		Message:   msg,
		Object:    data,
	}
	event.GenerateVersion()

	_, err = e.Create(common.DefaultNamespace, common.EventCollection, event)
	if err != nil {
		return err
	}

	return nil
}
