package service

import (
	"github.com/yametech/echoer/pkg/core"
	"github.com/yametech/echoer/pkg/storage"
	pb "github.com/yametech/echoer/proto"
)

var service *Service

type Service struct {
	storage.IStorage
}

func NewService(stage storage.IStorage) *Service {
	if service == nil {
		service = &Service{stage}
	}
	return service
}

func (s *Service) RecordEvent(evntype pb.EventType, object core.IObject, msg string) error {
	return (&eventService{s}).RecordEvent(evntype, object, msg)
}
