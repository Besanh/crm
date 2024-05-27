package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type (
	IContactGroup interface {
		PostContactGroup(ctx context.Context, domainUuid, userUuid string, contactGroup *model.ContactGroup) (int, any)
		GetContactGroups(ctx context.Context, domainUuid string, limit, offset int, filter model.ContactGroupFilter) (int, any)
		GetContactGroupById(ctx context.Context, domainUuid, id string) (int, any)
		PutContactGroupById(ctx context.Context, domainUuid, userUuid, id string, contactGroup model.ContactGroup) (int, any)
		DeleteContactGroup(ctx context.Context, domainUuid, id string) (int, any)
		ExportContactGroups(ctx context.Context, domainUuid, userUuid, fileType string, filter model.ContactGroupFilter) (string, error)
	}
	ContactGroup struct {
	}
)

func NewContactGroup() IContactGroup {
	return &ContactGroup{}
}

func (s *ContactGroup) PostContactGroup(ctx context.Context, domainUuid, userUuid string, contactGroup *model.ContactGroup) (int, any) {
	contactGroupExist, err := repository.ContactTagRepo.GetContacTagByTagName(ctx, domainUuid, contactGroup.GroupName)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if contactGroupExist != nil {
		return response.BadRequestMsg("contact group name already exist")
	}

	contactGroup.DomainUuid = domainUuid
	contactGroup.ContactGroupUuid = uuid.NewString()
	contactGroup.CreatedBy = userUuid
	contactGroup.CreatedAt = time.Now()

	contactGroupUsers := []model.ContactGroupUser{}
	if len(contactGroup.Member) > 0 {
		for _, val := range contactGroup.Member {
			contactGroupUsers = append(contactGroupUsers, model.ContactGroupUser{
				DomainUuid:           domainUuid,
				ContactGroupUserUuid: uuid.NewString(),
				ContactGroupUuid:     contactGroup.ContactGroupUuid,
				EntityType:           "member",
				UserUuid:             val,
			})
		}
	}

	if len(contactGroup.Staff) > 0 {
		for _, val := range contactGroup.Staff {
			contactGroupUsers = append(contactGroupUsers, model.ContactGroupUser{
				DomainUuid:           domainUuid,
				ContactGroupUserUuid: uuid.NewString(),
				EntityType:           "staff",
				ContactGroupUuid:     contactGroup.ContactGroupUuid,
				UserUuid:             val,
			})
		}
	}

	if err := repository.ContactGroupRepo.InsertContactGroupTransaction(ctx, contactGroup, &contactGroupUsers); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Created(map[string]any{
		"id": contactGroup.ContactGroupUuid,
	})
}

func (s *ContactGroup) GetContactGroups(ctx context.Context, domainUuid string, limit, offset int, filter model.ContactGroupFilter) (int, any) {
	total, contactGroups, err := repository.ContactGroupRepo.GetContactGroups(ctx, domainUuid, limit, offset, filter)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Pagination(contactGroups, total, limit, offset)
}

func (s *ContactGroup) GetContactGroupById(ctx context.Context, domainUuid, id string) (int, any) {
	contactGroup, err := repository.ContactGroupRepo.GetContactGroupById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if contactGroup == nil {
		return response.NotFoundMsg("contact group not found")
	}

	return response.OK(contactGroup)
}

func (s *ContactGroup) PutContactGroupById(ctx context.Context, domainUuid, userUuid, id string, contactGroup model.ContactGroup) (int, any) {
	contactGroupExist, err := repository.ContactGroupRepo.GetContactGroupById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if contactGroupExist == nil {
		return response.NotFoundMsg("contact group not found")
	}

	// contactGroupInfoExist, err := repository.ContactGroupRepo.GetContactGroupByGroupName(ctx, domainUuid, contactGroup.GroupName)
	// if err != nil {
	// 	log.Error(err)
	// 	return response.ServiceUnavailableMsg(err.Error())
	// } else if contactGroupInfoExist == nil {
	// 	return response.BadRequestMsg("contact group name not found")
	// } else if contactGroupExist.GroupName != contactGroup.GroupName {
	// 	return response.BadRequestMsg("contact group name already exist")
	// }

	contactGroupExist.GroupName = contactGroup.GroupName
	contactGroupExist.Status = contactGroup.Status
	contactGroupExist.Description = contactGroup.Description
	contactGroupExist.UpdatedBy = userUuid
	contactGroupExist.UpdatedAt = time.Now()

	contactGroupUsers := []model.ContactGroupUser{}
	if len(contactGroup.Member) > 0 {
		for _, val := range contactGroup.Member {
			contactGroupUsers = append(contactGroupUsers, model.ContactGroupUser{
				DomainUuid:           domainUuid,
				ContactGroupUserUuid: uuid.NewString(),
				ContactGroupUuid:     contactGroup.ContactGroupUuid,
				EntityType:           "member",
				UserUuid:             val,
			})
		}
	}

	if len(contactGroup.Staff) > 0 {
		for _, val := range contactGroup.Staff {
			contactGroupUsers = append(contactGroupUsers, model.ContactGroupUser{
				DomainUuid:           domainUuid,
				ContactGroupUserUuid: uuid.NewString(),
				EntityType:           "staff",
				ContactGroupUuid:     contactGroup.ContactGroupUuid,
				UserUuid:             val,
			})
		}
	}

	if err := repository.ContactGroupRepo.PutContactGroup(ctx, contactGroupExist, &contactGroupUsers); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"id": id,
	})
}

func (s *ContactGroup) DeleteContactGroup(ctx context.Context, domainUuid, id string) (int, any) {
	contactGroupExist, err := repository.ContactGroupRepo.GetContactGroupById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if contactGroupExist == nil {
		return response.NotFoundMsg("contact group not found")
	}

	if err := repository.ContactGroupRepo.DeleteContactGroup(ctx, domainUuid, contactGroupExist); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"id": id,
	})
}
