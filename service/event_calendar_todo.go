package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type (
	IEventCalendarTodo interface {
		PutEventCalendarTodos(ctx context.Context, domainUuid, userUuid, ecUuid string, eventCalendarTodo []model.EventCalendarTodo) (int, any)
		PutEventCalendarTodoById(ctx context.Context, domainUuid, userUuid, ectUuid string, eventCalendarTodo model.EventCalendarTodo) (int, any)
		DeleteEventCalendarTodo(ctx context.Context, domainUuid, ectUuid string) (int, any)
	}
	EventCalendarTodo struct{}
)

func NewEventCalendarTodo() IEventCalendarTodo {
	return &EventCalendarTodo{}
}

func (s *EventCalendarTodo) PutEventCalendarTodos(ctx context.Context, domainUuid, userUuid, ecUuid string, eventCalendarTodo []model.EventCalendarTodo) (int, any) {
	if len(eventCalendarTodo) > 0 {
		// Delete all or find not exist in slice to delete
		if err := repository.EventCalendarTodoRepo.DeleteEventCalendarTodoByEventId(ctx, domainUuid, ecUuid); err != nil {
			log.Error()
			return response.ServiceUnavailableMsg(err.Error())
		}

		// Inser/update
		for _, item := range eventCalendarTodo {
			todo, err := repository.EventCalendarTodoRepo.GetEventCalendarTodoById(ctx, domainUuid, item.EctUuid)
			if err != nil {
				log.Error(err)
				continue
			} else if todo == nil {
				// Create new
				todo = &model.EventCalendarTodo{
					DomainUuid: domainUuid,
					EctUuid:    uuid.NewString(),
					EcUuid:     item.EcUuid,
					Content:    item.Content,
					IsDone:     item.IsDone,
				}
				if err := repository.EventCalendarTodoRepo.InsertEventCalendarTodo(ctx, *todo); err != nil {
					log.Error(err)
					continue
				}
			}
		}
	}

	return response.OK(map[string]any{})
}

func (s *EventCalendarTodo) PutEventCalendarTodoById(ctx context.Context, domainUuid, userUuid, ectUuid string, eventCalendarTodo model.EventCalendarTodo) (int, any) {
	eventCalendarTodoExist, err := repository.EventCalendarTodoRepo.GetEventCalendarTodoById(ctx, domainUuid, ectUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if eventCalendarTodoExist == nil {
		return response.ServiceUnavailableMsg("event calendar todo not found")
	}

	eventCalendarTodo.UpdatedBy = userUuid
	eventCalendarTodo.UpdatedAt = time.Now()
	if err := repository.EventCalendarTodoRepo.UpdateEventCalendarTodo(ctx, domainUuid, eventCalendarTodo); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"id": ectUuid,
	})
}

func (s *EventCalendarTodo) DeleteEventCalendarTodo(ctx context.Context, domainUuid, ectUuid string) (int, any) {
	eventCalendarTodo, err := repository.EventCalendarTodoRepo.GetEventCalendarTodoById(ctx, domainUuid, ectUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if eventCalendarTodo == nil {
		return response.ServiceUnavailableMsg("event calendar todo not found")
	}

	if err := repository.EventCalendarTodoRepo.DeleteEventCalendarTodoById(ctx, domainUuid, ectUuid); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{

		"id": ectUuid,
	})
}
