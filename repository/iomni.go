package repository

import (
	"contactcenter-api/common/model"
	"contactcenter-api/common/model/omni"
	"context"
)

type IOmni interface {
	GetOmnis(ctx context.Context, domainUuid string, limit, offset int, filter model.OmniFilter) ([]omni.Omni, int, error)
	GetOmniById(ctx context.Context, domainUuid, id string) (omni.Omni, error)
	PutOmni(ctx context.Context, domainUuid string, omni omni.Omni) error
}

var OmniRepo IOmni
