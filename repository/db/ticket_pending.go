package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type TicketPending struct{}

func NewTicketPendingRepo() repository.ITicketPending {
	repo := &TicketPending{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()
	return repo
}

func (repo *TicketPending) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	repository.FusionSqlClient.GetDB().RegisterModel((*model.TicketPending)(nil))
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.TicketPending)(nil)); err != nil {
		panic(err)
	}
}

func (repo *TicketPending) InitColumn() {

}

func (repo *TicketPending) InitIndex() {

}

func (repo *TicketPending) InsertTicketPending(ctx context.Context, TicketPending *model.TicketPending) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(TicketPending)
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return err
	}
	return nil
}

func (repo *TicketPending) GetTicketPendingInDay(ctx context.Context, channel []string) (*[]model.TicketPending, error) {
	ticketPendings := new([]model.TicketPending)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketPendings).
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
	return ticketPendings, nil
}

func (repo *TicketPending) GetTicketPendings(ctx context.Context, domainUuid string, limit, offset int, filter model.TicketFilter) (*[]model.TicketPending, int, error) {
	ticketPending := new([]model.TicketPending)
	subQueryUnit := repository.FusionSqlClient.GetDB().NewSelect().TableExpr("unit").
		ColumnExpr("unit.unit_name").
		Where("unit.unit_uuid = tp.unit_uuid::uuid").
		Limit(1)

	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketPending).
		ColumnExpr("tp.*, unit_tmp.unit_name").
		Relation("TicketCategory").
		Relation("User").
		Join("LEFT JOIN LATERAL (?) unit_tmp ON true", subQueryUnit)
	if len(domainUuid) > 1 {
		query = query.Where("tp.domain_uuid = ?", domainUuid)
	}
	if len(filter.FromDate) > 0 {
		query = query.Where("tp.created_at >= ?", filter.FromDate)
	}
	if len(filter.ToDate) > 0 {
		query = query.Where("tp.created_at <= ?", filter.ToDate)
	}
	if len(filter.Status) > 0 {
		filter.Status = strings.ToLower(filter.Status)
		query = query.Where("lower(tp.status) = ?", filter.Status)
	}
	if len(filter.AssigneeUuid) > 0 {
		query = query.Where("tp.assignee_uuid = ?", filter.AssigneeUuid)
	}
	if len(filter.CategoryUuid) > 0 {
		query = query.Where("tp.ticket_category_uuid = ?", filter.CategoryUuid)
	}
	if len(filter.CustomerId) > 0 {
		query = query.Where("tp.customer_id = ?", filter.CustomerId)
	}
	if len(filter.Subject) > 0 {
		query = query.Where("? LIKE ?", bun.Ident("tp.subject"), "%"+filter.Subject+"%")
	}
	if len(filter.PhoneNumber) > 0 {
		query = query.Where("? LIKE ?", bun.Ident("tp.customer_phone"), "%"+filter.PhoneNumber+"%")
	}
	if len(filter.CreatedBy) > 0 {
		query = query.Where("tp.created_by = ?", filter.CreatedBy)
	}
	if len(filter.Priority) > 0 {
		query = query.Where("tp.priority = ?", filter.Priority)
	}
	if len(filter.Content) > 0 {
		query = query.Where("? LIKE ?", bun.Ident("tp.content"), "%"+filter.Content+"%")
	}
	if len(filter.TicketCode) > 0 {
		query = query.Where("tp.ticket_code = ?", filter.TicketCode)
	}
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	query = query.Order("tp.created_at DESC")
	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return ticketPending, 0, nil
	} else if err != nil {
		return nil, 0, err
	} else {
		return ticketPending, total, nil
	}
}

func (repo *TicketPending) UpdateTicketPending(ctx context.Context, ticketPendings *model.TicketPending) error {
	query := repository.FusionSqlClient.GetDB().NewUpdate().Model(ticketPendings).WherePK()
	if !ticketPendings.IsProcessed {
		query.Set("is_processed = true").
			Set("updated_at = ?", ticketPendings.UpdatedAt).
			Set("updated_by = ?", ticketPendings.UpdatedBy)
	}
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("update ticket pending failed")
	}
	return nil
}

func (repo *TicketPending) GetTicketPendingById(ctx context.Context, ticketPendingUuid string) (*model.TicketPending, error) {
	result := new(model.TicketPending)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).Where("ticket_uuid = ?", ticketPendingUuid)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *TicketPending) DeleteTicketPendingById(ctx context.Context, ticketPendingUuid string) error {
	query := repository.FusionSqlClient.GetDB().NewDelete().Model((*model.TicketPending)(nil)).Where("ticket_uuid = ?", ticketPendingUuid)
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return nil
	}
	return nil
}
