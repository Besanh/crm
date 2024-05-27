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
	IContactCareer interface {
		PostContactCareer(ctx context.Context, domainUuid, userUuid string, contactCareer *model.ContactCareer) (int, any)
		GetContactCareers(ctx context.Context, domainUuid string, limit, offset int, filter model.ContactCareerFilter) (int, any)
		GetContactCareerById(ctx context.Context, domainUuid, id string) (int, any)
		PutContactCareerById(ctx context.Context, domainUuid, userUuid, id string, contactCareer model.ContactCareer) (int, any)
		DeleteContactCareer(ctx context.Context, domainUuid, id string) (int, any)
		ExportContactCareers(ctx context.Context, domainUuid, userUuid, fileType string, filter model.ContactCareerFilter) (string, error)
	}
	ContactCareer struct {
	}
)

func NewContactCareer() IContactCareer {
	return &ContactCareer{}
}

func (s *ContactCareer) PostContactCareer(ctx context.Context, domainUuid, userUuid string, contactCareer *model.ContactCareer) (int, any) {
	contactCareerExist, err := repository.ContactCareerRepo.GetContacCareerByCareerName(ctx, domainUuid, contactCareer.CareerName)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if contactCareerExist != nil {
		return response.BadRequestMsg("contact career name already exist")
	}

	contactCareer.DomainUuid = domainUuid
	contactCareer.ContactCareerUuid = uuid.NewString()
	contactCareer.CreatedBy = userUuid
	contactCareer.CreatedAt = time.Now()

	contactCareerUsers := []model.ContactCareerUser{}
	if len(contactCareer.Career) > 0 {
		for _, val := range contactCareer.Career {
			contactCareerUsers = append(contactCareerUsers, model.ContactCareerUser{
				DomainUuid:            domainUuid,
				ContactCareerUserUuid: uuid.NewString(),
				ContactCareerUuid:     contactCareer.ContactCareerUuid,
				EntityType:            "career",
				UserUuid:              val,
			})
		}
	}

	if err := repository.ContactCareerRepo.InsertContactCareerTransaction(ctx, contactCareer, &contactCareerUsers); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Created(map[string]any{
		"id": contactCareer.ContactCareerUuid,
	})
}

func (s *ContactCareer) GetContactCareers(ctx context.Context, domainUuid string, limit, offset int, filter model.ContactCareerFilter) (int, any) {
	total, contactCareers, err := repository.ContactCareerRepo.GetContactCareers(ctx, domainUuid, limit, offset, filter)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Pagination(contactCareers, total, limit, offset)
}

func (s *ContactCareer) GetContactCareerById(ctx context.Context, domainUuid, id string) (int, any) {
	contactCareer, err := repository.ContactCareerRepo.GetContactCareerById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if contactCareer == nil {
		return response.NotFoundMsg("contact career not found")
	}

	return response.OK(contactCareer)
}

func (s *ContactCareer) PutContactCareerById(ctx context.Context, domainUuid, userUuid, id string, contactCareer model.ContactCareer) (int, any) {
	contactCareerExist, err := repository.ContactCareerRepo.GetContactCareerById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if contactCareerExist == nil {
		return response.NotFoundMsg("contact career not found")
	}

	// contactCareerInfoExist, err := repository.ContactCareerRepo.GetContacCareerByCareerName(ctx, domainUuid, contactCareer.CareerName)
	// if err != nil {
	// 	log.Error(err)
	// 	return response.ServiceUnavailableMsg(err.Error())
	// } else if contactCareerInfoExist == nil {
	// 	return response.BadRequestMsg("contact career name not found")
	// } else if contactCareerExist.CareerName != contactCareer.CareerName {
	// 	return response.BadRequestMsg("contact career name already exist")
	// }

	contactCareerExist.CareerName = contactCareer.CareerName
	contactCareerExist.Status = contactCareer.Status
	contactCareerExist.Description = contactCareer.Description
	contactCareerExist.UpdatedBy = userUuid
	contactCareerExist.UpdatedAt = time.Now()

	contactCareerUsers := []model.ContactCareerUser{}
	if len(contactCareer.Career) > 0 {
		for _, val := range contactCareer.Career {
			contactCareerUsers = append(contactCareerUsers, model.ContactCareerUser{
				DomainUuid:            domainUuid,
				ContactCareerUserUuid: uuid.NewString(),
				ContactCareerUuid:     contactCareer.ContactCareerUuid,
				EntityType:            "career",
				UserUuid:              val,
			})
		}
	}

	if err := repository.ContactCareerRepo.PutContactCareer(ctx, contactCareerExist, &contactCareerUsers); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"id": id,
	})
}

func (s *ContactCareer) DeleteContactCareer(ctx context.Context, domainUuid, id string) (int, any) {
	contactCareerExist, err := repository.ContactCareerRepo.GetContactCareerById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if contactCareerExist == nil {
		return response.NotFoundMsg("contact career not found")
	}

	if err := repository.ContactCareerRepo.DeleteContactCareer(ctx, domainUuid, contactCareerExist); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"id": id,
	})
}
