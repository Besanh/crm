package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IEventCalendar interface {
	InsertEventCalendarTransaction(ctx context.Context, eventCalendar *model.EventCalendar, eventCalendarAttachment []model.EventCalendarAttachment, eventCalendarTodo []model.EventCalendarTodo) error
	GetEventCalendar(ctx context.Context, domainUuid string, filter model.EventCalendarFilter) ([]model.EventCalendar, error)
	UpdateEventCalendarById(ctx context.Context, domainUuid string, eventCalendar model.EventCalendar) error
	UpdateStatusEventCalendarById(ctx context.Context, domainUuid string, eventCalendar model.EventCalendar) error
	UpdatePatchEventCalendar(ctx context.Context, domainUuid string, eventCalendar model.EventCalendar) error
	GetEventCalendarById(ctx context.Context, domainUuid, id string) (model.EventCalendar, error)
}

var EventCalendarRepo IEventCalendar
