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

type EventCalendarAttachmentRepo struct{}

func NewEventCalendarAttachment() repository.IEventCalendarAttachment {
	repo := &EventCalendarAttachmentRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *EventCalendarAttachmentRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.EventCalendarAttachment)(nil)); err != nil {
		panic(err)
	}
}

func (repo *EventCalendarAttachmentRepo) InitColumn() {

}

func (repo *EventCalendarAttachmentRepo) InitIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.EventCalendarAttachment)(nil)).IfNotExists().Index("idx_event_calendar_attachment_domain_uuid").Column("domain_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.EventCalendarAttachment)(nil)).IfNotExists().Index("idx_event_calendar_attachment_ec_uuid").Column("ec_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
}

func (repo *EventCalendarAttachmentRepo) InsertEventCalendarAttachment(ctx context.Context, domainUuid string, eventCalendatAttachment []model.EventCalendarAttachment) error {
	res, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&eventCalendatAttachment).
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("insert event calendary attachment fail")
	}

	return nil
}

func (repo *EventCalendarAttachmentRepo) DeleteEventCalendarAttachmentById(ctx context.Context, domainUuid, ecaUuid string) error {
	query := repository.FusionSqlClient.GetDB().NewDelete().Model(&model.EventCalendarAttachment{}).
		Where("domain_uuid = ?", domainUuid).
		Where("eca_uuid = ?", ecaUuid)
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("delete event calendary attachment fail")
	}

	return nil
}

func (repo *EventCalendarAttachmentRepo) GetEventCalendarAttachmentById(ctx context.Context, domainUuid, ecaUuid string) (*model.EventCalendarAttachment, error) {
	result := new(model.EventCalendarAttachment)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result)
	if len(domainUuid) > 0 {
		query.Where("domain_uuid = ?", domainUuid)
	}
	query.Where("eca_uuid = ?", ecaUuid)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return result, nil
}
