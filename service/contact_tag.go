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
	IContactTag interface {
		PostContactTag(ctx context.Context, domainUuid, userUuid string, contactTag *model.ContactTag) (int, any)
		GetContactTags(ctx context.Context, domainUuid string, limit, offset int, filter model.ContactTagFilter) (int, any)
		GetContactTagById(ctx context.Context, domainUuid, id string) (int, any)
		PutContactTagById(ctx context.Context, domainUuid, userUuid, id string, contactTag model.ContactTag) (int, any)
		DeleteContactTag(ctx context.Context, domainUuid, id string) (int, any)
		ExportContactTags(ctx context.Context, domainUuid, userUuid, fileType string, filter model.ContactTagFilter) (string, error)
	}
	ContactTag struct {
	}
)

func NewContactTag() IContactTag {
	return &ContactTag{}
}

func (s *ContactTag) PostContactTag(ctx context.Context, domainUuid, userUuid string, contactTag *model.ContactTag) (int, any) {
	contactTagExist, err := repository.ContactTagRepo.GetContacTagByTagName(ctx, domainUuid, contactTag.TagName)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if contactTagExist != nil {
		return response.BadRequestMsg("contact tag name already exist")
	}

	contactTag.DomainUuid = domainUuid
	contactTag.ContactTagUuid = uuid.NewString()
	contactTag.CreatedBy = userUuid
	contactTag.CreatedAt = time.Now()

	contactTagUser := []model.ContactTagUser{}
	if len(contactTag.Member) > 0 {
		for _, val := range contactTag.Member {
			contactTagUser = append(contactTagUser, model.ContactTagUser{
				DomainUuid:         domainUuid,
				ContactTagUserUuid: uuid.NewString(),
				ContactTagUuid:     contactTag.ContactTagUuid,
				EntityType:         "member",
				UserUuid:           val,
			})
		}
	}

	if len(contactTag.Staff) > 0 {
		for _, val := range contactTag.Staff {
			contactTagUser = append(contactTagUser, model.ContactTagUser{
				DomainUuid:         domainUuid,
				ContactTagUserUuid: uuid.NewString(),
				EntityType:         "staff",
				ContactTagUuid:     contactTag.ContactTagUuid,
				UserUuid:           val,
			})
		}
	}

	if err := repository.ContactTagRepo.InsertContactTagTransaction(ctx, contactTag, &contactTagUser); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Created(map[string]any{
		"id": contactTag.ContactTagUuid,
	})
}

func (s *ContactTag) GetContactTags(ctx context.Context, domainUuid string, limit, offset int, filter model.ContactTagFilter) (int, any) {
	total, contactTags, err := repository.ContactTagRepo.GetContactTags(ctx, domainUuid, limit, offset, filter)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Pagination(contactTags, total, limit, offset)
}

func (s *ContactTag) GetContactTagById(ctx context.Context, domainUuid, id string) (int, any) {
	contactTag, err := repository.ContactTagRepo.GetContactTagById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(contactTag)
}

func (s *ContactTag) PutContactTagById(ctx context.Context, domainUuid, userUuid, id string, contactTag model.ContactTag) (int, any) {
	contactTagExist, err := repository.ContactTagRepo.GetContactTagById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	contactTagInfoExist, err := repository.ContactTagRepo.GetContacTagByTagName(ctx, domainUuid, contactTag.TagName)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if contactTagInfoExist == nil {
		return response.BadRequestMsg("contact tag name not found")
	} else if contactTagExist.TagName != contactTag.TagName {
		return response.BadRequestMsg("contact tag name already exist")
	}

	contactTagExist.TagName = contactTag.TagName
	contactTagExist.LimitedFunction = contactTag.LimitedFunction
	contactTagExist.Status = contactTag.Status
	contactTagExist.Description = contactTag.Description
	contactTagExist.UpdatedBy = userUuid
	contactTagExist.UpdatedAt = time.Now()

	contactTagUser := []model.ContactTagUser{}
	if len(contactTag.Member) > 0 {
		for _, val := range contactTag.Member {
			contactTagUser = append(contactTagUser, model.ContactTagUser{
				DomainUuid:         domainUuid,
				ContactTagUserUuid: uuid.NewString(),
				ContactTagUuid:     contactTag.ContactTagUuid,
				EntityType:         "member",
				UserUuid:           val,
			})
		}
	}

	if len(contactTag.Staff) > 0 {
		for _, val := range contactTag.Staff {
			contactTagUser = append(contactTagUser, model.ContactTagUser{
				DomainUuid:         domainUuid,
				ContactTagUserUuid: uuid.NewString(),
				EntityType:         "staff",
				ContactTagUuid:     contactTag.ContactTagUuid,
				UserUuid:           val,
			})
		}
	}

	if err := repository.ContactTagRepo.PutContactTagTransaction(ctx, contactTagExist, &contactTagUser); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"id": id,
	})
}

func (s *ContactTag) DeleteContactTag(ctx context.Context, domainUuid, id string) (int, any) {
	contactTagExist, err := repository.ContactTagRepo.GetContactTagById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if contactTagExist == nil {
		return response.NotFoundMsg("contact tag not found")
	}

	if err := repository.ContactTagRepo.DeleteContactTag(ctx, domainUuid, contactTagExist); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"id": id,
	})
}
