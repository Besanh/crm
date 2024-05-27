package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IWorkDay interface {
	GetWorkDays(ctx context.Context, domainUuid string, filter model.WorkDayFilter) (int, *[]model.WorkDay, error)
	GetWorkDayByWorkdayId(ctx context.Context, domainUuid, workdayId string) (*[]model.WorkDay, error)
	GetWorkDayById(ctx context.Context, domainUuid, id string) (*model.WorkDay, error)
	InsertWorkDays(ctx context.Context, workDays *[]model.WorkDay) error
	DeleteWorkDay(ctx context.Context, domainUuid, id string) error
	DeleteWorkdayByWorkdayId(ctx context.Context, domainUuid, workdayId string) error
}

var WorkDayRepo IWorkDay
