package db

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

type EventCalendarRepo struct{}

func NewEventCalendar() repository.IEventCalendar {
	repo := &EventCalendarRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *EventCalendarRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.EventCalendar)(nil)); err != nil {
		panic(err)
	}
}

func (repo *EventCalendarRepo) InitColumn() {

}

func (repo *EventCalendarRepo) InitIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.EventCalendar)(nil)).IfNotExists().Index("idx_event_calendar_domain_uuid").Column("domain_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.EventCalendar)(nil)).IfNotExists().Index("idx_event_calendar_ecc_uuid").Column("ecc_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.EventCalendar)(nil)).IfNotExists().Index("idx_event_calendar_title").Column("title").Exec(ctx); err != nil {
		log.Error(err)
	}
}

func (repo *EventCalendarRepo) InsertEventCalendarTransaction(ctx context.Context, eventCalendar *model.EventCalendar, eventCalendarAttachment []model.EventCalendarAttachment, eventCalendarTodo []model.EventCalendarTodo) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(eventCalendar).Exec(ctx); err != nil {
			return err
		}
		if len(eventCalendarAttachment) > 0 {
			if _, err := tx.NewInsert().Model(&eventCalendarAttachment).Exec(ctx); err != nil {
				return err
			}
		}
		if len(eventCalendarTodo) > 0 {
			if _, err := tx.NewInsert().Model(&eventCalendarTodo).Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *EventCalendarRepo) GetEventCalendar(ctx context.Context, domainUuid string, filter model.EventCalendarFilter) ([]model.EventCalendar, error) {
	result := new([]model.EventCalendar)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("ec.domain_uuid = ?", domainUuid).
		Relation("Category").
		Relation("Todo").
		Relation("Attachment")
	if len(filter.EccUuid) > 0 {
		query.Where("ec.ecc_uuid = ?", filter.EccUuid)
	}
	if len(filter.CreatedBy) > 0 {
		query.Where("ec.created_by = ?", filter.CreatedBy)
	}
	query.Where("ec.start_time >= ?", filter.StartTime).
		Where("ec.end_time <= ?", filter.EndTime)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return *result, nil
}

func (repo *EventCalendarRepo) GetEventCalendarById(ctx context.Context, domainUuid, id string) (model.EventCalendar, error) {
	result := new(model.EventCalendar)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("ec.domain_uuid = ?", domainUuid).
		Where("ec.ec_uuid = ?", id).
		Relation("Category").
		Relation("Todo").
		Relation("Attachment")
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return model.EventCalendar{}, nil
	} else if err != nil {
		return model.EventCalendar{}, err
	}

	return *result, nil
}

func (repo *EventCalendarRepo) UpdateEventCalendarById(ctx context.Context, domainUuid string, eventCalendar model.EventCalendar) error {
	query := repository.FusionSqlClient.GetDB().NewUpdate().Model(&eventCalendar).
		Where("domain_uuid = ?", domainUuid).
		WherePK()
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 0 {
		return nil
	} else {
		return nil
	}
}

func (repo *EventCalendarRepo) UpdatePatchEventCalendar(ctx context.Context, domainUuid string, eventCalendar model.EventCalendar) error {
	res, err := repository.FusionSqlClient.GetDB().NewUpdate().Model(&model.EventCalendar{}).
		Where("domain_uuid = ?", domainUuid).
		Where("ec_uuid = ?", eventCalendar.EcUuid).
		Set("start_time = ?", eventCalendar.StartTime).
		Set("end_time = ?", eventCalendar.EndTime).
		Set("updated_by = ?", eventCalendar.UpdatedBy).
		Set("updated_at = ?", eventCalendar.UpdatedAt).
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 0 {
		return nil
	} else {
		return nil
	}
}

func (repo *EventCalendarRepo) UpdateStatusEventCalendarById(ctx context.Context, domainUuid string, eventCalendar model.EventCalendar) error {
	query := repository.FusionSqlClient.GetDB().NewUpdate().Model(&eventCalendar).
		Where("domain_uuid = ?", domainUuid).
		Where("ec_uuid = ?", eventCalendar.EcUuid).
		Set("status = ?", eventCalendar.Status).
		Set("is_whole_day = ?", eventCalendar.IsWholeDay).
		Set("is_notify_web = ?", eventCalendar.IsNotifyWeb).
		Set("is_notify_email = ?", eventCalendar.IsNotifyEmail).
		Set("is_notify_sms = ?", eventCalendar.IsNotifySms).
		Set("is_notify_zns = ?", eventCalendar.IsNotifyZns).
		Set("is_notify_call = ?", eventCalendar.IsNotifyCall).
		Set("updated_by = ?", eventCalendar.UpdatedBy).
		Set("updated_at = ?", eventCalendar.UpdatedAt)
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 0 {
		return nil
	} else {
		return nil
	}
}
