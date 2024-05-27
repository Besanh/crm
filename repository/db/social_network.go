package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
)

type SocialNetwork struct{}

func NewSocialNetwork() repository.ISocialNetwork {
	return &SocialNetwork{}
}

func (repo *SocialNetwork) GetSocialNetworkByUserId(ctx context.Context, domainUuid, userUuid, statusType string, isActiveAll bool) (*model.SocialNetwork, error) {
	socialNetwork := new(model.SocialNetwork)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(socialNetwork).
		Join("INNER JOIN v_users u ON u.user_uuid = sn.user_uuid").
		Where("sn.domain_uuid = ?", domainUuid).
		Where("sn.user_uuid = ?", userUuid)
	if isActiveAll {
		query.Where("(sn.status_chat = true AND sn.status_chat = true)")
	} else {
		if statusType == "chat" {
			query = query.Where("sn.status_chat = true")
		} else if statusType == "email" {
			query = query.Where("sn.status_email = true")
		}
	}
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, errors.New("social network not found")
	} else if err != nil {
		return nil, err
	} else {
		return socialNetwork, nil
	}
}
