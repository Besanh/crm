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

type TicketLog struct{}

func NewTicketLogRepo() repository.ITicketLog {
	repo := &TicketLog{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *TicketLog) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	repository.FusionSqlClient.GetDB().RegisterModel((*model.TicketLog)(nil))
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.TicketLog)(nil)); err != nil {
		panic(err)
	}
}

func (repo *TicketLog) InitColumn() {
}

func (repo *TicketLog) InitIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.TicketLog)(nil)).IfNotExists().Index("idx_ticket_log_domain_uuid").Column("domain_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
}

func (repo *TicketLog) InsertNewTicketLog(ctx context.Context, ticketLog model.TicketLog) error {
	resp, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&ticketLog).Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("insert ticket_log failed")
	}
	return nil
}

func (repo *TicketLog) GetTicketLogs(ctx context.Context, domainUuid, ticketUuid, ticketLogType, status string, offset, limit int) (*[]model.TicketLog, int, error) {
	ticketLogs := new([]model.TicketLog)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketLogs).
		Where("domain_uuid = ?", domainUuid).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC")
	if len(ticketUuid) > 1 {
		query = query.Where("ticket_uuid = ?", ticketUuid)
	}
	if len(status) > 1 {
		query = query.Where("status = ?", status)
	}
	if len(ticketLogType) > 1 {
		query = query.Where("type = ?", ticketLogType)
	}
	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return nil, 0, nil
	} else if err != nil {
		return nil, 0, err
	} else {
		return ticketLogs, total, nil
	}
}

func (repo *TicketLog) GetTicketLogsByType(ctx context.Context, domainUuid string, ticketUuid string, ticketLogType string) (*[]model.TicketLog, error) {
	ticketLogs := new([]model.TicketLog)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketLogs).
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_uuid = ?", ticketUuid)
	if len(ticketLogType) > 1 {
		query = query.Where("type = ?", ticketLogType)
	}
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return ticketLogs, nil
	}
}

func (repo *TicketLog) DeleteTicketLog(ctx context.Context, domainUuid, ticketUuid string) error {
	ticketLog := new([]model.TicketLog)
	_, err := repository.FusionSqlClient.GetDB().NewDelete().
		Model(ticketLog).
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_uuid = ?", ticketUuid).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
