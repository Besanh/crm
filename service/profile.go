package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	"contactcenter-api/repository"
	"contactcenter-api/service/common"
	"context"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
)

const (
	PRIMARY_PROFILE_TYPE  = "primary"
	RELATION_PROFILE_TYPE = "relation"
)

type (
	IProfile interface {
		GetProfileInfoById(ctx context.Context, domainUuid, userUuid, profileUuid string) (int, any)
		GetProfilesInfo(ctx context.Context, domainUuid, userUuid string, filter model.ProfileFilter, limit, offset int) (int, any)
		PostProfile(ctx context.Context, domainUuid, userUuid string, profilePost model.ProfilePost) (int, any)
		PutProfile(ctx context.Context, domainUuid, userUuid, profileUuid string, profilePost model.ProfilePost) (int, any)
		DeleteProfile(ctx context.Context, domainUuid, userUuid, profileUuid string) (int, any)
		DeleteProfileTransaction(ctx context.Context, domainUuid, userUuid string, filter model.ProfileFilter) (int, any)
		PatchProfileAvatar(ctx context.Context, domainUuid, userUuid, profileUuid string, fileName string) (int, any)
		GetAvatar(ctx context.Context, fileName string) (string, error)
		PatchStatusProfile(ctx context.Context, domainUuid, userUuid, profileId string, status bool) (int, any)
		PatchSocialMappingContactProfile(ctx context.Context, domainUuid, userUuid, profileId string, data model.SocialMappingContact) (int, any)
		GetManageProfiles(ctx context.Context, domainUuid, userUuid string, filter model.ProfileFilter, limit, offset int) (int, any)
		PostConvertLeadToProfile(ctx context.Context, domainUuid, userUuid string, profilePost model.ProfilePost) (int, any)
	}
	Profile struct {
	}
)

func NewProfile() IProfile {
	return &Profile{}
}

func (p *Profile) GetProfileInfoById(ctx context.Context, domainUuid, userUuid, profileUuid string) (int, any) {
	_, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	profileInfo, err := repository.ProfileRepo.GetProfileInfoById(ctx, domainUuid, profileUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if profileInfo == nil {
		return response.ServiceUnavailableMsg(errors.New("profile empty"))
	}
	profile := common.ParseProfileInfo(profileInfo)
	return response.OK(profile)
}

func (p *Profile) GetProfilesInfo(ctx context.Context, domainUuid, userUuid string, filter model.ProfileFilter, limit, offset int) (int, any) {
	total, data, err := repository.ProfileRepo.GetProfilesInfo(ctx, domainUuid, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	result := make([]model.ProfileInfo, 0)
	for i := 0; i < len(*data); i++ {
		result = append(result, common.ParseProfileInfo(&(*data)[i]))
	}
	return response.Pagination(result, total, limit, offset)
}

func (p *Profile) PostProfile(ctx context.Context, domainUuid, userUuid string, profilePost model.ProfilePost) (int, any) {
	user, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	contactUuid := profilePost.ContactUuid
	if len(contactUuid) <= 0 {
		return response.BadRequestMsg("contact_uuid is empty")
	}
	contact, err := repository.ContactRepo.GetContactInfoById(ctx, domainUuid, contactUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if contact == nil {
		return response.BadRequestMsg("contact not found")
	}
	profilesInContact := contact.Profiles
	for i := 0; i < len(profilesInContact); i++ {
		if profilesInContact[i].ProfileType == PRIMARY_PROFILE_TYPE && profilePost.ProfileType == PRIMARY_PROFILE_TYPE {
			return response.BadRequestMsg("contact already has primary profile")
		}
	}
	if profilePost.ProfileType == PRIMARY_PROFILE_TYPE {
		phoneNumbers := make([]string, 0)
		phoneNumbers = append(phoneNumbers, profilePost.PhoneNumber)
		for _, e := range profilePost.Phones {
			phoneNumbers = append(phoneNumbers, e.Data)
		}
		if ct, err := repository.ProfileRepo.GetProfileByPhoneNumber(ctx, domainUuid, phoneNumbers...); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if ct != nil {
			return response.ServiceUnavailableMsg("phone number is existed")
		}
	} else if profilePost.ProfileType == RELATION_PROFILE_TYPE {
		if len(profilePost.RelatedProfileUuid) <= 0 {
			return response.BadRequestMsg("related_profile_uuid is empty")
		}
		profile, err := repository.ProfileRepo.GetProfileInfoById(ctx, domainUuid, profilePost.RelatedProfileUuid)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		if profile == nil {
			return response.BadRequestMsg("related_profile_uuid is not existed")
		}
		if profile.ProfileType != PRIMARY_PROFILE_TYPE {
			return response.BadRequestMsg("related_profile_uuid is not primary profile")
		}
	}
	moreInfoArr := make([]model.ContactMoreInformationData, 0)
	for _, v := range profilePost.MoreInformation {
		if len(v.Attribute) > 0 {
			moreInfoArr = append(moreInfoArr, v)
		}
	}
	moreInformation, _ := util.ParseMapToString(moreInfoArr)
	profileUuid := uuid.NewString()
	birthday := ""
	if len(profilePost.Birthday) > 0 {
		birthdayTmp, err := time.Parse("2006-01-02", profilePost.Birthday)
		if err != nil {
			birthday = profilePost.Birthday
		} else {
			birthday = birthdayTmp.Format("2006-01-02")
		}
	}
	profile := model.Profile{
		ProfileUuid:        profileUuid,
		ContactUuid:        contactUuid,
		DomainUuid:         domainUuid,
		ProfileCode:        profilePost.ProfileCode,
		Fullname:           profilePost.Fullname,
		PhoneNumber:        profilePost.PhoneNumber,
		Email:              profilePost.Email,
		Birthday:           birthday,
		Description:        profilePost.Description,
		Address:            profilePost.Address,
		IdentityNumber:     profilePost.IdentityNumber,
		IdentityIssueOn:    profilePost.IdentityIssueOn,
		IdentityIssueAt:    profilePost.IdentityIssueAt,
		Passport:           profilePost.Passport,
		JobTitle:           profilePost.JobTitle,
		Gender:             profilePost.Gender,
		RefId:              profilePost.RefId,
		RefCode:            profilePost.RefCode,
		MoreInformation:    moreInformation,
		Country:            profilePost.Country,
		Province:           profilePost.Province,
		District:           profilePost.District,
		Ward:               profilePost.Ward,
		ProfileType:        profilePost.ProfileType,
		ProfileName:        profilePost.ProfileName,
		RelatedProfileUuid: profilePost.RelatedProfileUuid,
		Status:             true,
		ContactSocial:      &profilePost.ContactSocial,
		CreatedBy:          userUuid,
		CreatedAt:          time.Now(),
	}
	profilePhones := make([]model.ProfilePhone, 0)
	for _, e := range profilePost.Phones {
		profilePhones = append(profilePhones, model.ProfilePhone{
			ProfileUuid: profile.ProfileUuid,
			DomainUuid:  domainUuid,
			Data:        e.Data,
			Type:        e.Type,
		})
	}
	profileEmails := make([]model.ProfileEmail, 0)
	for _, e := range profilePost.Emails {
		profileEmails = append(profileEmails, model.ProfileEmail{
			ProfileUuid: profile.ProfileUuid,
			DomainUuid:  domainUuid,
			Data:        e.Data,
			EmailType:   e.Type,
		})
	}
	profileOwners := make([]model.ProfileOwner, 0)
	for _, u := range profilePost.UserOwners {
		userAssign, err := repository.UserRepo.GetUserByIdOrUsername(ctx, domainUuid, u.Data)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if userAssign == nil {
			return response.BadRequestMsg("user_owner is not existed")
		}
		profileOwners = append(profileOwners, model.ProfileOwner{
			ProfileUuid: profile.ProfileUuid,
			DomainUuid:  domainUuid,
			UserUuid:    userAssign.UserUuid,
			Username:    userAssign.Username,
			AssignedAt:  time.Now(),
			Type:        u.Type,
		})
	}
	profileNotes := make([]model.ProfileNote, 0)
	if len(profilePost.Note) > 0 {
		profileNotes = append(profileNotes, model.ProfileNote{
			ProfileUuid: profile.ProfileUuid,
			DomainUuid:  domainUuid,
			UserUuid:    user.UserUuid,
			Content:     profilePost.Note,
			CreatedAt:   time.Now(),
		})
	}
	if len(profilePost.ProfileCode) > 0 {
		_, checkProfileCodeIsExisted, err := repository.ProfileRepo.GetProfilesInfo(ctx, domainUuid, model.ProfileFilter{ProfileCode: profile.ProfileCode}, 1, 0)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		if len(*checkProfileCodeIsExisted) > 0 {
			return response.ServiceUnavailableMsg("profile_code is existed")
		}
	}
	if err := repository.ProfileRepo.InsertProfileTransaction(ctx, &profile, profilePhones, profileEmails, profileOwners, profileNotes); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.Created(map[string]any{
		"contact_uuid": contact.ContactUuid,
		"profile_uuid": profile.ProfileUuid,
	})

}

func (p *Profile) PutProfile(ctx context.Context, domainUuid, userUuid, profileUuid string, profilePut model.ProfilePost) (int, any) {
	_, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	profile, err := repository.ProfileRepo.GetProfileById(ctx, domainUuid, profileUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if profile == nil {
		return response.BadRequestMsg("profile not found")
	}
	contactUuid := profilePut.ContactUuid
	if len(contactUuid) <= 0 {
		return response.BadRequestMsg("contact_uuid is empty")
	}
	contact, err := repository.ContactRepo.GetContactInfoById(ctx, domainUuid, contactUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if contact == nil {
		return response.BadRequestMsg("contact not found")
	}
	profilesInContact := contact.Profiles
	for i := 0; i < len(profilesInContact); i++ {
		if profilesInContact[i].ProfileType == PRIMARY_PROFILE_TYPE && profilePut.ProfileType == PRIMARY_PROFILE_TYPE && profileUuid != profilesInContact[i].ProfileUuid {
			return response.BadRequestMsg("contact already has primary profile")
		}
	}
	if profilePut.ProfileType == PRIMARY_PROFILE_TYPE {
		phoneNumbers := make([]string, 0)
		phoneNumbers = append(phoneNumbers, profilePut.PhoneNumber)
		for _, e := range profilePut.Phones {
			phoneNumbers = append(phoneNumbers, e.Data)
		}
		if ct, err := repository.ProfileRepo.GetProfileByPhoneNumber(ctx, domainUuid, phoneNumbers...); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if ct != nil {
			if ct.ProfileUuid != profileUuid {
				return response.ServiceUnavailableMsg("phone number is existed")
			}
		}
	} else if profilePut.ProfileType == RELATION_PROFILE_TYPE {
		if len(profilePut.RelatedProfileUuid) <= 0 {
			return response.BadRequestMsg("related_profile_uuid is empty")
		}
		profile, err := repository.ProfileRepo.GetProfileInfoById(ctx, domainUuid, profilePut.RelatedProfileUuid)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		if profile == nil {
			return response.BadRequestMsg("related_profile_uuid is not existed")
		}
		if profile.ProfileType != PRIMARY_PROFILE_TYPE {
			return response.BadRequestMsg("related_profile_uuid is not primary profile")
		}
	}
	moreInfoArr := make([]model.ContactMoreInformationData, 0)
	for _, v := range profilePut.MoreInformation {
		if len(v.Attribute) > 0 {
			moreInfoArr = append(moreInfoArr, v)
		}
	}
	moreInformation, _ := util.ParseMapToString(moreInfoArr)
	profile.Fullname = profilePut.Fullname
	birthday := ""
	if len(profilePut.Birthday) > 0 {
		birthdayTmp, err := time.Parse("2006-01-02", profilePut.Birthday)
		if err != nil {
			birthday = ""
		} else {
			birthday = birthdayTmp.Format("2006-01-02")
		}
	}
	profile.Status = profilePut.Status
	profile.Birthday = birthday
	profile.PhoneNumber = profilePut.PhoneNumber
	profile.Email = profilePut.Email
	profile.Description = profilePut.Description
	profile.Address = profilePut.Address
	profile.IdentityNumber = profilePut.IdentityNumber
	profile.IdentityIssueOn = profilePut.IdentityIssueOn
	profile.IdentityIssueAt = profilePut.IdentityIssueAt
	profile.Passport = profilePut.Passport
	profile.JobTitle = profilePut.JobTitle
	profile.Gender = profilePut.Gender
	profile.RefId = profilePut.RefId
	profile.RefCode = profilePut.RefCode
	profile.MoreInformation = moreInformation
	profile.Country = profilePut.Country
	profile.Province = profilePut.Province
	profile.District = profilePut.District
	profile.Ward = profilePut.Ward
	profile.ContactSocial = &profilePut.ContactSocial
	profile.UpdatedAt = time.Now()
	profile.ProfileName = profilePut.ProfileName
	profile.ProfileType = profilePut.ProfileType
	profile.ProfileCode = profilePut.ProfileCode
	if len(profilePut.RelatedProfileUuid) > 0 {
		profile.RelatedProfileUuid = profilePut.RelatedProfileUuid
	} else {
		profile.RelatedProfileUuid = profileUuid
	}
	profilePhones := make([]model.ProfilePhone, 0)
	for _, e := range profilePut.Phones {
		profilePhones = append(profilePhones, model.ProfilePhone{
			ProfileUuid: profile.ProfileUuid,
			DomainUuid:  domainUuid,
			Data:        e.Data,
			Type:        e.Type,
		})
	}
	profileEmails := make([]model.ProfileEmail, 0)
	for _, e := range profilePut.Emails {
		profileEmails = append(profileEmails, model.ProfileEmail{
			ProfileUuid: profile.ProfileUuid,
			DomainUuid:  domainUuid,
			Data:        e.Data,
			EmailType:   e.Type,
		})
	}
	profileOwners := make([]model.ProfileOwner, 0)
	for _, u := range profilePut.UserOwners {
		userAssign, err := repository.UserRepo.GetUserByIdOrUsername(ctx, domainUuid, u.Data)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if userAssign == nil {
			return response.BadRequestMsg("user_owner is not existed")
		}
		profileOwners = append(profileOwners, model.ProfileOwner{
			ProfileUuid: profile.ProfileUuid,
			DomainUuid:  domainUuid,
			UserUuid:    userAssign.UserUuid,
			Username:    userAssign.Username,
			AssignedAt:  time.Now(),
			Type:        u.Type,
		})
	}
	// contactNotes := make([]model.ContactNote, 0)
	// if len(profilePut.Note) > 0 {
	// 	contactNotes = append(contactNotes, model.ContactNote{
	// 		ContactUuid: contact.ContactUuid,
	// 		DomainUuid:  domainUuid,
	// 		UserUuid:    user.UserUuid,
	// 		Content:     profilePut.Note,
	// 		CreatedAt:   time.Now(),
	// 	})
	// }
	if len(profilePut.ProfileCode) > 0 {
		_, checkProfileCodeIsExisted, err := repository.ProfileRepo.GetProfilesInfo(ctx, domainUuid, model.ProfileFilter{ProfileCode: profile.ProfileCode}, 1, 0)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		if len(*checkProfileCodeIsExisted) > 0 {
			if (*checkProfileCodeIsExisted)[0].ProfileUuid != profileUuid {
				return response.ServiceUnavailableMsg("profile_code is existed")
			}
		}
	}

	if err := repository.ProfileRepo.UpdateProfileTransaction(ctx, profileUuid, profile, profilePhones, profileEmails, profileOwners); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(map[string]any{
		"contact_uuid": contact.ContactUuid,
		"profile_uuid": profile.ProfileUuid,
	})
}

func (p *Profile) DeleteProfile(ctx context.Context, domainUuid, userUuid, profileUuid string) (int, any) {
	_, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	profile, err := repository.ProfileRepo.GetProfileById(ctx, domainUuid, profileUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if profile == nil {
		return response.BadRequestMsg("profile not found")
	}

	if err := repository.ProfileRepo.DeleteProfileTransaction(ctx, profile.ProfileUuid); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(map[string]any{
		"id": profile.ProfileUuid,
	})
}

func (p *Profile) PatchProfileAvatar(ctx context.Context, domainUuid, userUuid, profileUuid, data string) (int, any) {
	_, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	profile, err := repository.ProfileRepo.GetProfileById(ctx, domainUuid, profileUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if profile == nil {
		return response.BadRequestMsg("profile is not existed")
	}
	dir := util.UPLOAD_DIR + "avatar/"
	fileName := "avatar_" + profile.ProfileUuid
	if fileName, err = util.DecodeAndSaveImageBase64(data, dir, fileName); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err)
	}
	profile.AvatarUrl = fileName
	if err := repository.ProfileRepo.UpdateProfile(ctx, profile); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(map[string]any{
		"profile_uuid": profile.ContactUuid,
	})
}

func (p *Profile) GetAvatar(ctx context.Context, fileName string) (string, error) {
	avatar := ""
	if _, err := os.Stat(util.UPLOAD_DIR + "avatar/" + fileName); err != nil {
		return "", errors.New("file " + fileName + " is not exist")

	}
	avatar = util.UPLOAD_DIR + "avatar/" + fileName
	return avatar, nil
}

func (p *Profile) PatchStatusProfile(ctx context.Context, domainUuid, userUuid, profileId string, status bool) (int, any) {
	profile, err := repository.ProfileRepo.GetProfileById(ctx, domainUuid, profileId)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if profile == nil {
		return response.BadRequestMsg("profile not found")
	}
	profile.Status = status
	profile.UpdatedAt = time.Now()

	if err := repository.ProfileRepo.UpdateProfile(ctx, profile); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error)
	}

	return response.OK(map[string]any{
		"id": profileId,
	})
}

func (s *Profile) PatchSocialMappingContactProfile(ctx context.Context, domainUuid, userUuid, profileId string, data model.SocialMappingContact) (int, any) {
	profile, err := repository.ProfileRepo.GetProfileById(ctx, domainUuid, profileId)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if profile == nil {
		return response.BadRequestMsg("profile not found")
	}

	profile.SocialMappingContact = data
	profile.UpdatedAt = time.Now()

	if err := repository.ProfileRepo.UpdateProfileByField(ctx, domainUuid, *profile); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error)
	}

	return response.OK(map[string]any{
		"id": profileId,
	})
}

func (s *Profile) GetManageProfiles(ctx context.Context, domainUuid, userUuid string, filter model.ProfileFilter, limit, offset int) (int, any) {
	total, data, err := repository.ProfileRepo.GetManageProfiles(ctx, domainUuid, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Pagination(data, total, limit, offset)
}

func (s *Profile) DeleteProfileTransaction(ctx context.Context, domainUuid, userUuid string, filter model.ProfileFilter) (int, any) {
	var profileUuids []string
	var manageProfiles []model.ProfileManageView
	if filter.IsAll {
		filter.ProfileUuids = []string{}
		_, data, err := repository.ProfileRepo.GetManageProfiles(ctx, domainUuid, filter, -1, 0)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		if len(*data) > 0 {
			if err := util.ParseAnyToAny(data, &manageProfiles); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
		}
	} else {
		filter = model.ProfileFilter{
			ProfileUuids: filter.ProfileUuids,
		}
		limit := -1
		if len(filter.ProfileUuids) > 0 {
			limit = len(filter.ProfileUuids)
		}
		_, data, err := repository.ProfileRepo.GetManageProfiles(ctx, domainUuid, filter, limit, 0)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		if len(*data) > 0 {
			if err := util.ParseAnyToAny(data, &manageProfiles); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
		}
	}

	if len(manageProfiles) > 0 {
		var profiles []model.Profile
		for _, item := range manageProfiles {
			profileUuids = append(profileUuids, item.ProfileUuid)
			var profile model.Profile
			if err := util.ParseAnyToAny(item, &profile); err != nil {
				log.Error(err)
				break
			}
			profiles = append(profiles, profile)
		}

		var tickets []model.Ticket
		if len(profileUuids) > 0 {
			ticketTmps, err := repository.TicketRepo.GetTicketByProfileUuids(ctx, domainUuid, profileUuids)
			if err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
			tickets = *ticketTmps
		}

		if err := repository.ProfileRepo.DeleteProfileWithTicketTransaction(ctx, domainUuid, profiles, tickets); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
	}

	return response.NewOKResponse("success")
}

func (s *Profile) PostConvertLeadToProfile(ctx context.Context, domainUuid, userUuid string, profilePost model.ProfilePost) (int, any) {
	user, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	leadExist, err := repository.LeadRepo.GetLeadById(ctx, domainUuid, profilePost.LeadUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if leadExist == nil {
		return response.ServiceUnavailableMsg("lead not found")
	} else if leadExist.IsConvertProfile {
		return response.ServiceUnavailableMsg("lead has been converted to profile")
	}

	profile := model.Profile{}
	if err := util.ParseAnyToAny(profilePost, &profile); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if len(profilePost.RelatedProfileUuid) > 0 {
		filter := model.ProfileFilter{
			ProfileType:  "primary",
			ProfileUuids: []string{profilePost.RelatedProfileUuid},
		}
		_, profiles, err := repository.ProfileRepo.GetManageProfiles(ctx, domainUuid, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		if len(*profiles) < 1 {
			return response.ServiceUnavailableMsg("profile relation does not exist")
		}
		profile.ContactUuid = (*profiles)[0].ContactUuid
		if profile.ProfileType != "primary" {
			profile.ProfileType = "primary"
			profile.ProfileName = ""
			leadExist.RelatedProfileType = "primary"
		}
	} else {
		profileFilter := model.ProfileFilter{
			PhoneNumber: profilePost.PhoneNumber,
			ProfileType: "primary",
		}
		_, profiles, err := repository.ProfileRepo.GetProfilesInfo(ctx, domainUuid, profileFilter, 1, 0)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		contactUuid := ""
		if len(*profiles) > 0 && (profilePost.ProfileType == "primary" || len(profilePost.ProfileName) < 1) {
			return response.ServiceUnavailableMsg("profile exist with profile name " + (*profiles)[0].Fullname)
		}
		if len(*profiles) > 0 {
			contactUuid = (*profiles)[0].ContactUuid
		} else if len(*profiles) < 1 {
			contactUuid = uuid.NewString()
			contact := model.Contact{
				ContactUuid: contactUuid,
				DomainUuid:  domainUuid,
				UnitUuid:    user.UnitUuid,
				ContactType: "",
				ContactName: profile.ProfileName,
				Status:      true,
				SourceName:  "",
				SourceUuid:  "",
				CreatedBy:   userUuid,
				CreatedAt:   time.Now(),
			}
			if err := repository.ContactRepo.InsertContact(ctx, &contact); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
			if profile.ProfileType != "primary" {
				profile.ProfileType = "primary"
				profile.ProfileName = ""
				leadExist.RelatedProfileType = "primary"
			}
		}

		profile.ContactUuid = contactUuid
	}

	profile.DomainUuid = domainUuid
	profile.ProfileUuid = uuid.NewString()
	profile.CreatedAt = time.Now()
	profile.CreatedBy = userUuid

	if err := repository.ProfileRepo.InsertProfileTransaction(ctx, &profile, []model.ProfilePhone{}, []model.ProfileEmail{}, []model.ProfileOwner{}, []model.ProfileNote{}); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	leadExist.RelatedProfileUuid = profile.ProfileUuid
	leadExist.RelatedProfileType = profile.ProfileName
	leadExist.IsConvertProfile = true
	leadExist.UpdatedAt = time.Now()
	if err := repository.LeadRepo.UpdateLead(ctx, leadExist); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]string{
		"id": profile.ProfileUuid,
	})
}
