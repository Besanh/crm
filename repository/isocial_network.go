package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type (
	ISocialNetwork interface {
		GetSocialNetworkByUserId(ctx context.Context, domainUuid, userUuid, statusType string, isActiveAll bool) (*model.SocialNetwork, error)
	}
)

var SocialNetworkRepo ISocialNetwork
