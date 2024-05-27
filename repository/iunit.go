package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IUnit interface {
	InsertUnit(ctx context.Context, unit model.Unit) error
	GetUnits(ctx context.Context, domainUuid string, limit, offset int, filter model.UnitFilter) (int, *[]model.UnitInfo, error)
	GetUnitById(ctx context.Context, domainUuid, id string) (*model.UnitInfo, error)
	GetUnitRelationById(ctx context.Context, domainUuid, id string) (*model.UnitInfo, error)
	PutUnit(ctx context.Context, domainUuid string, unit model.UnitInfo) error
	DeleteUnitById(ctx context.Context, domainUuid, id string) error
}

var UnitRepo IUnit
