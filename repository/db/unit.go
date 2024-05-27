package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type UnitRepo struct {
}

func NewUnit() repository.IUnit {
	repo := &UnitRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()
	return repo
}

func (repo *UnitRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	repository.FusionSqlClient.GetDB().RegisterModel((*model.Unit)(nil))
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.Unit)(nil)); err != nil {
		panic(err)
	}
}

func (repo *UnitRepo) InitColumn() {
}

func (repo *UnitRepo) InitIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Unit)(nil)).IfNotExists().Index("idx_unit_domain_uuid").Column("domain_uuid").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Unit)(nil)).IfNotExists().Index("idx_unit_parent_unit_uuid").Column("parent_unit_uuid").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Unit)(nil)).IfNotExists().Index("idx_unit_uuid").Column("unit_uuid").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Unit)(nil)).IfNotExists().Index("idx_unit_name").Column("unit_name").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Unit)(nil)).IfNotExists().Index("idx_unit_code").Column("unit_code").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Unit)(nil)).IfNotExists().Index("idx_unit_leader").Column("unit_leader").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Unit)(nil)).IfNotExists().Index("idx_unit_basis").Column("unit_basis").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Unit)(nil)).IfNotExists().Index("idx_status").Column("status").Exec(ctx); err != nil {
		panic(err)
	}
}

func (repo *UnitRepo) InsertUnit(ctx context.Context, unit model.Unit) error {
	resp, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&unit).Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("insert unit failed")
	}
	return nil
}

func (repo *UnitRepo) GetUnits(ctx context.Context, domainUuid string, limit, offset int, filter model.UnitFilter) (int, *[]model.UnitInfo, error) {
	result := new([]model.UnitInfo)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Relation("Users").
		Where("unit.domain_uuid = ?", domainUuid)
	if len(filter.ParentUnitUuid) > 0 {
		query.Where("parent_unit_uuid = ?", filter.ParentUnitUuid)
	}
	if filter.IsParent {
		query.Where("parent_unit_uuid = ''")
	}
	if len(filter.UnitUuid) > 0 {
		query.Where("unit_uuid = ?", filter.UnitUuid)
	}
	if len(filter.UnitName) > 0 {
		query.Where("unit_name = ?", filter.UnitName)
	}
	if len(filter.UnitCode) > 0 {
		query.Where("unit_code = ?", filter.UnitCode)
	}
	if filter.UnitBasis.Valid {
		query.Where("unit_basis = ?", filter.UnitBasis.Bool)
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status.Bool)
	}
	if len(filter.Level) > 0 {
		query.Where("level = ?", filter.Level)
	}

	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	query.Order("unit.created_at desc")

	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return 0, nil, err
	} else if err == sql.ErrNoRows {
		return 0, nil, nil
	}

	return total, result, nil
}

func (repo *UnitRepo) GetUnitById(ctx context.Context, domainUuid, id string) (*model.UnitInfo, error) {
	result := new(model.UnitInfo)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Relation("Users").
		Where("unit.domain_uuid = ?", domainUuid).
		Where("unit_uuid = ?", id)

	err := query.Scan(ctx)
	if err != nil {
		return result, err
	} else if err == sql.ErrNoRows {
		return result, errors.New("unit not found")
	}

	return result, nil
}

func (repo *UnitRepo) GetUnitRelationById(ctx context.Context, domainUuid, id string) (*model.UnitInfo, error) {
	result := new(model.UnitInfo)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Relation("Users").
		Where("unit.domain_uuid = ?", domainUuid).
		Where("unit_uuid = ?", id)
	err := query.Scan(ctx)
	if err != nil {
		return result, err
	} else if err == sql.ErrNoRows {
		return result, errors.New("unit not found")
	}

	return result, nil
}

func (repo *UnitRepo) PutUnit(ctx context.Context, domainUuid string, unit model.UnitInfo) error {
	query := repository.FusionSqlClient.GetDB().NewUpdate().Model(&unit).
		Where("domain_uuid = ?", domainUuid).
		Where("unit_uuid = ?", unit.UnitUuid)
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("update unit failed")
	}
	return nil
}

func (repo *UnitRepo) DeleteUnitById(ctx context.Context, domainUuid, id string) error {
	query := repository.FusionSqlClient.GetDB().NewDelete().Model(&model.Unit{}).
		Where("domain_uuid = ?", domainUuid).
		Where("unit_uuid = ?", id)
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("delete unit failed")
	}
	return nil
}
