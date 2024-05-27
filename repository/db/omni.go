package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/common/model/omni"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

type Omni struct {
}

func NewOmniRepo() repository.IOmni {
	repo := &Omni{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()
	return repo
}

func (repo *Omni) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*omni.Omni)(nil)); err != nil {
		panic(err)
	}
}

func (repo *Omni) InitColumn() {
}

func (repo *Omni) InitIndex() {
}

func (repo *Omni) GetOmnis(ctx context.Context, domainUuid string, limit, offset int, filter model.OmniFilter) ([]omni.Omni, int, error) {
	var omnis []omni.Omni
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(&omnis).Where("domain_uuid = ?", domainUuid).
		Where("domain_uuid = ?", domainUuid)
	if len(filter.OmniName) > 0 {
		query.Where("? LIKE ?", bun.Ident("omni_name"), "%"+filter.OmniName+"%")
	}
	if len(filter.OmniType) > 0 {
		query.Where("omni_type = ?", filter.OmniType)
	}
	if len(filter.Supplier) > 0 {
		query.Where("supplier = ?", filter.Supplier)
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status.Bool)
	}
	if !filter.StartTime.IsZero() {
		query.Where("c.created_at >= ?", filter.StartTime)
	}
	if !filter.EndTime.IsZero() {
		query.Where("c.created_at <= ?", filter.EndTime)
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return omnis, 0, nil
	}
	return omnis, total, err
}

func (repo *Omni) GetOmniById(ctx context.Context, domainUuid, id string) (omni.Omni, error) {
	omni := new(omni.Omni)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(omni).
		Where("domain_uuid = ?", domainUuid).
		Where("omni_uuid = ?", id)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return *omni, nil
	} else if err != nil {
		return *omni, err
	} else {
		return *omni, nil
	}
}

func (repo *Omni) PutOmni(ctx context.Context, domainUuid string, omni omni.Omni) error {
	res, err := repository.FusionSqlClient.GetDB().NewUpdate().Model(&omni).
		Where("domain_uuid = ?", domainUuid).
		Where("id = ?", omni.OmniUuid).
		WherePK().Exec(ctx)

	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 0 {
		return nil
	} else {
		return nil
	}
}
