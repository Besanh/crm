package repository

import (
	"contactcenter-api/common/model/location"
	"context"
)

type ILocationWard interface {
	GetWardByEntity(ctx context.Context, districtCode string, entity string, limit, offset int) (int, *[]location.LocationWard, error)
	GetWardByCode(ctx context.Context, wardCode string) (*location.LocationWard, error)
}

var LocationWardRepo ILocationWard
