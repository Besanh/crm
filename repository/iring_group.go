package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IRingGroup interface {
	GetRingGroupById(ctx context.Context, domainUuid, ringGroupUuid string) (*model.RingGroup, error)
	GetRingGroupByExtension(ctx context.Context, domainUuid, ringGroupExtension string) (*model.RingGroup, error)
	GetRingGroupDestinationByRingGroupUuid(ctx context.Context, domainUuid, ringGroupUuid string) (*model.RingGroupDestination, error)
	GetRingGroupDestinationOfExtension(ctx context.Context, domainUuid, ringGroupUuid, extension string) (*model.RingGroupDestination, error)
}

var RingGroupRepo IRingGroup
