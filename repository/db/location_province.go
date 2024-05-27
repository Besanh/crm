package db

import (
	"contactcenter-api/common/model/location"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type LocationProvince struct {
}

func NewLocationProvinceRepo() repository.ILocationProvince {
	repo := &LocationProvince{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *LocationProvince) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*location.LocationProvince)(nil)); err != nil {
		panic(err)
	}
}

func (repo *LocationProvince) InitColumn() {

}

func (repo *LocationProvince) InitIndex() {
}

func (repo *LocationProvince) GetProvinceByNameOrCode(ctx context.Context, entity string, limit, offset int) (int, *[]location.LocationProvince, error) {
	result := new([]location.LocationProvince)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result)
	if len(entity) > 0 {
		query.WhereOr("province_name = ?", entity).
			WhereOr("province_code = ?", entity)
	}

	query.Where("status = ?", true)

	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, nil, errors.New("province not found")
	} else if err != nil {
		return 0, nil, err
	}

	return total, result, nil
}

func (repo *LocationProvince) GetProvinceByCode(ctx context.Context, provinceCode string) (*location.LocationProvince, error) {
	result := new(location.LocationProvince)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result)
	if len(provinceCode) > 0 {
		query.Where("province_code = ?", provinceCode)
	}
	query.Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, errors.New("province not found")
	} else if err != nil {
		return nil, err
	}

	return result, nil
}
