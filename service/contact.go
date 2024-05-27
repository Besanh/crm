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
	IContact interface {
		GetContacts(ctx context.Context, domainUuid, userUuid string, filter model.ContactFilter, limit, offset int) (int, any)
		GetContactById(ctx context.Context, domainUuid, userUuid, contactUuid string) (int, any)
		GetContactInfo(ctx context.Context, domainUuid, userUuid string, filter model.ContactFilter) (int, any)
		PostContact(ctx context.Context, domainUuid, userUuid string, contactPost model.ContactPost) (int, any)
		PostFileImportContacts(ctx context.Context, domainUuid, userUuid, filePath, extension string, isUpdateContact bool) (int, any)
		PutContact(ctx context.Context, domainUuid, userUuid, contactUuid string, contactPost model.ContactPost) (int, any)
		DeleteContact(ctx context.Context, domainUuid, userUuid, contactUuid string) (int, any)
		GetContactNotesOfId(ctx context.Context, domainUuid, userUuid, contactUuid string, limit, offset int) (int, any)
		PostContactNote(ctx context.Context, domainUuid, userUuid string, contactNotePost model.ContactNotePost) (int, any)
		PatchContactAvatar(ctx context.Context, domainUuid, userUuid, contactUuid, data string) (int, any)
		GetAvatar(ctx context.Context, fileName string) (string, error)
		PatchContactWithCallId(ctx context.Context, domainUuid, userUuid, contactUuid, callId string) (int, any)
	}
	Contact struct {
	}
)

func NewContact() IContact {
	return &Contact{}
}

func (s *Contact) GetContacts(ctx context.Context, domainUuid, userUuid string, filter model.ContactFilter, limit, offset int) (int, any) {
	total, data, err := repository.ContactRepo.GetContactsInfo(ctx, domainUuid, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	result := make([]model.ContactInfo, 0)
	for i := 0; i < len(*data); i++ {
		result = append(result, common.ParseContactInfo(&(*data)[i]))
	}
	return response.Pagination(result, total, limit, offset)
}

func (s *Contact) GetContactById(ctx context.Context, domainUuid, userUuid, contactUuid string) (int, any) {
	_, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	contactInfo, err := repository.ContactRepo.GetContactInfoById(ctx, domainUuid, contactUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if contactInfo == nil {
		return response.NotFoundMsg("contact is not existed")
	}
	contact := common.ParseContactInfo(contactInfo)

	return response.OK(contact)
}

func (s *Contact) PostContact(ctx context.Context, domainUuid, userUuid string, contactPost model.ContactPost) (int, any) {
	unitUuid := ""
	userExist, err := repository.UserCrmRepo.GetUserCrmById(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(userExist.UserUuid) < 1 {
		return response.ServiceUnavailableMsg("user is not existed")
	}
	unitUuid = userExist.UnitUuid

	contact := model.Contact{
		ContactUuid: uuid.New().String(),
		DomainUuid:  domainUuid,
		UnitUuid:    unitUuid,
		ContactType: contactPost.ContactType,
		ContactName: contactPost.ContactName,
		Status:      true,
		SourceName:  contactPost.SourceName,
		SourceUuid:  contactPost.SourceUuid,
		CreatedBy:   userUuid,
		CreatedAt:   time.Now(),
	}
	if err := repository.ContactRepo.InsertContact(ctx, &contact); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.Created(map[string]any{
		"id": contact.ContactUuid,
	})
}

func (s *Contact) PutContact(ctx context.Context, domainUuid, userUuid, contactUuid string, contactPut model.ContactPost) (int, any) {
	_, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	contact, err := repository.ContactRepo.GetContactById(ctx, domainUuid, contactUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if contact == nil {
		return response.BadRequestMsg("contact is not existed")
	}
	contact.UpdatedAt = time.Now()
	contact.ContactName = contactPut.ContactName
	contact.ContactType = contactPut.ContactType
	contact.UpdatedBy = userUuid
	contact.Status = contactPut.Status
	if err := repository.ContactRepo.UpdateContact(ctx, contact); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(map[string]any{
		"id": contact.ContactUuid,
	})
}

func (s *Contact) DeleteContact(ctx context.Context, domainUuid string, userUuid string, contactUuid string) (int, any) {
	_, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	contact, err := repository.ContactRepo.GetContactById(ctx, domainUuid, contactUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if contact == nil {
		return response.BadRequestMsg("contact is not existed")
	}
	if err := repository.ContactRepo.DeleteContactTransaction(ctx, contact.ContactUuid); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(map[string]any{
		"id": contact.ContactUuid,
	})
}

func (s *Contact) GetContactInfo(ctx context.Context, domainUuid, userUuid string, filter model.ContactFilter) (int, any) {
	contact, err := repository.ContactRepo.GetContactInfo(ctx, domainUuid, filter)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"data": contact,
	})
}

func (s *Contact) PatchContactWithCallId(ctx context.Context, domainUuid, userUuid, contactUuid, callId string) (int, any) {
	contact, err := repository.ContactRepo.GetContactById(ctx, domainUuid, contactUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(contact.ContactUuid) < 1 {
		return response.BadRequestMsg("contact is not existed")
	}

	contact.UpdatedBy = userUuid
	contact.UpdatedAt = time.Now()

	if err := repository.ContactRepo.UpdateContact(ctx, contact); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"contact_uuid": contactUuid,
	})
}
