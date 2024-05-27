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
	IWorkDay interface {
		GetWorkDays(ctx context.Context, domainUuid, userUuid string, limit, offset int, filter model.WorkDayFilter) (int, any)
		GetWorkdayByWorkdayId(ctx context.Context, domainUuid, userUuid, workdayId string) (int, any)
		PostWorkDays(ctx context.Context, domainUuid, userUuid string, workDays []model.WorkDay) (int, any)
		DeleteWorkDay(ctx context.Context, domainUuid, userUuid, workDayUuid string) (int, any)
		DeleteWorkdayByWorkdayId(ctx context.Context, domainUuid, userUuid, workdayId string) (int, any)
		PostHoliday(ctx context.Context, domainUuid, userUuid string, holiday model.Holiday) (int, any)
	}
	WorkDay struct{}
)

func NewWorkDay() IWorkDay {
	s := &WorkDay{}
	return s
}

func (s *WorkDay) GetWorkDays(ctx context.Context, domainUuid, userUuid string, limit, offset int, filter model.WorkDayFilter) (int, any) {
	filter.Limit = limit
	filter.Offset = offset
	total, workdays, err := repository.WorkDayRepo.GetWorkDays(ctx, domainUuid, filter)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.Pagination(workdays, total, limit, offset)
}

func (s *WorkDay) PostWorkDays(ctx context.Context, domainUuid, userUuid string, workDays []model.WorkDay) (int, any) {
	workdayId := uuid.NewString()
	if len(workDays) > 0 {
		for i, item := range workDays {
			item.DomainUuid = domainUuid
			item.WorkdayUuid = uuid.NewString()
			item.WorkDayId = workdayId
			item.WorkDayName = workDays[0].WorkDayName
			item.Status = workDays[0].Status
			item.UnitUuid = workDays[0].UnitUuid
			item.Description = workDays[0].Description
			item.CreatedAt = time.Now()
			item.CreatedBy = userUuid
			workDays[i] = item
		}
	}
	if err := repository.WorkDayRepo.InsertWorkDays(ctx, &workDays); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	// bytes, err := json.Marshal(workDays)
	// if err != nil {
	// 	log.Error(err)
	// 	return response.ServiceUnavailableMsg(err.Error())
	// }
	// go func() {
	// 	auditLog := model.Transaction{
	// 		TransactionUuid: uuid.NewString(),
	// 		DomainUuid:      domainUuid,
	// 		Entity:          "workday",
	// 		Action:          "post",
	// 		EntityUuid:      workDay.WorkdayUuid,
	// 		Status:          "",
	// 		OldData:         "",
	// 		NewData:         string(bytes),
	// 		CreatedAt:       time.Now(),
	// 	}

	// 	if err := repository.AuditLogRepo.InsertNewAuditLog(ctx, &auditLog); err != nil {
	// 		log.Error(err)
	// 	}
	// }()
	return response.Created(map[string]any{
		"total": len(workDays),
	})
}

func (s *WorkDay) DeleteWorkDay(ctx context.Context, domainUuid, userUuid, id string) (int, any) {
	workDay, err := repository.WorkDayRepo.GetWorkDayById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if workDay == nil {
		return response.NotFoundMsg("WorkDay not found")
	}
	if err := repository.WorkDayRepo.DeleteWorkDay(ctx, domainUuid, id); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	go func(workDay *model.WorkDay) {
		auditLog := model.Transaction{
			TransactionUuid: uuid.NewString(),
			DomainUuid:      domainUuid,
			Entity:          "workday",
			Action:          "delete",
			EntityUuid:      workDay.WorkdayUuid,
			Status:          "",
			OldData:         "",
			NewData:         "",
			CreatedAt:       time.Now(),
		}
		if err := repository.TransactionRepo.InsertTransaction(ctx, &auditLog); err != nil {
			log.Error(err)
		}
	}(workDay)

	return response.OK(map[string]any{
		"id": id,
	})
}

func (s *WorkDay) DeleteWorkdayByWorkdayId(ctx context.Context, domainUuid, userUuid, workdayId string) (int, any) {
	workDay, err := repository.WorkDayRepo.GetWorkDayByWorkdayId(ctx, domainUuid, workdayId)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if workDay == nil {
		return response.NotFoundMsg("WorkDay not found")
	}
	if err := repository.WorkDayRepo.DeleteWorkdayByWorkdayId(ctx, domainUuid, workdayId); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"workday_id": workdayId,
	})
}

func (s *WorkDay) PostHoliday(ctx context.Context, domainUuid, userUuid string, holiday model.Holiday) (int, any) {
	holidays := make([]model.WorkDay, 0)
	for d := holiday.StartTime; !d.After(holiday.EndTime); d = d.AddDate(0, 0, 1) {
		startTime := time.Date(d.Year(), d.Month(), d.Day(), 7, 0, 0, 0, time.UTC)
		endTime := time.Date(d.Year(), d.Month(), d.Day(), 21, 30, 0, 0, time.UTC)
		startTimeDuration, _ := time.ParseDuration(startTime.Format("15:04:05"))
		endTimeDuration, _ := time.ParseDuration(endTime.Format("15:04:05"))
		data := &model.WorkDay{
			DomainUuid:  domainUuid,
			WorkdayUuid: uuid.NewString(),
			Day:         d.Format("2006-01-02"),
			StartTime:   startTimeDuration,
			EndTime:     endTimeDuration,
			WorkdayType: holiday.WorkdayType,
			IsWork:      false,
			Offset:      holiday.Offset,
			Description: holiday.Description,
			CreatedBy:   userUuid,
			CreatedAt:   time.Now(),
		}
		holidays = append(holidays, *data)
	}

	if err := repository.WorkDayRepo.InsertWorkDays(ctx, &holidays); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Created(map[string]any{

		"total": len(holidays),
	})
}

func (s *WorkDay) GetWorkdayByWorkdayId(ctx context.Context, domainUuid, userUuid, workdayId string) (int, any) {
	workDay, err := repository.WorkDayRepo.GetWorkDayByWorkdayId(ctx, domainUuid, workdayId)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if workDay == nil {
		return response.NotFoundMsg("work day not found")
	}
	return response.OK(workDay)
}
