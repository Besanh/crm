package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IEventCalendarTodo interface {
	InsertEventCalendarTodo(ctx context.Context, eventCalendarTodo model.EventCalendarTodo) error
	GetEventCalendarTodoById(ctx context.Context, domainUuid, ectUuid string) (*model.EventCalendarTodo, error)
	UpdateEventCalendarTodo(ctx context.Context, domainUuid string, eventCalendarTodo model.EventCalendarTodo) error
	GetEventCalendarTodosByEcUuid(ctx context.Context, domainUuid, ecUuid string) ([]model.EventCalendarTodo, error)
	DeleteEventCalendarTodoById(ctx context.Context, domainUuid, ecUuid string) error
	DeleteEventCalendarTodoByEventId(ctx context.Context, domainUuid, ecUuid string) error
}

var EventCalendarTodoRepo IEventCalendarTodo
