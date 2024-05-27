package db

import (
	"contactcenter-api/common/model/location"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type LocationDistrict struct {
}

func NewLocationDistrictRepo() repository.ILocationDistrict {
	repo := &LocationDistrict{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *LocationDistrict) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*location.LocationDistrict)(nil)); err != nil {
		panic(err)
	}
}

func (repo *LocationDistrict) InitColumn() {

}

func (repo *LocationDistrict) InitIndex() {
}

func (repo *LocationDistrict) GetDistrictByEntity(ctx context.Context, provinceCode string, entity string, limit, offset int) (int, *[]location.LocationDistrict, error) {
	result := new([]location.LocationDistrict)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result)
	if len(provinceCode) > 0 {
		query.Where("province_code = ?", provinceCode)
	}
	if len(entity) > 0 {
		query.WhereOr("district_name = ?", entity).
			WhereOr("district_code = ?", entity)
	}

	query.Where("status = ?", true)
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, nil, errors.New("district not found")
	} else if err != nil {
		return 0, nil, err
	}

	return total, result, nil
}

func (repo *LocationDistrict) GetDistrictByCode(ctx context.Context, districtCode string) (*location.LocationDistrict, error) {
	result := new(location.LocationDistrict)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result)
	if len(districtCode) > 0 {
		query.Where("district_code = ?", districtCode)
	}
	query.Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, errors.New("district not found")
	} else if err != nil {
		return nil, err
	}

	return result, nil
}
