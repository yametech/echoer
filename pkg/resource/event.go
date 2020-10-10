package resource

import (
	"github.com/yametech/echoer/pkg/core"
	"github.com/yametech/echoer/pkg/storage"
	"github.com/yametech/echoer/pkg/storage/gtm"
	pb "github.com/yametech/echoer/proto"
)

var _ core.IObject = &Event{}

type Event struct {
	core.Metadata `json:"metadata" bson:"metadata"`
	pb.EventType  `json:"event_type" bson:"event_type"`
	Message       string                 `json:"message" bson:"message"`
	Object        map[string]interface{} `json:"object" bson:"object"`
}

func (e *Event) Clone() core.IObject {
	result := &Event{}
	core.Clone(e, result)
	return result
}

var _ storage.Coder = &Event{}

// Event impl Coder
func (*Event) Decode(op *gtm.Op) (core.IObject, error) {
	event := &Event{}
	switch op.Operation {
	case "c", "i":
		event.EventType = pb.EventType_Added
	case "u":
		event.EventType = pb.EventType_Modified
	case "d":
		event.EventType = pb.EventType_Deleted
	}
	if err := core.ObjectToResource(op.Data, event); err != nil {
		return nil, err
	}
	return event, nil
}

func init() {
	storage.AddResourceCoder("event", &Event{})
}
