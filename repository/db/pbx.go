package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type Pbx struct{}

func NewPbxRepo() repository.IPbx {
	repo := &Pbx{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *Pbx) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.Pbx)(nil)); err != nil {
		panic(err)
	}
}

func (repo *Pbx) InitColumn() {

}

func (repo *Pbx) InitIndex() {

}

func (repo *Pbx) InsertPbx(ctx context.Context, pbx model.Pbx) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(&pbx)
	resp, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected != 1 {
		return errors.New("insert pbx failed")
	}
	return nil
}

func (repo *Pbx) GetPbxs(ctx context.Context, domainUuid, unitUuid string, filter model.PbxFilter) (*[]model.Pbx, error) {
	result := new([]model.Pbx)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("domain_uuid = ?", domainUuid)
	if len(unitUuid) > 0 {
		query.Where("unit_uuid = ?", unitUuid)
	}

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return result, nil
}

func (repo *Pbx) GetPbxByUnitId(ctx context.Context, domainUuid, unitUuid string) (model.Pbx, error) {
	result := new(model.Pbx)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("domain_uuid = ?", domainUuid).
		Where("unit_uuid = ?", unitUuid)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return *result, nil
	} else if err != nil {
		return *result, err
	}

	return *result, nil
}

func (repo *Pbx) GetPbxById(ctx context.Context, domainUuid, id string) (model.Pbx, error) {
	result := new(model.Pbx)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("domain_uuid = ?", domainUuid).
		Where("pbx_uuid = ?", id)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return *result, nil
	} else if err != nil {
		return *result, err
	}

	return *result, nil
}

func (repo *Pbx) PutPbxById(ctx context.Context, domainUuid string, pbx model.Pbx) error {
	query := repository.FusionSqlClient.GetDB().NewUpdate().Model(&pbx).
		WherePK()
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 0 {
		return nil
	}

	return nil
}

func (repo *Pbx) DeletePbxById(ctx context.Context, domainUuid, id string) error {
	query := repository.FusionSqlClient.GetDB().NewDelete().Model(&model.Pbx{}).
		Where("pbx_uuid = ?", id)
	if err := query.Scan(ctx); err != nil {
		return err
	} else if err == sql.ErrNoRows {
		return nil
	}
	return nil
}
