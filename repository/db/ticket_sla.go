package db

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"errors"
	"time"

	"github.com/uptrace/bun"
)

type TicketSLA struct{}

func NewTicketSLARepo() repository.ITicketSLA {
	repo := &TicketSLA{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *TicketSLA) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	repository.FusionSqlClient.GetDB().RegisterModel((*model.TicketSLA)(nil))
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.TicketSLA)(nil)); err != nil {
		panic(err)
	}
}

func (repo *TicketSLA) InitColumn() {

}

func (repo *TicketSLA) InitIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.TicketSLA)(nil)).IfNotExists().Index("idx_ticket_sla_domain_uuid").Column("domain_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
}

func (repo *TicketSLA) UpdateTicketSLA(ctx context.Context, ticketSLA *model.TicketSLA) error {
	resp, err := repository.FusionSqlClient.GetDB().NewUpdate().Model(ticketSLA).WherePK().
		Where("domain_uuid = ?", ticketSLA.DomainUuid).
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("insert ticket_sla failed")
	}
	return nil
}

func (repo *TicketSLA) GetTicketSLAByTicketId(ctx context.Context, domainUuid, ticketUuid string, state string) (*model.TicketSLA, error) {
	ticketSla := new(model.TicketSLA)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketSla).
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_uuid = ?", ticketUuid)
	if len(state) > 0 {
		query = query.Where("stage = ?", state)
	}
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return ticketSla, nil
}

func (repo *TicketSLA) UpdateTicketSLAByTicketIds(ctx context.Context, ticketUuids []string) error {
	if len(ticketUuids) > 0 {
		ticketSla := new(model.TicketSLA)
		resp, err := repository.FusionSqlClient.GetDB().NewUpdate().Model(ticketSla).
			Set("status = ?", "STOP").
			Where("ticket_uuid IN (?)", bun.In(ticketUuids)).Exec(ctx)
		if err != nil {
			return err
		} else if affected, _ := resp.RowsAffected(); affected < 1 {
			return errors.New("update ticket_sla failed")
		}
		return nil
	}
	return nil
}

func (repo *TicketSLA) GetTicketSLAs(ctx context.Context, domainUuid, ticketUuid string) (*[]model.TicketSLA, error) {
	ticketSla := new([]model.TicketSLA)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketSla).
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_uuid = ?", ticketUuid).
		Order("created_at DESC")

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return ticketSla, nil
}

func (repo *TicketSLA) GetTicketSLAById(ctx context.Context, domainUuid, ticketSlaUuid string) (*model.TicketSLA, error) {
	ticketSla := new(model.TicketSLA)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketSla).
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_sla_uuid = ?", ticketSlaUuid)

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return ticketSla, nil
}

func (repo *TicketSLA) GetTicketSLALatestByTicketId(ctx context.Context, domainUuid, ticketUuid string) (*model.TicketSLA, error) {
	ticketSla := new(model.TicketSLA)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketSla).
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_uuid = ?", ticketUuid).
		Order("created_at DESC").
		Limit(1)

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return ticketSla, nil
}

func (repo *TicketSLA) GetTicketSlaByStatus(ctx context.Context, domainUuid, ticketSlaUuid, status, stage string) (*model.TicketSLA, error) {
	ticketSla := new(model.TicketSLA)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketSla).
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_uuid = ?", ticketSlaUuid).
		Where("status = (?)", status)
	if stage != "" {
		query = query.Where("stage = (?)", stage)
	}
	query = query.Order("created_at DESC").
		Limit(1)
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return ticketSla, nil
}

func (repo *TicketSLA) GetTicketSlaLatestByStatusTicketAndId(ctx context.Context, domainUuid, ticketSlaUuid, status string) (*model.TicketSLA, error) {
	ticketSla := new(model.TicketSLA)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketSla).
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_uuid = ?", ticketSlaUuid).
		Where("ticket_status = (?)", status)
	query = query.Order("created_at DESC").
		Limit(1)
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return ticketSla, nil
}

func (repo *TicketSLA) DeleteTicketSla(ctx context.Context, domainUuid, ticketUuid string) error {
	ticketSla := new([]model.TicketSLA)
	_, err := repository.FusionSqlClient.GetDB().NewDelete().
		Model(ticketSla).
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_uuid = ?", ticketUuid).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
