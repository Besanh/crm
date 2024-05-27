package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type SLAPolicy struct {
}

func NewSLAPolicy() repository.ISLAPolicy {
	if _, err := repository.FusionSqlClient.GetDB().Query(
		`DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'sla_priority') THEN
				CREATE TYPE sla_priority AS ENUM ('normal', 'high', 'urgent');
			END IF;
		END
		$$
		`,
		nil,
	); err != nil {
		panic(err)
	}
	repo := &SLAPolicy{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *SLAPolicy) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.SlaPolicy)(nil)); err != nil {
		panic(err)
	}
}
func (repo *SLAPolicy) InitColumn() {

}
func (repo *SLAPolicy) InitIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.SlaPolicy)(nil)).IfNotExists().Index("idx_sla_policy_domain_uuid").Column("ticket_category_uuid").Exec(ctx); err != nil {
		panic(err)
	}
}

func (repo *SLAPolicy) InsertSLAPolicies(ctx context.Context, slaPolicies []model.SlaPolicy) error {
	query := repository.FusionSqlClient.GetDB().
		NewInsert().
		Model(&slaPolicies)
	resp, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("insert sla_policies failed")
	}
	return nil
}

func (repo *SLAPolicy) GetSLAPolicyOfFilter(ctx context.Context, ticketCategoryUuid, priority, status string) (*model.SlaPolicy, error) {
	slaPolicy := new(model.SlaPolicy)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(slaPolicy).
		Where("ticket_category_uuid = ?", ticketCategoryUuid).
		Where("status = ?", status).
		Where("priority = ?", priority).
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return slaPolicy, nil
	}
}

func (repo *SLAPolicy) GetSLAPolicyByInfo(ctx context.Context, domainUuid, ticketCategoryUuid, status, priority string) (*model.SlaPolicyInfo, error) {
	slaPolicy := new(model.SlaPolicyInfo)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(slaPolicy).
		Where("ticket_category_uuid = ?", ticketCategoryUuid).
		Where("status = ?", status).
		Where("priority = ?", priority)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return slaPolicy, nil
	}
}

func (repo *SLAPolicy) GetSLAPolicyById(ctx context.Context, domainUuid, ticketCategoryUuid string) (*[]model.SlaPolicy, error) {
	slaPolicy := new([]model.SlaPolicy)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(slaPolicy).
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_category_uuid = ?", ticketCategoryUuid).
		Order("created_at DESC")
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return slaPolicy, nil
	}
}

func (repo *SLAPolicy) GetSLAPolicyInfoById(ctx context.Context, domainUuid, ticketCategoryUuid string) (*[]model.SlaPolicyInfo, error) {
	slaPolicy := new([]model.SlaPolicyInfo)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(slaPolicy).
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_category_uuid = ?", ticketCategoryUuid).
		Order("created_at DESC")

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return slaPolicy, nil
	}
}
