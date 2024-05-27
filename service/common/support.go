package common

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/util"
)

func ParseContactInfo(data *model.ContactView) model.ContactInfo {
	profile := data.Profiles
	for _, e := range profile {
		moreInformation := make([]any, 0)
		if len(e.MoreInformationStr) > 0 {
			if err := util.ParsesStringToStruct(e.MoreInformationStr, &moreInformation); err != nil {
				log.Error(err)
			}
		}
	}

	contact := model.ContactInfo{
		ContactUuid: data.ContactUuid,
		Status:      data.Status,
		ContactType: data.ContactType,
		ContactName: data.ContactName,
		Profiles:    profile,
		UnitUuid:    data.UnitUuid,
	}
	return contact
}
