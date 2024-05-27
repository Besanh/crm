package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/repository"
	"contactcenter-api/service/common"
	"context"
	"time"

	"github.com/google/uuid"
)

type (
	IRoleGroup interface {
		GetRoleGroup(ctx context.Context, domainUuid, userUuid string, limit, offset int, filter model.RoleGroupFilter) (int, any)
		PostRoleGroup(ctx context.Context, domainUuid, userUuid string, roleGroup model.RoleGroup) (int, any)
		GetRoleGroupById(ctx context.Context, domainUuid, userUuid, id string) (int, any)
		PutRoleGroup(ctx context.Context, domainUuid, userUuid, id string, roleGroup model.RoleGroup) (int, any)
		DeleteRoleGroupById(ctx context.Context, domainUuid, userUuid, id string) (int, any)
		ExportRoleGroups(ctx context.Context, domainUuid, userUuid string, filter model.RoleGroupFilter) (string, error)
	}
	RoleGroup struct{}
)

func NewRoleGroup() IRoleGroup {
	return &RoleGroup{}
}

func (s *RoleGroup) GetRoleGroup(ctx context.Context, domainUuid, userUuid string, limit, offset int, filter model.RoleGroupFilter) (int, any) {
	total, result, err := repository.RoleGroupRepo.GetRoleGroup(ctx, domainUuid, limit, offset, filter)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Pagination(result, total, limit, offset)
}

func (s *RoleGroup) PostRoleGroup(ctx context.Context, domainUuid, userUuid string, roleGroup model.RoleGroup) (int, any) {
	roleGroup.DomainUuid = domainUuid
	roleGroup.RoleGroupUuid = uuid.New().String()
	roleGroup.CreatedBy = userUuid
	roleGroup.CreatedAt = time.Now()

	// Permission module

	roleGroup.PermissionMain = common.HandleCreatePermissionMain()
	roleGroup.PermissionAdvance = common.HandleCreatePermissionAdvance()

	if err := repository.RoleGroupRepo.InsertRoleGroup(ctx, domainUuid, userUuid, roleGroup); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Created(map[string]any{
		"id": roleGroup.RoleGroupUuid,
	})
}

func (s *RoleGroup) GetRoleGroupById(ctx context.Context, domainUuid, userUuid, id string) (int, any) {
	roleGroup, err := repository.RoleGroupRepo.GetRoleGroupById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if roleGroup == nil {
		return response.NotFoundMsg("role group is not found")
	}

	return response.OK(map[string]any{
		"data": roleGroup,
	})
}

func (s *RoleGroup) PutRoleGroup(ctx context.Context, domainUuid, userUuid, id string, roleGroup model.RoleGroup) (int, any) {
	roleGroupExist, err := repository.RoleGroupRepo.GetRoleGroupById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if roleGroupExist == nil {
		return response.NotFoundMsg("role group is not found")
	}
	roleGroupExist.RoleGroupName = roleGroup.RoleGroupName
	roleGroupExist.Status = roleGroup.Status
	roleGroup.UpdatedBy = userUuid
	roleGroup.UpdatedAt = time.Now()

	// Permission module
	roleGroupExist.PermissionMain = roleGroup.PermissionMain
	roleGroupExist.PermissionAdvance = roleGroup.PermissionAdvance

	if err := repository.RoleGroupRepo.UpdateRoleGroupById(ctx, domainUuid, roleGroup); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"id": id,
	})
}

func (s *RoleGroup) DeleteRoleGroupById(ctx context.Context, domainUuid, userUuid, id string) (int, any) {
	roleGroupExist, err := repository.RoleGroupRepo.GetRoleGroupById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if roleGroupExist == nil {
		return response.NotFoundMsg("role group is not found")
	}

	if err := repository.RoleGroupRepo.DeleteRoleGroupById(ctx, domainUuid, id); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"id": id,
	})
}
