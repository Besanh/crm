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

type EventCalendarTodoRepo struct{}

func NewEventCalendarTodo() repository.IEventCalendarTodo {
	repo := &EventCalendarTodoRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *EventCalendarTodoRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.EventCalendarTodo)(nil)); err != nil {
		panic(err)
	}
}

func (repo *EventCalendarTodoRepo) InitColumn() {

}

func (repo *EventCalendarTodoRepo) InitIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.EventCalendarTodo)(nil)).IfNotExists().Index("idx_event_calendar_todo_domain_uuid").Column("domain_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.EventCalendarTodo)(nil)).IfNotExists().Index("idx_event_calendar_todo_ec_uuid").Column("ec_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.EventCalendarTodo)(nil)).IfNotExists().Index("idx_event_calendar_todo_ect_uuid").Column("ect_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.EventCalendarTodo)(nil)).IfNotExists().Index("idx_event_calendar_todo_content").Column("content").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.EventCalendarTodo)(nil)).IfNotExists().Index("idx_event_calendar_todo_is_done").Column("is_done").Exec(ctx); err != nil {
		log.Error(err)
	}
}

func (repo *EventCalendarTodoRepo) InsertEventCalendarTodo(ctx context.Context, eventCalendarTodo model.EventCalendarTodo) error {
	res, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&eventCalendarTodo).
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("insert event calendar todo fail")
	}

	return nil
}

func (repo *EventCalendarTodoRepo) GetEventCalendarTodoById(ctx context.Context, domainUuid, ectUuid string) (*model.EventCalendarTodo, error) {
	result := new(model.EventCalendarTodo)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("domain_uuid = ?", domainUuid).
		Where("ect_uuid = ?", ectUuid)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		// return nil, err	// comment for update new item
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *EventCalendarTodoRepo) UpdateEventCalendarTodo(ctx context.Context, domainUuid string, eventCalendarTodo model.EventCalendarTodo) error {
	query := repository.FusionSqlClient.GetDB().NewUpdate().Model(&eventCalendarTodo).
		Where("domain_uuid = ?", domainUuid).
		WherePK()
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return nil
	}
	return nil
}

func (repo *EventCalendarTodoRepo) GetEventCalendarTodosByEcUuid(ctx context.Context, domainUuid, ecUuid string) ([]model.EventCalendarTodo, error) {
	result := new([]model.EventCalendarTodo)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("domain_uuid = ?", domainUuid).
		Where("ec_uuid = ?", ecUuid)
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return *result, nil
}

func (repo *EventCalendarTodoRepo) DeleteEventCalendarTodoById(ctx context.Context, domainUuid, ectUuid string) error {
	query := repository.FusionSqlClient.GetDB().NewDelete().Model(&model.EventCalendarTodo{}).
		Where("domain_uuid = ?", domainUuid).
		Where("ect_uuid = ?", ectUuid)
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("delete event calendar todo fail")
	}

	return nil
}

func (repo *EventCalendarTodoRepo) DeleteEventCalendarTodoByEventId(ctx context.Context, domainUuid, ecUuid string) error {
	query := repository.FusionSqlClient.GetDB().NewDelete().Model(&model.EventCalendarTodo{}).
		Where("domain_uuid = ?", domainUuid).
		Where("ec_uuid = ?", ecUuid)
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("delete event calendar todo fail")
	}

	return nil
}
