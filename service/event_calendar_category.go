package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type (
	IEventCalendarCategory interface {
		GetEventCalendarCategory(ctx context.Context, domainUuid string, filter model.EventCalendarCategoryFilter, limit, offset int) (int, any)
		InsertEventCalendarCategory(ctx context.Context, domainUuid string, eventCalendarCategory model.EventCalendarCategory) (int, any)
		GetEventCalendarCategoryById(ctx context.Context, domainUuid string, eccUuid string) (int, any)
		UpdateEventCalendarCategoryById(ctx context.Context, domainUuid string, eventCalendarCategory model.EventCalendarCategory) (int, any)
		DeleteEventCalendarCategoryById(ctx context.Context, domainUuid, eccUuid string) (int, any)
	}
	EventCalendarCategory struct{}
)

func NewEventCalendarCategory() IEventCalendarCategory {
	return &EventCalendarCategory{}
}

func (s *EventCalendarCategory) GetEventCalendarCategory(ctx context.Context, domainUuid string, filter model.EventCalendarCategoryFilter, limit, offset int) (int, any) {
	total, data, err := repository.EventCalendarCategoryRepo.GetEventCalendarCategories(ctx, domainUuid, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.Pagination(data, total, limit, offset)
}

func (s *EventCalendarCategory) InsertEventCalendarCategory(ctx context.Context, domainUuid string, eventCalendarCategory model.EventCalendarCategory) (int, any) {
	filter := model.EventCalendarCategoryFilter{
		Title: eventCalendarCategory.Title,
	}
	total, _, err := repository.EventCalendarCategoryRepo.GetEventCalendarCategories(ctx, domainUuid, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if total > 0 {
		err = errors.New(eventCalendarCategory.Title + " is existed")
		return response.ServiceUnavailableMsg(err)
	}

	eventCalendarCategory.DomainUuid = domainUuid
	eventCalendarCategory.EccUuid = uuid.NewString()
	eventCalendarCategory.Status = true
	eventCalendarCategory.CreatedAt = time.Now()
	if err := repository.EventCalendarCategoryRepo.InsertEventCalendarCategory(ctx, domainUuid, eventCalendarCategory); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Created(map[string]any{
		"id":      eventCalendarCategory.EccUuid,
		"message": "created successfull",
	})
}

func (s *EventCalendarCategory) GetEventCalendarCategoryById(ctx context.Context, domainUuid string, eccUuid string) (int, any) {
	eventCalendarCategory, err := repository.EventCalendarCategoryRepo.GetEventCalendarCategoryById(ctx, domainUuid, eccUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(eventCalendarCategory.EccUuid) <= 0 {
		err = errors.New(eventCalendarCategory.Title + " is not exist")
		return response.ServiceUnavailableMsg(err)
	}

	return response.OK(eventCalendarCategory)
}

func (s *EventCalendarCategory) UpdateEventCalendarCategoryById(ctx context.Context, domainUuid string, eventCalendarCategory model.EventCalendarCategory) (int, any) {
	eventCalendarCategoryExist, err := repository.EventCalendarCategoryRepo.GetEventCalendarCategoryById(ctx, domainUuid, eventCalendarCategory.EccUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(eventCalendarCategoryExist.EccUuid) <= 0 {
		err = errors.New(eventCalendarCategory.Title + " is not exist")
		return response.ServiceUnavailableMsg(err)
	}

	eventCalendarCategory.UpdatedAt = time.Now()
	if err := repository.EventCalendarCategoryRepo.UpdateEventCalendarCategoryById(ctx, domainUuid, eventCalendarCategory); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"id":      eventCalendarCategory.EccUuid,
		"message": "updated successfull",
	})
}

func (s *EventCalendarCategory) DeleteEventCalendarCategoryById(ctx context.Context, domainUuid, eccUuid string) (int, any) {
	eventCalendarCategoryExist, err := repository.EventCalendarCategoryRepo.GetEventCalendarCategoryById(ctx, domainUuid, eccUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(eventCalendarCategoryExist.EccUuid) <= 0 {
		err = errors.New(eventCalendarCategoryExist.Title + " is not exist")
		return response.ServiceUnavailableMsg(err)
	}

	if err := repository.EventCalendarCategoryRepo.DeleteEventCalendarCategoryById(ctx, domainUuid, eccUuid); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"id":      eccUuid,
		"message": "updated successfull",
	})
}
