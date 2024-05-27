package repository

import (
	"contactcenter-api/common/model/location"
	"context"
)

type ILocationDistrict interface {
	GetDistrictByEntity(ctx context.Context, provinceUuid string, entity string, limit, offset int) (int, *[]location.LocationDistrict, error)
	GetDistrictByCode(ctx context.Context, districtCode string) (*location.LocationDistrict, error)
}

var LocationDistrictRepo ILocationDistrict
