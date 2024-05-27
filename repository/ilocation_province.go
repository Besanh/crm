package repository

import (
	"contactcenter-api/common/model/location"
	"context"
)

type ILocationProvince interface {
	GetProvinceByNameOrCode(ctx context.Context, entity string, limit, offset int) (int, *[]location.LocationProvince, error)
	GetProvinceByCode(ctx context.Context, provinceCode string) (*location.LocationProvince, error)
}

var LocationProvinceRepo ILocationProvince
