package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IFollowMe interface {
	GetFollowMeById(ctx context.Context, domainUuid, followMeUuid string) (*model.FollowMe, error)
	GetFollowMeDestinationByFollowMeUuid(ctx context.Context, domainUuid, followMeUuid string) (*model.FollowMeDestination, error)
}

var FollowMeRepo IFollowMe
