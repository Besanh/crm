package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IPbx interface {
	InsertPbx(ctx context.Context, pbx model.Pbx) error
	GetPbxs(ctx context.Context, domainUuid, unitUuid string, pbx model.PbxFilter) (*[]model.Pbx, error)
	GetPbxById(ctx context.Context, domainUuid, id string) (model.Pbx, error)
	GetPbxByUnitId(ctx context.Context, domainUuid, unitUuid string) (model.Pbx, error)
	PutPbxById(ctx context.Context, domainUuid string, pbx model.Pbx) error
	DeletePbxById(ctx context.Context, domainUuid, id string) error
}

var PbxRepo IPbx
