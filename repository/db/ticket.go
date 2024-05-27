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

type Ticket struct{}

func NewTicketRepo() repository.ITicket {
	repo := &Ticket{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *Ticket) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	repository.FusionSqlClient.GetDB().RegisterModel((*model.Ticket)(nil))
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.Ticket)(nil)); err != nil {
		panic(err)
	}
}

func (repo *Ticket) InitColumn() {

}

func (repo *Ticket) InitIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Ticket)(nil)).IfNotExists().Index("idx_ticket_ticket_category_uuid").Column("ticket_category_uuid").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Ticket)(nil)).IfNotExists().Index("idx_ticket_domain_uuid").Column("domain_uuid").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Ticket)(nil)).IfNotExists().Index("idx_ticket_ticket_uuid").Column("ticket_uuid").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Ticket)(nil)).IfNotExists().Index("idx_ticket_customer_id").Column("customer_id").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Ticket)(nil)).IfNotExists().Index("idx_ticket_customer_phone").Column("customer_phone").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Ticket)(nil)).IfNotExists().Index("idx_ticket_assignee_uuid").Column("assignee_uuid").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Ticket)(nil)).IfNotExists().Index("idx_ticket_unit_uuid").Column("unit_uuid").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Ticket)(nil)).IfNotExists().Index("idx_ticket_solution_uuid").Column("solution_uuid").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Ticket)(nil)).IfNotExists().Index("idx_ticket_status").Column("status").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Ticket)(nil)).IfNotExists().Index("idx_ticket_ticket_category_uuid").Column("ticket_category_uuid").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Ticket)(nil)).IfNotExists().Index("idx_ticket_ticket_code").Column("ticket_code").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Ticket)(nil)).IfNotExists().Index("idx_ticket_priority").Column("priority").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Ticket)(nil)).IfNotExists().Index("idx_ticket_channel").Column("channel").Exec(ctx); err != nil {
		panic(err)
	}
}

func (repo *Ticket) InsertTicket(ctx context.Context, ticket *model.Ticket, ticketSla *model.TicketSLA, ticketLog *model.TicketLog, ticketComment *model.TicketComment, checkExistTicketComment bool) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(ticket).Exec(ctx); err != nil {
			return err
		}
		if len(ticketSla.Status) > 0 {
			if _, err := tx.NewInsert().Model(ticketSla).Exec(ctx); err != nil {
				return err
			}
		}
		if _, err := tx.NewInsert().Model(ticketLog).Exec(ctx); err != nil {
			return err
		}
		if checkExistTicketComment {
			if _, err := tx.NewInsert().Model(ticketComment).Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *Ticket) PutTicket(ctx context.Context, domainUuid string, ticket *model.Ticket,
	ticketSlaUpdate *model.TicketSLA, ticketSlaNew *model.TicketSLA, ticketLog *model.TicketLog) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		resp, err := tx.NewUpdate().Model(ticket).WherePK().
			Where("domain_uuid = ?", domainUuid).
			Where("ticket_uuid = ?", ticket.TicketUuid).Exec(ctx)
		if err != nil {
			return err
		} else if affected, _ := resp.RowsAffected(); affected < 1 {
			return nil
		}
		if ticketSlaUpdate != nil {
			resp, err = tx.NewUpdate().Model(ticketSlaUpdate).WherePK().
				Where("domain_uuid = ?", domainUuid).
				Where("ticket_uuid = ?", ticketSlaUpdate.TicketUuid).Exec(ctx)
			if err != nil {
				return err
			} else if affected, _ := resp.RowsAffected(); affected < 1 {
				return nil
			}
		}
		if ticketSlaNew != nil {
			if resp, err = tx.NewInsert().Model(ticketSlaNew).Exec(ctx); err != nil {
				return err
			} else if affected, _ := resp.RowsAffected(); affected < 1 {
				return nil
			}
		}
		if ticketLog != nil {
			if resp, err = tx.NewInsert().Model(ticketLog).Exec(ctx); err != nil {
				return err
			} else if affected, _ := resp.RowsAffected(); affected < 1 {
				return nil
			}
		}
		return nil
	})
}

func (repo *Ticket) GetTicketById(ctx context.Context, domainUuid, ticketId string) (*model.Ticket, error) {
	ticket := new(model.Ticket)
	subQueryUnit := repository.FusionSqlClient.GetDB().NewSelect().TableExpr("unit").
		ColumnExpr("unit.unit_name").
		Where("unit.unit_uuid = ticket.unit_uuid").
		Limit(1)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticket).
		ColumnExpr("ticket.*, unit_tmp.unit_name").
		Relation("TicketCategory").
		Relation("SolutionItem").
		Relation("TicketSla", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Order("created_at ASC")
		}).
		Relation("TicketComment").
		Where("ticket.ticket_uuid = ?", ticketId).
		Join("LEFT JOIN LATERAL (?) unit_tmp ON true", subQueryUnit).Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return ticket, nil
	}
}

func (repo *Ticket) GetTicketsInfo(ctx context.Context, domainUuid string, limit, offset int, ticketFilter *model.TicketFilter) (*[]model.Ticket, int, error) {
	tickets := new([]model.Ticket)
	subQueryUnit := repository.FusionSqlClient.GetDB().NewSelect().TableExpr("unit").
		ColumnExpr("unit.unit_name").
		Where("unit.unit_uuid = ticket.unit_uuid").
		Limit(1)

	query := repository.FusionSqlClient.GetDB().NewSelect().Model(tickets).
		ColumnExpr("ticket.*, unit_tmp.unit_name").
		Relation("TicketCategory").
		Relation("SolutionItem").
		Relation("TicketSla").
		Relation("TicketComment", func(q *bun.SelectQuery) *bun.SelectQuery {
			q.Relation("User")
			return q
		}).
		Relation("User").
		Join("LEFT JOIN LATERAL (?) unit_tmp ON true", subQueryUnit)
	if len(domainUuid) > 1 {
		query.Where("ticket.domain_uuid = ?", domainUuid)
	}
	if len(ticketFilter.Common) > 0 {
		query.WhereOr("? LIKE ?", bun.Ident("ticket.subject"), "%"+ticketFilter.Common+"%").
			WhereOr("? LIKE ?", bun.Ident("ticket.content"), "%"+ticketFilter.Common+"%").
			WhereOr("? LIKE ?", bun.Ident("ticket.solving_content"), "%"+ticketFilter.Common+"%")
	} else {
		if len(ticketFilter.FromDate) > 0 {
			query.Where("ticket.created_at >= ?", ticketFilter.FromDate)
		}
		if len(ticketFilter.ToDate) > 0 {
			query.Where("ticket.created_at <= ?", ticketFilter.ToDate)
		}
		if len(ticketFilter.Status) > 0 {
			ticketFilter.Status = strings.ToLower(ticketFilter.Status)
			query.Where("lower(ticket.status) = ?", ticketFilter.Status)
		}
		if len(ticketFilter.AssigneeUuid) > 0 {
			query.Where("ticket.assignee_uuid = ?", ticketFilter.AssigneeUuid)
		}
		if len(ticketFilter.CategoryUuid) > 0 {
			query.Where("ticket.ticket_category_uuid = ?", ticketFilter.CategoryUuid)
		}
		if len(ticketFilter.CustomerId) > 0 {
			query.Where("ticket.customer_id = ?", ticketFilter.CustomerId)
		}
		if len(ticketFilter.Subject) > 0 {
			query.Where("? LIKE ?", bun.Ident("ticket.subject"), "%"+ticketFilter.Subject+"%")
		}
		if len(ticketFilter.PhoneNumber) > 0 {
			query.Where("? LIKE ?", bun.Ident("ticket.customer_phone"), "%"+ticketFilter.PhoneNumber+"%")
		}
		if len(ticketFilter.CreatedBy) > 0 {
			query.Where("ticket.created_by = ?", ticketFilter.CreatedBy)
		}
		if len(ticketFilter.Priority) > 0 {
			query.Where("ticket.priority = ?", ticketFilter.Priority)
		}
		if len(ticketFilter.Content) > 0 {
			query.Where("? LIKE ?", bun.Ident("ticket.content"), "%"+ticketFilter.Content+"%")
		}
		if len(ticketFilter.SenderId) > 0 {
			query.Where("ticket.sender_id = ?", ticketFilter.SenderId)
		}
		if len(ticketFilter.FullName) > 0 {
			query.Where("? LIKE ?", bun.Ident("ticket.full_name"), "%"+ticketFilter.FullName+"%")
		}
		if len(ticketFilter.TicketCode) > 0 {
			query.Where("ticket.ticket_code = ?", ticketFilter.TicketCode)
		}
	}

	if limit > 0 {
		query.Offset(offset).Limit(limit)
	}
	query.Order("ticket.created_at DESC")

	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return tickets, 0, nil
	} else if err != nil {
		return nil, 0, err
	} else {
		return tickets, total, nil
	}
}

func (repo *Ticket) PutInfoTicket(ctx context.Context, domainUuid string, ticket *model.Ticket) error {
	resp, err := repository.FusionSqlClient.GetDB().NewUpdate().Model(ticket).
		Where("domain_uuid = ?", domainUuid).
		WherePK().Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 0 {
		return nil
	} else {
		return nil
	}
}

func (repo *Ticket) GetLatestTicket(ctx context.Context, domainUuid string) (*model.Ticket, error) {
	currentTime := time.Now().Local().Format("2006-01-02")
	ticket := new(model.Ticket)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticket).
		Where("domain_uuid = ?", domainUuid).
		Where("created_at >= ?", currentTime+" 00:00:00").
		Order("created_at DESC").
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return ticket, nil
	}
}

func (repo *Ticket) UpdateConversationIdTicketByTicketId(ctx context.Context, domainUuid, ticketUuid, conversationId string) error {
	ticket := new(model.TicketInfo)
	query := repository.FusionSqlClient.GetDB().NewUpdate().Model(ticket)
	_, err := query.Set("conversation_id = ?", conversationId).
		Where("ticket_uuid = ?", ticketUuid).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Ticket) GetTicketsExport(ctx context.Context, domainUuid string, limit int, offset int, ticketFilter *model.TicketFilter) (*[]model.TicketExport, int, error) {
	tickets := new([]model.TicketExport)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(tickets).
		// TableExpr("ticket as ticket").
		ColumnExpr("ticket.*")
	if len(domainUuid) > 1 {
		query.Where("ticket.domain_uuid = ?", domainUuid)
	}
	if limit > 0 || offset > 0 {
		query.Limit(limit).Offset(offset)
	}
	if ticketFilter.FromDate != "" && ticketFilter.ToDate != "" {
		query.Where("ticket.created_at >= ? AND ticket.created_at <= ?", ticketFilter.FromDate, ticketFilter.ToDate)
	}

	if ticketFilter.Status != "" {
		ticketFilter.Status = strings.ToLower(ticketFilter.Status)
		query.Where("lower(ticket.status) = ?", ticketFilter.Status)
	}

	if ticketFilter.AssigneeUuid != "" {
		query.Where("ticket.assignee_uuid = ?", ticketFilter.AssigneeUuid)
	}

	if ticketFilter.CategoryUuid != "" {
		query.Where("ticket.ticket_category_uuid = ?", ticketFilter.CategoryUuid)
	}

	if ticketFilter.CustomerId != "" {
		query.Where("ticket.customer_id = ?", ticketFilter.CustomerId)
	}

	if ticketFilter.Subject != "" {
		ticketFilter.Subject = strings.ToLower(ticketFilter.Subject)
		query.Where("lower(ticket.subject) = ?", strings.ToLower(ticketFilter.Subject))
	}

	if ticketFilter.PhoneNumber != "" {
		query.Where("ticket.customer_phone = ?", ticketFilter.PhoneNumber)
	}

	subQueryGetCategory := repository.FusionSqlClient.GetDB().NewSelect().TableExpr("ticket_category as ticket_category").
		ColumnExpr("ticket_category.ticket_category_uuid ,ticket_category.ticket_category_name").
		Where("ticket_category.domain_uuid = ?", domainUuid).
		Where("ticket_category.ticket_category_uuid = ticket.ticket_category_uuid")

	subQueryGetSolution := repository.FusionSqlClient.GetDB().NewSelect().TableExpr("solution as solution").
		ColumnExpr("solution.solution_name").
		Where("solution.solution_uuid = ticket.solution_uuid")

	subQueryGetAssignee := repository.FusionSqlClient.GetDB().NewSelect().TableExpr("v_users as v_users").
		ColumnExpr("v_users.username,v_users.user_uuid").
		Where("v_users.user_uuid = ticket.assignee_uuid").
		Where("v_users.user_enabled= 'true'").
		Limit(1)
	subQueryGetUserCreated := repository.FusionSqlClient.GetDB().NewSelect().TableExpr("v_users as v_users").
		ColumnExpr("v_users.username,v_users.user_uuid").
		Where("v_users.user_uuid = ticket.created_by").
		Where("v_users.user_enabled= 'true'").
		Limit(1)
	subQueryGetUserUpdated := repository.FusionSqlClient.GetDB().NewSelect().TableExpr("v_users as v_users").
		ColumnExpr("v_users.username,v_users.user_uuid").
		Where("v_users.user_uuid = ticket.updated_by").
		Where("v_users.user_enabled= 'true'").
		Limit(1)
	query.ColumnExpr("ticket_category.ticket_category_name as ticket_category_name").Join("LEFT JOIN LATERAL (?) ticket_category ON true", subQueryGetCategory)
	query.ColumnExpr("solution.solution_name as solution_name").Join("LEFT JOIN LATERAL (?) solution ON true", subQueryGetSolution)
	query.ColumnExpr("user_assignee.username as user_assignee").Join("LEFT JOIN LATERAL (?) user_assignee ON true", subQueryGetAssignee)
	query.ColumnExpr("user_created.username as user_created ").Join("LEFT JOIN LATERAL (?) user_created ON true", subQueryGetUserCreated)
	query.ColumnExpr("user_updated.username as user_updated").Join("LEFT JOIN LATERAL (?) user_updated ON true", subQueryGetUserUpdated)

	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return nil, 0, nil
	} else if err != nil {
		return nil, 0, err
	}
	return tickets, total, nil
}

func (repo *Ticket) DeleteTicket(ctx context.Context, domainUuid, ticketUuid string) error {
	ticket := model.Ticket{}
	resp, err := repository.FusionSqlClient.GetDB().NewDelete().
		Model(&ticket).
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_uuid = ?", ticketUuid).
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("delete ticket failed")
	}
	return nil
}

func (repo *Ticket) PatchTicketAttachment(ctx context.Context, domainUuid string, ticket model.Ticket) error {
	resp, err := repository.FusionSqlClient.GetDB().NewUpdate().
		Model(&ticket).
		Where("domain_uuid = ?", domainUuid).
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("update ticket attachment failed")
	}
	return nil
}

func (repo *Ticket) GetLatestTicketEmail(ctx context.Context, domainUuid string) (*model.Ticket, error) {
	currentTime := time.Now().Local().Format("2006-01-02")
	ticket := new(model.Ticket)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticket).
		Where("domain_uuid = ?", domainUuid).
		Where("channel = ?", "email").
		Where("created_at >= ?", currentTime+" 00:00:00").
		Order("created_at DESC").
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return ticket, nil
	}
}

func (repo *Ticket) GetTicketByProfileUuids(ctx context.Context, domainUuid string, profileUuids []string) (*[]model.Ticket, error) {
	result := new([]model.Ticket)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("domain_uuid = ?", domainUuid).
		Where("profile_uuid IN (?)", bun.In(profileUuids))
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}
