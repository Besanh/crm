package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IUserCrm interface {
	GetUserCrms(ctx context.Context, domainUuid string, limit, offset int, filter model.UserFilter) (int, []model.UserView, error)
	GetUserCrmById(ctx context.Context, id string) (*model.UserView, error)
	GetUsersInfoOfUnit(ctx context.Context, domainUuid string, userUuid, userlevel string, unitUuids []string, level []string, isExtension bool) (*[]model.UserInfoData, error)
	PatchUserCrm(ctx context.Context, domainUuid, id, unitUuid, roleUuid string) error
}

var UserCrmRepo IUserCrm
