package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/uptrace/bun"
)

type TicketCategoryRepo struct{}

func NewTicketCategory() repository.ITicketCategory {
	repo := &TicketCategoryRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *TicketCategoryRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	repository.FusionSqlClient.GetDB().RegisterModel((*model.TicketCategoryInfo)(nil))
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.TicketCategory)(nil)); err != nil {
		panic(err)
	}
}

func (repo *TicketCategoryRepo) InitColumn() {

}

func (repo *TicketCategoryRepo) InitIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.TicketCategory)(nil)).IfNotExists().Index("idx_ticket_category_domain_uuid").Column("domain_uuid").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.TicketCategory)(nil)).IfNotExists().Index("idx_ticket_category_ticket_category_uuid").Column("ticket_category_uuid").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.TicketCategory)(nil)).IfNotExists().Index("idx_ticket_category_code").Column("ticket_category_code").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.TicketCategory)(nil)).IfNotExists().Index("idx_ticket_category_name").Column("ticket_category_name").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.TicketCategory)(nil)).IfNotExists().Index("idx_ticket_category_parent_ticket_category_uuid").Column("parent_ticket_category_uuid").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.TicketCategory)(nil)).IfNotExists().Index("idx_ticket_category_status").Column("status").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.TicketCategory)(nil)).IfNotExists().Index("idx_ticket_category_sla").Column("sla").Exec(ctx); err != nil {
		panic(err)
	}
}

func (repo *TicketCategoryRepo) InsertTicketCategory(ctx context.Context, ticketCategory *model.TicketCategory) error {
	resp, err := repository.FusionSqlClient.GetDB().NewInsert().Model(ticketCategory).Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("insert ticket_category failed")
	}
	return nil
}
func (repo *TicketCategoryRepo) InsertTicketCategoryAndSLAPolicy(ctx context.Context, ticketCategory *model.TicketCategory, slaPolicies *[]model.SlaPolicy) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(ticketCategory).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewInsert().Model(slaPolicies).Exec(ctx); err != nil {
			return err
		}
		return nil
	})
}

func (repo *TicketCategoryRepo) UpdateTicketCategoryAndSLAPolicy(ctx context.Context, ticketCategory *model.TicketCategory, slaPolicies *[]model.SlaPolicyInfo) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewUpdate().Model(ticketCategory).WherePK().Exec(ctx)
		if err != nil {
			return err
		}
		for _, slaPolicicy := range *slaPolicies {
			tmp := slaPolicicy
			existed, err := tx.NewSelect().Model(&slaPolicicy).
				Where("ticket_category_uuid = ?", slaPolicicy.TicketCategoryUuid).
				Where("status = ?", slaPolicicy.Status).
				Where("priority = ?", slaPolicicy.Priority).Exists(ctx)
			if err != nil {
				return err
			}
			if !existed {
				slaPolicicy = tmp
				_, err = tx.NewInsert().Model(&slaPolicicy).
					Exec(ctx)
				if err != nil {
					return err
				}
			} else {
				slaPolicicy = tmp
				_, err := tx.NewUpdate().Model(&slaPolicicy).
					Where("ticket_category_uuid = ?", slaPolicicy.TicketCategoryUuid).
					Where("status = ?", slaPolicicy.Status).
					Where("priority = ?", slaPolicicy.Priority).
					Column("response_time").
					Column("response_type").
					Exec(ctx)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (repo *TicketCategoryRepo) GetTicketCategories(ctx context.Context, domainUuid string, categoryCode string) (*[]model.TicketCategory, error) {
	ticketCategories := new([]model.TicketCategory)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketCategories).
		Where("domain_uuid = ?", domainUuid).
		Where("deleted_at IS NULL")
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return ticketCategories, nil
	}
}

func (repo *TicketCategoryRepo) GetTicketCategoriesInfo(ctx context.Context, domainUuid string, limit, offset int, filter model.TicketCategoryFilter) (*[]model.TicketCategoryInfo, int, error) {
	ticketCategories := new([]model.TicketCategoryInfo)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketCategories).
		ColumnExpr("tc.*").
		ColumnExpr("tca.ticket_category_name as parent_ticket_category_name").
		ColumnExpr("(CASE WHEN tc.parent_ticket_category_uuid IS NULL THEN true ELSE false END) as is_parent").
		Where("tc.domain_uuid = ?", domainUuid).
		Join("LEFT JOIN ticket_category tca ON tc.parent_ticket_category_uuid = tca.ticket_category_uuid").
		Order("tc.parent_ticket_category_uuid DESC", "tc.ticket_category_code", "tc.created_at")
	if len(filter.ParentTicketCategoryUuid) > 0 {
		query = query.Where("tc.parent_ticket_category_uuid = ?", filter.ParentTicketCategoryUuid)
	}
	if len(filter.TicketCategoryName) > 0 {
		query.Where("tc.ticket_category_name = ?", filter.TicketCategoryName)
	}
	if len(filter.TicketCategoryCode) > 0 {
		query.Where("tc.ticket_category_code = ?", filter.TicketCategoryCode)
	}
	if filter.Active.Valid {
		query = query.Where("tc.active = ?", filter.Active.Bool)
	}
	if filter.IsParent.Valid {
		query = query.Where("tc.parent_ticket_category_uuid IS NULL")
	}
	if len(filter.Level) > 0 {
		query = query.Where("tc.level = ?", filter.Level)
	}
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return nil, 0, nil
	} else if err != nil {
		return nil, 0, err
	} else {
		return ticketCategories, total, nil
	}
}
func (repo *TicketCategoryRepo) GetTicketCategoryByCode(ctx context.Context, domainUuid string, ticketCategoryCode string) (*model.TicketCategory, error) {
	ticketCategory := new(model.TicketCategory)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketCategory).
		Where("domain_uuid = ?", domainUuid).
		Limit(1)
	if len(ticketCategoryCode) > 1 {
		query.Where("ticket_category_code = ?", ticketCategoryCode)
	}
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return ticketCategory, nil
	}
}

func (repo *TicketCategoryRepo) GetParentTicketCategoryById(ctx context.Context, domainUuid string, ticketCategoryUuid string) (*model.TicketCategory, error) {
	ticketCategory := new(model.TicketCategory)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketCategory).
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_category_uuid = ?", ticketCategoryUuid).
		// Where("parent_ticket_category_uuid IS NULL").
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return ticketCategory, nil
	}
}

func (repo *TicketCategoryRepo) GetParentTicketCategoriesById(ctx context.Context, domainUuid string, ticketCategoryUuid string) (*[]model.TicketCategory, error) {
	ticketCategory := new([]model.TicketCategory)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketCategory).
		Where("domain_uuid = ?", domainUuid).
		Where("parent_ticket_category_uuid = ?", ticketCategoryUuid)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return ticketCategory, nil
	}
}

func (repo *TicketCategoryRepo) GetTicketCategoryById(ctx context.Context, domainUuid string, ticketCategoryUuid string) (*model.TicketCategory, error) {
	ticketCategory := new(model.TicketCategory)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketCategory).
		Relation("SLAPolicies").
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_category_uuid = ?", ticketCategoryUuid).
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return ticketCategory, nil
	}
}

func (repo *TicketCategoryRepo) UpdateTicketCategory(ctx context.Context, ticketCategory *model.TicketCategory) error {
	resp, err := repository.FusionSqlClient.GetDB().NewUpdate().Model(ticketCategory).
		Where("domain_uuid = ?", ticketCategory.DomainUuid).
		Where("ticket_category_uuid = ?", ticketCategory.TicketCategoryUuid).WherePK().
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("update ticket_category failed")
	}
	return nil
}
func (repo *TicketCategoryRepo) DeleteTicketCategory(ctx context.Context, domainUuid, ticketCategoryUuid string) error {
	category := model.TicketCategory{}
	resp, err := repository.FusionSqlClient.GetDB().NewDelete().
		Model(&category).
		Where("domain_uuid = ? ", domainUuid).
		Where("ticket_category_uuid = ?", ticketCategoryUuid).
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("delete ticket_category failed")
	}
	return nil
}

func (repo *TicketCategoryRepo) GetTicketCategoryInfoById(ctx context.Context, domainUuid string, ticketCategoryUuid string) (*model.TicketCategoryInfo, error) {
	ticketCategory := new(model.TicketCategoryInfo)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ticketCategory).
		Relation("SLAPolicies").
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_category_uuid = ?", ticketCategoryUuid).
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return ticketCategory, nil
	}
}
