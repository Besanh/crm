package db

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type EventCalendarCategoryRepo struct{}

func NewEventCalendarCategoryRepo() repository.IEventCalendarCategory {
	repo := &EventCalendarCategoryRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *EventCalendarCategoryRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.EventCalendarCategory)(nil)); err != nil {
		panic(err)
	}
}

func (repo *EventCalendarCategoryRepo) InitColumn() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.EventCalendarCategory)(nil)).IfNotExists().Index("idx_event_calendar_category_domain_uuid").Column("domain_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.EventCalendarCategory)(nil)).IfNotExists().Index("idx_event_calendar_category_ecc_uuid").Column("ecc_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.EventCalendarCategory)(nil)).IfNotExists().Index("idx_event_calendar_category_title").Column("title").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.EventCalendarCategory)(nil)).IfNotExists().Index("idx_event_calendar_category_color").Column("color").Exec(ctx); err != nil {
		log.Error(err)
	}
}

func (repo *EventCalendarCategoryRepo) InitIndex() {

}

func (repo *EventCalendarCategoryRepo) InsertEventCalendarCategory(ctx context.Context, domainUuid string, eventCalendarCategory model.EventCalendarCategory) error {
	res, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&eventCalendarCategory).
		Where("domain_uuid = ?", domainUuid).
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("insert event calendary category fail")
	}

	return nil
}

func (repo *EventCalendarCategoryRepo) GetEventCalendarCategories(ctx context.Context, domainUuid string, filter model.EventCalendarCategoryFilter, limit, offset int) (int, *[]model.EventCalendarCategory, error) {
	result := new([]model.EventCalendarCategory)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("domain_uuid = ?", domainUuid)
	if len(filter.Title) > 0 {
		query.Where("title = ?", filter.Title)
	}
	if len(filter.Color) > 0 {
		query.Where("color = ?", filter.Color)
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status.Bool)
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return 0, result, err
	} else if err == sql.ErrNoRows {
		return 0, result, nil
	}

	return total, result, nil
}

func (repo *EventCalendarCategoryRepo) GetEventCalendarCategoryById(ctx context.Context, domainUuid, eventCalendarCategoryUuid string) (*model.EventCalendarCategory, error) {
	result := new(model.EventCalendarCategory)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("domain_uuid = ?", domainUuid).
		Where("ecc_uuid = ?", eventCalendarCategoryUuid)
	err := query.Scan(ctx)
	if err != nil {
		return result, err
	} else if err == sql.ErrNoRows {
		return result, nil
	}

	return result, nil
}

func (repo *EventCalendarCategoryRepo) UpdateEventCalendarCategoryById(ctx context.Context, domainUuid string, eventCaEventCalendarCategory model.EventCalendarCategory) error {
	result := new(model.EventCalendarCategory)
	query := repository.FusionSqlClient.GetDB().NewUpdate().Model(result).
		Where("domain_uuid = ?", domainUuid).
		Where("ecc_uuid = ?", eventCaEventCalendarCategory.EccUuid).
		Set("title = ?", eventCaEventCalendarCategory.Title).
		Set("color = ?", eventCaEventCalendarCategory.Color).
		Set("status = ?", eventCaEventCalendarCategory.Status)
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return nil
	}
	return nil
}

func (repo *EventCalendarCategoryRepo) DeleteEventCalendarCategoryById(ctx context.Context, domainUuid, eventCalendarCategoryUuid string) error {
	_, err := repository.FusionSqlClient.GetDB().NewDelete().Model((*model.EventCalendarCategory)(nil)).
		Where("domain_uuid = ?", domainUuid).
		Where("ecc_uuid = ?", eventCalendarCategoryUuid).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
