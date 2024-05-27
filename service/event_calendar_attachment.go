package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	"contactcenter-api/repository"
	"context"
	"mime/multipart"
	"os"
	"time"

	"github.com/google/uuid"
)

type (
	IEventCalendarAttachment interface {
		PostEventCalendarAttachment(ctx context.Context, domainUuid, userUuid, ecUuid, dir string, files []*multipart.FileHeader) (int, any)
		DownloadEventCalendarAttachment(ctx context.Context, id string) (model.EventCalendarAttachment, error)
		DeleteEventCalendarAttachmentById(ctx context.Context, domainUuid, userUuid, id string) (int, any)
	}
	EventCalendarAttachment struct{}
)

func NewEventCalendarAttachment() IEventCalendarAttachment {
	return &EventCalendarAttachment{}
}

func (s *EventCalendarAttachment) PostEventCalendarAttachment(ctx context.Context, domainUuid, userUuid, ecUuid, dir string, files []*multipart.FileHeader) (int, any) {
	var eventCalendarAttachment []model.EventCalendarAttachment
	for _, file := range files {
		fileName := util.TimeToStringLayout(time.Now(), "2006_01_02_15_04_05") + "_" + file.Filename
		pathFile := dir + "/" + fileName
		attachment := model.EventCalendarAttachment{
			DomainUuid: domainUuid,
			EcaUuid:    uuid.NewString(),
			EcUuid:     ecUuid,
			FileName:   fileName,
			PathFile:   pathFile,
			CreatedBy:  userUuid,
			CreatedAt:  time.Now(),
		}
		eventCalendarAttachment = append(eventCalendarAttachment, attachment)
	}

	if err := repository.EventCalendarAttachmentRepo.InsertEventCalendarAttachment(ctx, domainUuid, eventCalendarAttachment); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Created(map[string]any{})
}

func (s *EventCalendarAttachment) DownloadEventCalendarAttachment(ctx context.Context, id string) (model.EventCalendarAttachment, error) {
	attachment, err := repository.EventCalendarAttachmentRepo.GetEventCalendarAttachmentById(ctx, "", id)
	if err != nil {
		log.Error(err)
		return model.EventCalendarAttachment{}, err
	}

	return *attachment, nil
}

func (s *EventCalendarAttachment) DeleteEventCalendarAttachmentById(ctx context.Context, domainUuid, userUuid, id string) (int, any) {
	attachment, err := repository.EventCalendarAttachmentRepo.GetEventCalendarAttachmentById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	if err := repository.EventCalendarAttachmentRepo.DeleteEventCalendarAttachmentById(ctx, domainUuid, id); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	// Delete from os
	if err := os.Remove(attachment.PathFile); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{

		"id": attachment.EcaUuid,
	})
}
