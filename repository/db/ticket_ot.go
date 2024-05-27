package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"github.com/uptrace/bun"
)

type TicketOT struct{}

func NewTicketOT() repository.ITicketOT {
	if _, err := repository.FusionSqlClient.GetDB().NewCreateTable().Model((*model.TicketOT)(nil)).IfNotExists().Exec(context.Background()); err != nil {
		panic(err)
	}
	return &TicketOT{}
}

func (repo *TicketOT) InsertTicketOTs(ctx context.Context, ticketOT *model.TicketOT) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(ticketOT)
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return err
	}
	return nil
}

func (repo *TicketOT) GetTicketOTInDay(ctx context.Context, channel []string) (*[]model.TicketOT, error) {
	ticketOTs := new([]model.TicketOT)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketOTs).
		Where("channel IN (?)", bun.In(channel)).
		Where("is_processed = false").
		Where("((created_at::time >= '08:00' AND created_at::time <= '12:00') OR (created_at::time >= '13:00' AND created_at::time <= '17:30'))").
		Where("date(created_at) = CURRENT_DATE")

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return ticketOTs, nil
}

func (repo *TicketOT) GetTicketOTs(ctx context.Context) (*[]model.TicketOT, error) {
	ticketOTs := new([]model.TicketOT)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketOTs).
		Where("is_processed = false").
		Order("created_at DESC")

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return ticketOTs, nil
}

func (repo *TicketOT) UpdateTicketOT(ctx context.Context, ticketOT *model.TicketOT) error {
	query := repository.FusionSqlClient.GetDB().NewUpdate().Model(ticketOT).WherePK()
	if !ticketOT.IsProcessed {
		query.Set("is_processed = true").
			Set("updated_at = ?", ticketOT.UpdatedAt).
			Set("updated_by = ?", ticketOT.UpdatedBy)
	}
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("update ticketOT failed")
	}
	return nil
}

func (repo *TicketOT) GetTicketOTById(ctx context.Context, ticketOTUuid string) (*model.TicketOT, error) {
	ticketOT := new(model.TicketOT)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketOT).Where("ticket_uuid=?", ticketOTUuid)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return ticketOT, nil
}
