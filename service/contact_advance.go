package service

import (
	"contactcenter-api/common/constants"
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	"contactcenter-api/repository"
	"context"
	"errors"
	"os"
	"time"
)

func (s *Contact) GetContactNotesOfId(ctx context.Context, domainUuid string, userUuid string, contactUuid string, limit, offset int) (int, any) {
	user, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	_ = user
	result, total, err := repository.ContactRepo.GetContactNotes(ctx, domainUuid, contactUuid, limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.Pagination(*result, total, limit, offset)
}

func (s *Contact) PostContactNote(ctx context.Context, domainUuid, userUuid string, contactNotePost model.ContactNotePost) (int, any) {
	user, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	userCreatedUuid := user.UserUuid
	if util.InArray(user.Level, []string{constants.ADMIN, constants.SUPERADMIN, constants.MANAGER, constants.LEADER}) && len(contactNotePost.NoteBy) > 0 {
		userTmp, err := repository.UserRepo.GetUserByIdOrUsername(ctx, domainUuid, contactNotePost.NoteBy)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if userTmp == nil {
			return response.BadRequestMsg("user is not existed")
		}
		userCreatedUuid = userTmp.UserUuid
	}
	contact := new(model.Contact)
	if len(contactNotePost.ContactUuid) > 0 {
		contact, err = repository.ContactRepo.GetContactById(ctx, domainUuid, contactNotePost.ContactUuid)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if contact == nil {
			return response.BadRequestMsg("contact is not existed")
		}
	} else {
		contact, err = repository.ContactRepo.GetContactByPhoneNumber(ctx, domainUuid, contactNotePost.Phone)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if contact == nil {
			return response.BadRequestMsg("contact is not existed")
		}
	}
	contactNote := model.ContactNote{
		DomainUuid:  domainUuid,
		ContactUuid: contact.ContactUuid,
		Content:     contactNotePost.Content,
		Type:        contactNotePost.Type,
		UserUuid:    userCreatedUuid,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := repository.ContactRepo.InsertContactNote(ctx, &contactNote); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(map[string]any{
		"contact_uuid": contact.ContactUuid,
		"message":      "success",
	})
}

func (s *Contact) PatchContactAvatar(ctx context.Context, domainUuid string, userUuid string, contactUuid string, data string) (int, any) {
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
	dir := util.UPLOAD_DIR + "avatar/"
	fileName := "avatar_" + contact.ContactUuid
	if fileName, err = util.DecodeAndSaveImageBase64(data, dir, fileName); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err)
	}
	//contact.AvatarUrl = fmt.Sprintf("%s/%s/%s", API_HOST, "v1/crm/contact/avatar_file", fileName)
	if err := repository.ContactRepo.UpdateContact(ctx, contact); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(map[string]any{
		"contact_uuid": contact.ContactUuid,
	})
}

func (s *Contact) GetAvatar(ctx context.Context, fileName string) (string, error) {
	avatar := ""
	if _, err := os.Stat(util.UPLOAD_DIR + "avatar/" + fileName); err != nil {
		return "", errors.New("file " + fileName + " is not exist")

	}
	avatar = util.UPLOAD_DIR + "avatar/" + fileName
	return avatar, nil
}
