package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IEventCalendarAttachment interface {
	InsertEventCalendarAttachment(ctx context.Context, domainUuid string, eventCalendatAttachment []model.EventCalendarAttachment) error
	DeleteEventCalendarAttachmentById(ctx context.Context, domainUuid, ecaUuid string) error
	GetEventCalendarAttachmentById(ctx context.Context, domainUuid, ecaUuid string) (*model.EventCalendarAttachment, error)
}

var EventCalendarAttachmentRepo IEventCalendarAttachment
