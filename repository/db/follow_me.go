package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
)

type FollowMe struct {
}

func NewFollowMe() repository.IFollowMe {
	return &FollowMe{}
}

func (repo *FollowMe) GetFollowMeById(ctx context.Context, domainUuid, followMeUuid string) (*model.FollowMe, error) {
	followMe := new(model.FollowMe)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(followMe).
		Where("domain_uuid = ?", domainUuid).
		Where("follow_me_uuid = ?", followMeUuid).
		Limit(1)
	err := query.Scan(ctx)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return followMe, err
}

func (repo *FollowMe) GetFollowMeDestinationByFollowMeUuid(ctx context.Context, domainUuid, followMeUuid string) (*model.FollowMeDestination, error) {
	followMeDestination := new(model.FollowMeDestination)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(followMeDestination).
		Where("domain_uuid = ?", domainUuid).
		Where("follow_me_uuid = ?", followMeUuid).
		Limit(1)
	err := query.Scan(ctx)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return followMeDestination, err
}
