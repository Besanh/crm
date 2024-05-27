package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IEventCalendarCategory interface {
	InsertEventCalendarCategory(ctx context.Context, domainUuid string, eventCalendarCategory model.EventCalendarCategory) error
	GetEventCalendarCategories(ctx context.Context, domainUuid string, filter model.EventCalendarCategoryFilter, limit, offset int) (int, *[]model.EventCalendarCategory, error)
	GetEventCalendarCategoryById(ctx context.Context, domainUuid, eventCalendarCategoryUuid string) (*model.EventCalendarCategory, error)
	UpdateEventCalendarCategoryById(ctx context.Context, domainUuid string, eventCaEventCalendarCategory model.EventCalendarCategory) error
	DeleteEventCalendarCategoryById(ctx context.Context, domainUuid, eventCalendarCategoryUuid string) error
}

var EventCalendarCategoryRepo IEventCalendarCategory
