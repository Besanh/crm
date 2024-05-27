package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	"contactcenter-api/repository"
	"context"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type (
	IEventCalendar interface {
		PostEventCalendar(ctx context.Context, domainUuid, userUuid string, eventCalendar model.EventCalendar, eventCalendarTodo []model.EventCalendarTodo, files []*multipart.FileHeader) (int, any)
		GetEventCalendar(ctx context.Context, domainUuid, userUuid string, filter model.EventCalendarFilter) (int, any)
		PutEventCalendarById(ctx context.Context, domainUuid, userUuid, id string, eventCalendar model.EventCalendar) (int, any)
		PatchEventCalendarById(ctx context.Context, domainUuid, userUuid, id string, eventCalendar model.EventCalendar) (int, any)
		GetEventCalendarById(ctx context.Context, domainUuid, userUuid, id string) (int, any)
		PutEventCalendarInfoAndTodoById(ctx context.Context, domainUuid, userUuid, id string, eventCalendar model.EventCalendar) (int, any)
		PutEventCalendarStatusById(ctx context.Context, domainUuid, userUuid, id string, status bool) (int, any)
	}

	EventCalendar struct{}
)

func NewEventCalendar() IEventCalendar {
	return &EventCalendar{}
}

func (s *EventCalendar) PostEventCalendar(ctx context.Context, domainUuid, userUuid string, eventCalendar model.EventCalendar, eventCalendarTodo []model.EventCalendarTodo, files []*multipart.FileHeader) (int, any) {
	// event calendar
	if eventCalendar.RemindTypeEvent == "after" || eventCalendar.RemindTypeEvent == "" {
		now := time.Now()
		endTime := now
		if eventCalendar.RemindType == "minute" {
			endTime = endTime.Add(time.Duration(eventCalendar.RemindTime) * time.Minute)
		} else if eventCalendar.RemindType == "hour" {
			endTime = endTime.Add(time.Duration(eventCalendar.RemindTime) * time.Hour)
		} else if eventCalendar.RemindType == "day" {
			endTime = endTime.AddDate(0, 0, eventCalendar.RemindTime)
		}
		eventCalendar.EndTime = endTime
	}
	if eventCalendar.StartTime.IsZero() {
		eventCalendar.StartTime = time.Now()
	}
	eventCalendar.DomainUuid = domainUuid
	eventCalendar.EcUuid = uuid.NewString()
	eventCalendar.CreatedAt = time.Now()

	// todo
	if len(eventCalendarTodo) > 0 {
		for i, item := range eventCalendarTodo {
			item.DomainUuid = domainUuid
			item.EctUuid = uuid.NewString()
			item.EcUuid = eventCalendar.EcUuid
			item.CreatedBy = userUuid
			item.CreatedAt = time.Now()
			eventCalendarTodo[i] = item
		}
	}

	// attachment
	var eventCalendarAttachment []model.EventCalendarAttachment
	for _, file := range files {
		dir := util.PUBLIC_DIR + "event_calendar_attachments/" + util.TimeToStringLayout(time.Now(), "2006_01_02")
		fileName := util.TimeToStringLayout(time.Now(), "2006_01_02_15_04_05") + "_" + file.Filename
		pathFile := dir + "/" + fileName
		attachment := model.EventCalendarAttachment{
			DomainUuid: domainUuid,
			EcaUuid:    uuid.NewString(),
			EcUuid:     eventCalendar.EcUuid,
			FileName:   fileName,
			PathFile:   pathFile,
			CreatedBy:  userUuid,
			CreatedAt:  time.Now(),
		}
		eventCalendarAttachment = append(eventCalendarAttachment, attachment)
	}

	if err := repository.EventCalendarRepo.InsertEventCalendarTransaction(ctx, &eventCalendar, eventCalendarAttachment, eventCalendarTodo); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Created(map[string]any{

		"id": eventCalendar.EcUuid,
	})
}

func (s *EventCalendar) GetEventCalendar(ctx context.Context, domainUuid, userUuid string, filter model.EventCalendarFilter) (int, any) {
	eventCalendars, err := repository.EventCalendarRepo.GetEventCalendar(ctx, domainUuid, filter)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"data":  eventCalendars,
		"total": len(eventCalendars),
	})
}

func (s *EventCalendar) GetEventCalendarById(ctx context.Context, domainUuid, userUuid, id string) (int, any) {
	eventCalendar, err := repository.EventCalendarRepo.GetEventCalendarById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{

		"data": eventCalendar,
	})
}

func (s *EventCalendar) PutEventCalendarById(ctx context.Context, domainUuid, userUuid, id string, eventCalendar model.EventCalendar) (int, any) {
	eventCalendarExist, err := repository.EventCalendarRepo.GetEventCalendarById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if eventCalendarExist.EcUuid == "" {
		return response.ServiceUnavailableMsg("event calendar not found")
	}
	// event calendar
	if eventCalendar.RemindTypeEvent == "after" {
		now := time.Now()
		endTime := now
		if eventCalendar.RemindType == "minute" {
			endTime = endTime.Add(time.Duration(eventCalendar.RemindTime) * time.Minute)
		} else if eventCalendar.RemindType == "hour" {
			endTime = endTime.Add(time.Duration(eventCalendar.RemindTime) * time.Hour)
		} else if eventCalendar.RemindType == "day" {
			endTime = endTime.AddDate(0, 0, eventCalendar.RemindTime)
		}
		eventCalendar.EndTime = endTime
	}
	eventCalendar.UpdatedBy = userUuid
	eventCalendar.UpdatedAt = time.Now()

	if err := repository.EventCalendarRepo.UpdateEventCalendarById(ctx, domainUuid, eventCalendar); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{

		"id": id,
	})
}

/**
* Use for drag and drop event calendar
 */
func (s *EventCalendar) PatchEventCalendarById(ctx context.Context, domainUuid, userUuid, id string, eventCalendar model.EventCalendar) (int, any) {
	eventCalendarExist, err := repository.EventCalendarRepo.GetEventCalendarById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(eventCalendarExist.EcUuid) < 1 {
		return response.ServiceUnavailableMsg("event calendar not found")
	}
	eventCalendarExist.StartTime = eventCalendar.StartTime
	eventCalendarExist.EndTime = eventCalendar.EndTime
	eventCalendarExist.UpdatedBy = userUuid
	eventCalendarExist.UpdatedAt = time.Now()

	if err := repository.EventCalendarRepo.UpdatePatchEventCalendar(ctx, domainUuid, eventCalendarExist); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{

		"id": eventCalendarExist.EcUuid,
	})
}

func (s *EventCalendar) PutEventCalendarInfoAndTodoById(ctx context.Context, domainUuid, userUuid, id string, eventCalendar model.EventCalendar) (int, any) {
	eventCalendarExist, err := repository.EventCalendarRepo.GetEventCalendarById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(eventCalendarExist.EcUuid) < 1 {
		return response.ServiceUnavailableMsg("event calendar not found")
	}
	eventCalendarExist.Title = eventCalendar.Title
	eventCalendarExist.Description = eventCalendar.Description
	eventCalendarExist.UpdatedBy = userUuid
	eventCalendarExist.UpdatedAt = time.Now()

	if err := repository.EventCalendarRepo.UpdateEventCalendarById(ctx, domainUuid, eventCalendarExist); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	// Todo
	// Delete all or find not exist in slice to delete
	if len(eventCalendar.Todo) > 0 {
		if err := repository.EventCalendarTodoRepo.DeleteEventCalendarTodoByEventId(ctx, domainUuid, eventCalendar.EcUuid); err != nil {
			log.Error()
			return response.ServiceUnavailableMsg(err.Error())
		}

		// Inser/update
		for i, item := range eventCalendar.Todo {
			if len(item.EctUuid) < 1 {
				// Create new
				todo := &model.EventCalendarTodo{
					DomainUuid: domainUuid,
					EctUuid:    uuid.NewString(),
					EcUuid:     id,
					Content:    item.Content,
					IsDone:     item.IsDone,
					CreatedBy:  userUuid,
					CreatedAt:  time.Now(),
				}
				if err := repository.EventCalendarTodoRepo.InsertEventCalendarTodo(ctx, *todo); err != nil {
					log.Error(err)
					continue
				}
			} else {
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
						CreatedBy:  userUuid,
						CreatedAt:  time.Now(),
					}
					if err := repository.EventCalendarTodoRepo.InsertEventCalendarTodo(ctx, *todo); err != nil {
						log.Error(err)
						continue
					}
				} else {
					todo.Content = item.Content
					todo.IsDone = item.IsDone
					todo.UpdatedBy = userUuid
					todo.UpdatedAt = time.Now()
					eventCalendar.Todo[i] = todo
					if err := repository.EventCalendarTodoRepo.UpdateEventCalendarTodo(ctx, domainUuid, *todo); err != nil {
						log.Error(err)
						continue
					}
				}
			}
		}
	}

	return response.OK(map[string]any{

		"id": eventCalendar.EcUuid,
	})
}

func (s *EventCalendar) PutEventCalendarStatusById(ctx context.Context, domainUuid, userUuid, id string, status bool) (int, any) {
	eventCalendarExist, err := repository.EventCalendarRepo.GetEventCalendarById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(eventCalendarExist.EcUuid) < 1 {
		return response.ServiceUnavailableMsg("event calendar not found")
	}
	eventCalendarExist.Status = status
	eventCalendarExist.UpdatedBy = userUuid
	eventCalendarExist.UpdatedAt = time.Now()

	if err := repository.EventCalendarRepo.UpdateStatusEventCalendarById(ctx, domainUuid, eventCalendarExist); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{

		"id": eventCalendarExist.EcUuid,
	})
}
