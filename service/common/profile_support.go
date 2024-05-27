package common

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/util"
	"time"
)

func ParseProfileInfo(data *model.ProfileView) model.ProfileInfo {
	moreInformation := make([]any, 0)
	if len(data.MoreInformationStr) > 0 {
		if err := util.ParsesStringToStruct(data.MoreInformationStr, &moreInformation); err != nil {
			log.Error(err)
		}
	}

	profile := model.ProfileInfo{
		ProfileUuid:          data.ProfileUuid,
		AvatarUrl:            data.AvatarUrl,
		Fullname:             data.Fullname,
		Birthday:             data.Birthday,
		PhoneNumber:          data.PhoneNumber,
		Email:                data.Email,
		Description:          data.Description,
		Address:              data.Address,
		IdentityNumber:       data.IdentityNumber,
		IdentityIssueOn:      data.IdentityIssueOn,
		IdentityIssueAt:      data.IdentityIssueAt,
		Passport:             data.Passport,
		JobTitle:             data.JobTitle,
		MoreInformation:      moreInformation,
		RefId:                data.RefId,
		RefCode:              data.RefCode,
		Gender:               data.Gender,
		Country:              data.Country,
		Province:             data.Province,
		District:             data.District,
		Ward:                 data.Ward,
		Status:               data.Status,
		ContactSocial:        data.ContactSocial,
		ProfileType:          data.ProfileType,
		ProfileName:          data.ProfileName,
		ProfileCode:          data.ProfileCode,
		ListRelatedProfile:   data.ListRelatedProfile,
		ContactUuid:          data.ContactUuid,
		RelatedProfileUuid:   data.RelatedProfileUuid,
		SocialMappingContact: &data.SocialMappingContact,
	}

	profile.UserOwners = make([]model.UserOwnerShortInfo, 0)
	for _, e := range data.UserOwners {
		profile.UserOwners = append(profile.UserOwners, model.UserOwnerShortInfo{
			Data: e.UserUuid,
			Type: e.Type,
			Info: struct {
				UserUuid   string    "json:\"user_uuid\" bun:\"user_uuid,pk\""
				Username   string    "json:\"username\" bun:\"username\""
				AssignedAt time.Time "json:\"assigned_at\" bun:\"assigned_at\""
			}{
				UserUuid:   e.UserUuid,
				Username:   e.Username,
				AssignedAt: e.AssignedAt,
			},
		})
	}
	profile.Phones = make([]model.ContactMapData, 0)
	if len(data.Phones) > 0 {
		for _, e := range data.Phones {
			profile.Phones = append(profile.Phones, model.ContactMapData{
				Data: e.Data,
				Type: e.Type,
			})
		}
	} else {
		profile.Phones = append(profile.Phones, model.ContactMapData{
			Data: "",
			Type: "personal",
		})
	}

	profile.Emails = make([]model.ContactMapData, 0)
	if len(data.Emails) > 0 {
		for _, e := range data.Emails {
			profile.Emails = append(profile.Emails, model.ContactMapData{
				Data: e.Data,
				Type: e.EmailType,
			})
		}
	} else {
		profile.Emails = append(profile.Emails, model.ContactMapData{
			Data: "",
			Type: "personal",
		})
	}
	return profile

}
