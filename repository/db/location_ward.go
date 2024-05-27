package db

import (
	"contactcenter-api/common/model/location"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/uptrace/bun"
)

type LocationWard struct{}

func NewLocationWardRepo() repository.ILocationWard {
	repo := &LocationWard{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *LocationWard) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*location.LocationWard)(nil)); err != nil {
		panic(err)
	}
}

func (repo *LocationWard) InitColumn() {

}

func (repo *LocationWard) InitIndex() {

}

/**
* entity: name, code
 */
func (repo *LocationWard) GetWardByEntity(ctx context.Context, districtCode string, entity string, limit, offset int) (int, *[]location.LocationWard, error) {
	result := new([]location.LocationWard)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Relation("LocationDistrict").Relation("LocationDistrict.LocationProvince")
	if len(districtCode) > 0 {
		query.Where("lw.district_code = ?", districtCode)
	}
	if len(entity) > 0 {
		query.WhereOr("? ILIKE ?", bun.Ident("lw.ward_name"), "%"+entity+"%").
			WhereOr("lw.ward_code = ?", entity)
	}
	query.Where("lw.status = ?", true)

	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, nil, errors.New("ward not found")
	} else if err != nil {
		return 0, nil, err
	}

	return total, result, nil
}

func (repo *LocationWard) GetWardByCode(ctx context.Context, wardUuid string) (*location.LocationWard, error) {
	result := new(location.LocationWard)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result)
	if len(wardUuid) > 0 {
		query.WhereOr("ward_uuid = ?", wardUuid)
	}
	query.Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, errors.New("ward not found")
	} else if err != nil {
		return nil, err
	}

	return result, nil
}
