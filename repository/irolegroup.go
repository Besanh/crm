package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IRoleGroup interface {
	InsertRoleGroup(ctx context.Context, domainUuid string, userUuid string, roleGroup model.RoleGroup) error
	GetRoleGroup(ctx context.Context, domainUuid string, limit, offset int, filter model.RoleGroupFilter) (int, *[]model.RoleGroup, error)
	GetRoleGroupById(ctx context.Context, domainUuid, id string) (*model.RoleGroup, error)
	UpdateRoleGroupById(ctx context.Context, domainUuid string, roleGroup model.RoleGroup) error
	DeleteRoleGroupById(ctx context.Context, domainUuid, id string) error
}

var RoleGroupRepo IRoleGroup
