package service

import (
	"contactcenter-api/common/constants"
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	"contactcenter-api/internal/freeswitch"
	"contactcenter-api/repository"
	"contactcenter-api/service/common"
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/dchest/uniuri"
	"github.com/google/uuid"
)

type (
	IUserCrm interface {
		GetUserCrms(ctx context.Context, domainUuid, userUuid string, limit, offset int, filter model.UserFilter) (int, any)
		PostUserCrm(ctx context.Context, domainUuid, userUuid string, userPost model.UserPost) (int, any)
		GetUserCrmById(ctx context.Context, domainUuid, userUuid, id string) (int, any)
		GetUserCrmViewById(ctx context.Context, domainUuid, userUuid, id string) (int, any)
		PutUserCrmById(ctx context.Context, domainUuuid, userUuid, id string, data model.UserPost) (int, any)
		PatchUserCrm(ctx context.Context, domainUuid, userUuid, id, unitUuid, roleGroupUuid string) (int, any)
		DeleteUserCrmById(ctx context.Context, domainUuid, userUuid, id string) (int, any)
		PostExportUsers(ctx context.Context, domainUuid, userUuid string, limit, offset int, filter model.UserFilter) (int, any)
	}
	UserCrm struct{}
)

func NewUserCrm() IUserCrm {
	return &UserCrm{}
}

func (s *UserCrm) PostUserCrm(ctx context.Context, domainUuid, userUuid string, userPost model.UserPost) (int, any) {
	domain, err := repository.DomainRepo.GetDomainById(ctx, domainUuid)
	if err != nil {
		return response.ServiceUnavailableMsg(err.Error())
	} else if domain == nil {
		return response.BadRequestMsg("domain error")
	}
	active := "true"
	if userPost.UserEnabled == "false" {
		active = "false"
	}

	if userRes, err := repository.UserRepo.GetUserOfDomainByName(ctx, domainUuid, userPost.Username); err != nil {
		return response.ServiceUnavailableMsg(err.Error())
	} else if userRes != nil {
		return response.BadRequestMsg("username is already taken")
	}
	roleGroup, err := repository.RoleGroupRepo.GetRoleGroupById(ctx, domainUuid, userPost.RoleUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if roleGroup == nil {
		return response.BadRequestMsg("id is not exists")
	}
	unit, err := repository.UnitRepo.GetUnitById(ctx, domainUuid, userPost.UnitUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(unit.UnitUuid) < 1 {
		return response.BadRequestMsg("unit_uuid is not exists")
	}
	user := model.User{
		UserUuid:     uuid.NewString(),
		DomainUuid:   domainUuid,
		Username:     userPost.Username,
		Salt:         uuid.NewString(),
		UserStatus:   "Logged Out",
		ApiKey:       uuid.NewString(),
		AddDate:      util.TimeToString(time.Now()),
		UserEnabled:  active,
		AddUser:      "",
		Password:     userPost.Password,
		Level:        userPost.Level,
		RoleUuid:     userPost.RoleUuid,
		UnitUuid:     userPost.UnitUuid,
		EnableWebrtc: userPost.EnableWebrtc,
	}
	dbPassword := []byte(user.Salt + user.Password)
	user.Password = fmt.Sprintf("%x", md5.Sum(dbPassword))

	vContact := model.VContact{
		DomainUuid:        domainUuid,
		ContactUuid:       uuid.NewString(),
		ContactType:       "user",
		ContactNameGiven:  userPost.FirstName,
		ContactNameMiddle: userPost.MiddleName,
		ContactNameFamily: userPost.LastName,
		ContactNickname:   userPost.Username,
	}

	// contact crm
	contactEmail := model.VContactEmail{
		DomainUuid:       domainUuid,
		ContactEmailUuid: uuid.NewString(),
		ContactUuid:      vContact.ContactUuid,
		EmailAddress:     userPost.Email,
		EmailPrimary:     1,
	}
	roleGroup.RoleGroupUuid = userPost.RoleUuid

	user.ContactUuid = vContact.ContactUuid
	extension := model.Extension{}
	callCenterAgent := model.CallCenterAgent{}
	if userPost.IsCreateExtension {
		extensionChars := "0123456789"
		extension = model.Extension{
			DomainUuid:            domainUuid,
			UserContext:           domain.DomainName,
			ExtensionUuid:         uuid.NewString(),
			LimitMax:              "5",
			DirectoryVisible:      "true",
			DirectoryExtenVisible: "true",
			LimitDestination:      "error/user_busy",
			CallTimeout:           30,
			CallScreenEnabled:     "false",
			UserRecord:            "all",
			Enabled:               active,
			Password:              userPost.Password,
		}
		if len(userPost.Extension) > 0 {
			if exten, err := repository.ExtensionRepo.GetExtensionByExten(ctx, domainUuid, userPost.Extension); err != nil {
				return response.ServiceUnavailableMsg(err.Error())
			} else if exten != nil {
				return response.BadRequestMsg("extension is existed")
			}
			extension.Extension = userPost.Extension
		} else {
			extension.Extension = uniuri.NewLenChars(8, []byte(extensionChars))
		}
		// Loop to get unique extension in domain
		for {
			if res, err := repository.ExtensionRepo.GetExtensionByExten(ctx, domainUuid, extension.Extension); err != nil {
				return response.ServiceUnavailableMsg(err.Error())
			} else if res != nil {
				extension.Extension = uniuri.NewLenChars(8, []byte(extensionChars))
			} else {
				break
			}
		}
		extensionUser := model.ExtensionUser{
			DomainUuid:        domainUuid,
			ExtensionUuid:     extension.ExtensionUuid,
			ExtensionUserUuid: uuid.NewString(),
			UserUuid:          user.UserUuid,
		}
		callCenterAgent = model.CallCenterAgent{
			DomainUuid:             domainUuid,
			CallCenterAgentUuid:    uuid.NewString(),
			UserUuid:               user.UserUuid,
			AgentName:              user.Username,
			AgentType:              "callback",
			AgentCallTimeout:       15,
			AgentId:                extension.Extension,
			AgentPassword:          "",
			AgentContact:           "user" + "/" + extension.Extension + "@" + domain.DomainName,
			AgentStatus:            "Logged Out",
			AgentMaxNoAnswer:       0,
			AgentWrapUpTime:        1,
			AgentRejectDelayTime:   1,
			AgentNoAnswerDelayTime: "1",
			AgentBusyDelayTime:     1,
		}
		userPost.Extension = extension.Extension
		if err := repository.AgentRepo.InsertAgentTransaction(ctx, &user, &vContact, nil, &extension,
			&extensionUser, &callCenterAgent, &contactEmail, roleGroup); err != nil {
			return response.ServiceUnavailableMsg(err.Error())
		}
	} else {
		if err := repository.UserRepo.InsertUserTransaction(ctx, &user, &vContact, nil, &contactEmail, roleGroup); err != nil {
			return response.ServiceUnavailableMsg(err.Error())
		}
	}
	go AddTransaction(domainUuid, userUuid, "user_crm", user.UserUuid, "add", "done", "", nil, map[string]any{
		"user":              user,
		"extension":         extension,
		"contact":           vContact,
		"contact_email":     contactEmail,
		"call_center_agent": callCenterAgent,
		"role_group":        roleGroup,
	})
	return response.Created(map[string]any{
		"user_uuid":      user.UserUuid,
		"extension_uuid": extension.ExtensionUuid,
	})
}

func (s *UserCrm) GetUserCrms(ctx context.Context, domainUuid, userUuid string, limit, offset int, filter model.UserFilter) (int, any) {
	tree := map[string]any{}
	if len(filter.UnitUuid) > 0 {
		unit, err := repository.UnitRepo.GetUnitById(ctx, domainUuid, filter.UnitUuid)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if unit == nil {
			return response.BadRequestMsg("unit_uuid is not exists")
		}
		common.ProcessTree(ctx, domainUuid, tree, *unit, unit.Level)
	}
	if len(filter.UnitUuid) > 0 {
		for _, val := range tree {
			if val != nil {
				items := []model.UnitInfo{}
				if err := util.ParseAnyToAny(val, &items); err != nil {
					log.Error(err)
					continue
				}
				for _, item := range items {
					if len(item.Users) > 0 {
						for _, v := range item.Users {
							filter.ManageUserUuids = append(filter.ManageUserUuids, v.UserUuid)
						}
					}
				}
			}
		}
	}
	if len(filter.ManageUserUuids) > 0 {
		filter.UnitUuid = ""
	}

	total, result, err := repository.UserCrmRepo.GetUserCrms(ctx, domainUuid, limit, offset, filter)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.Pagination(result, total, limit, offset)
}

func (s *UserCrm) GetUserCrmById(ctx context.Context, domainUuid, userUuid, id string) (int, any) {
	userCrmExist, err := repository.UserCrmRepo.GetUserCrmById(ctx, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(userCrmExist)
}

func (s *UserCrm) GetUserCrmViewById(ctx context.Context, domainUuid, userUuid, id string) (int, any) {
	userCrmExist, err := repository.UserCrmRepo.GetUserCrmById(ctx, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if userCrmExist == nil {
		return response.NotFoundMsg("user_crm is not exists")
	}
	userCrmExist.RoleGroups.PermissionMainOptimize = common.HandlePermissionMainOptimize(userCrmExist.RoleGroups.PermissionMain)
	userCrmExist.RoleGroups.PermissionAdvanceOptimize = common.HandlePermissionAdvanceOptimize(userCrmExist.RoleGroups.PermissionAdvance)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(userCrmExist)
}

func (s *UserCrm) PutUserCrmById(ctx context.Context, domainUuid, userUuid, id string, data model.UserPost) (int, any) {
	domain, err := repository.DomainRepo.GetDomainById(ctx, domainUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if domain == nil {
		return response.BadRequestMsg("domain error")
	}
	user, err := repository.UserRepo.GetUserById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if user == nil {
		return response.NotFoundMsg("user is not exist")
	}
	oldData := map[string]any{
		"user": user,
	}
	// check username is existed
	if userTmp, err := repository.UserRepo.GetUserByUsername(ctx, domainUuid, data.Username); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if userTmp != nil && id != user.UserUuid {
		return response.BadRequestMsg("username is already taken")
	}
	roleGroup, err := repository.RoleGroupRepo.GetRoleGroupById(ctx, domainUuid, data.RoleUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if roleGroup == nil {
		return response.BadRequestMsg("role_uuid is not exists")
	}
	unit, err := repository.UnitRepo.GetUnitById(ctx, domainUuid, data.UnitUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(unit.UnitUuid) < 1 {
		return response.BadRequestMsg("unit_uuid is not exists")
	}
	userUpdate := model.User{
		UserUuid:     user.UserUuid,
		DomainUuid:   user.DomainUuid,
		Username:     data.Username,
		UserEnabled:  data.UserEnabled,
		Level:        data.Level,
		RoleUuid:     data.RoleUuid,
		UnitUuid:     data.UnitUuid,
		EnableWebrtc: data.EnableWebrtc,
	}
	if data.Password != "default" {
		userUpdate.Salt = uuid.NewString()
		dbPassword := []byte(userUpdate.Salt + data.Password)
		userUpdate.Password = fmt.Sprintf("%x", md5.Sum(dbPassword))
	} else {
		userUpdate.Salt = user.Salt
		userUpdate.Password = user.Password
	}
	contact := new(model.VContact)
	if len(user.ContactUuid) > 0 {
		contact, err = repository.UserRepo.GetContactById(ctx, domainUuid, user.ContactUuid)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if contact != nil {
			oldData["contact"] = *contact
			contact.ContactNameGiven = data.FirstName
			contact.ContactNameMiddle = data.MiddleName
			contact.ContactNameFamily = data.LastName
			contact.ContactNickname = data.Username
		} else {
			contact = &model.VContact{
				DomainUuid:        domainUuid,
				ContactUuid:       uuid.NewString(),
				ContactType:       "user",
				ContactNameGiven:  data.FirstName,
				ContactNameMiddle: data.MiddleName,
				ContactNameFamily: data.LastName,
				ContactNickname:   data.Username,
			}
		}
	} else {
		contact = &model.VContact{
			DomainUuid:        domainUuid,
			ContactUuid:       uuid.NewString(),
			ContactType:       "user",
			ContactNameGiven:  data.FirstName,
			ContactNameMiddle: data.MiddleName,
			ContactNameFamily: data.LastName,
			ContactNickname:   data.Username,
		}
	}
	contactEmail := model.VContactEmail{
		DomainUuid:       domainUuid,
		ContactEmailUuid: uuid.NewString(),
		ContactUuid:      contact.ContactUuid,
		EmailAddress:     data.Email,
		EmailPrimary:     1,
	}
	userUpdate.ContactUuid = contact.ContactUuid
	extension := model.Extension{}
	callCenterAgent := model.CallCenterAgent{}
	extensionUuid := ""
	if !data.IsRemapping {
		extensionOld, err := repository.ExtensionRepo.GetExtensionByUserUuid(ctx, domainUuid, id)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if extensionOld != nil && extensionOld.Extension != data.Extension {
			return response.BadRequestMsg("user is already has extension")
		} else if extensionOld != nil && extensionOld.Extension == data.Extension {
			oldData["extension"] = *extensionOld
			data.IsCreateExtension = false
		}
		if data.IsCreateExtension {
			extensionChars := "0123456789"
			extension = model.Extension{
				DomainUuid:            domainUuid,
				UserContext:           domain.DomainName,
				ExtensionUuid:         uuid.NewString(),
				LimitMax:              "5",
				DirectoryVisible:      "true",
				DirectoryExtenVisible: "true",
				LimitDestination:      "error/user_busy",
				CallTimeout:           30,
				CallScreenEnabled:     "false",
				UserRecord:            "all",
				Enabled:               data.UserEnabled,
				Password:              data.Password,
			}
			if len(data.Extension) > 0 {
				if exten, err := repository.ExtensionRepo.GetExtensionByExten(ctx, domainUuid, data.Extension); err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				} else if exten != nil {
					return response.BadRequestMsg("extension is existed")
				}
				extension.Extension = data.Extension
			} else {
				extension.Extension = uniuri.NewLenChars(8, []byte(extensionChars))
			}

			// Loop to get unique extension in domain
			for {
				if res, err := repository.ExtensionRepo.GetExtensionByExten(ctx, domainUuid, extension.Extension); err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				} else if res != nil {
					extension.Extension = uniuri.NewLenChars(8, []byte(extensionChars))
				} else {
					break
				}
			}
			extensionUser := model.ExtensionUser{
				DomainUuid:        domainUuid,
				ExtensionUuid:     extension.ExtensionUuid,
				ExtensionUserUuid: uuid.NewString(),
				UserUuid:          user.UserUuid,
			}
			callCenterAgent = model.CallCenterAgent{
				DomainUuid:             domainUuid,
				CallCenterAgentUuid:    uuid.NewString(),
				UserUuid:               user.UserUuid,
				AgentName:              user.Username,
				AgentType:              "callback",
				AgentCallTimeout:       15,
				AgentId:                extension.Extension,
				AgentPassword:          "",
				AgentContact:           "user" + "/" + extension.Extension + "@" + domain.DomainName,
				AgentStatus:            "Logged Out",
				AgentMaxNoAnswer:       0,
				AgentWrapUpTime:        1,
				AgentRejectDelayTime:   1,
				AgentNoAnswerDelayTime: "1",
				AgentBusyDelayTime:     1,
			}
			if err := repository.AgentRepo.UpdateAgentTransaction(ctx, &userUpdate, contact, &extension, &extensionUser, &callCenterAgent, &contactEmail); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
		} else {
			if err := repository.UserRepo.UpdateUserTransaction(ctx, &userUpdate, contact, &contactEmail); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
			if extensionOld != nil {
				extensionUuid = extensionOld.ExtensionUuid
			}
		}
	} else {
		if len(data.Extension) < 1 {
			return response.BadRequestMsg("extension is missing")
		}
		// isDeleteExtension := false
		// extensionOld, err := repository.ExtensionRepo.GetExtensionByUserUuid(ctx, domainUuid, id)
		// if err != nil {
		// 	log.Error(err)
		// 	return response.ServiceUnavailableMsg(err.Error())
		// } else if extensionOld != nil && extensionOld.Extension != data.Extension {
		// 	isDeleteExtension = true
		// }
		extensionTmp, err := repository.ExtensionRepo.FindExtensionDataOfExten(ctx, domainUuid, data.Extension)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if extensionTmp == nil {
			return response.BadRequestMsg("extension is invalid")
		} else if extensionTmp != nil && len(extensionTmp.UserUuid) > 0 && extensionTmp.UserUuid != id {
			return response.BadRequestMsg("extension is mapping with other user")
		}
		extensionUuid = extensionTmp.ExtensionUuid
		extensionUser := model.ExtensionUser{
			DomainUuid:        domainUuid,
			ExtensionUuid:     extensionTmp.ExtensionUuid,
			ExtensionUserUuid: uuid.NewString(),
			UserUuid:          userUpdate.UserUuid,
		}
		callCenterAgent = model.CallCenterAgent{
			DomainUuid:             domainUuid,
			CallCenterAgentUuid:    uuid.NewString(),
			UserUuid:               userUpdate.UserUuid,
			AgentName:              userUpdate.Username,
			AgentType:              "callback",
			AgentCallTimeout:       15,
			AgentId:                extensionTmp.Extension,
			AgentPassword:          "",
			AgentContact:           "user" + "/" + extensionTmp.Extension + "@" + domain.DomainName,
			AgentStatus:            "Logged Out",
			AgentMaxNoAnswer:       0,
			AgentWrapUpTime:        1,
			AgentRejectDelayTime:   1,
			AgentNoAnswerDelayTime: "1",
			AgentBusyDelayTime:     1,
		}
		if err := repository.AgentRepo.UpdateAgentTransaction(ctx, &userUpdate, contact, nil, &extensionUser, &callCenterAgent, &contactEmail); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
	}
	if err := repository.UserRepo.UpdateUserActive(ctx, domainUuid, id, extensionUuid, data.UserEnabled == "true"); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	freeswitch.ESLClient.ClearCacheOfDomain(domain.DomainName)
	go AddTransaction(domainUuid, userUuid, "user", id, "update", "done", "", oldData, map[string]any{
		"user":      *user,
		"extension": extension,
		"contact":   *contact,
	})
	return response.OK(map[string]any{
		"extension_uuid": extensionUuid,
		"user_uuid":      user.UserUuid,
	})
}

func (s *UserCrm) DeleteUserCrmById(ctx context.Context, domainUuid, userUuid, id string) (int, any) {
	user, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		return response.ServiceUnavailableMsg(err.Error())
	}
	if user.Level == constants.LEADER || user.Level == constants.MANAGER {
		manageUserUuids := make([]string, 0)
		for _, u := range user.ManageUsers {
			manageUserUuids = append(manageUserUuids, u.UserUuid)
		}
		if !util.InArray(userUuid, manageUserUuids) {
			return response.Forbidden()
		}
	}
	userDelete, err := repository.UserCrmRepo.GetUserCrmById(ctx, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if userDelete == nil {
		return response.NotFoundMsg("user is not exist")
	} else if userDelete.UserEnabled != "false" {
		return response.BadRequestMsg("user must be deactive")
	}
	if err := repository.UserRepo.DeleteUserTransaction(ctx, domainUuid, id, ""); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	go AddTransaction(domainUuid, userUuid, "user", id, "delete", "done", "", *userDelete, nil)

	return response.OK(map[string]any{
		"user_uuid": user.UserUuid,
	})
}

func (s *UserCrm) PatchUserCrm(ctx context.Context, domainUuid, userUuid, id, unitUuid, roleGroupUuid string) (int, any) {
	unitExist, err := repository.UnitRepo.GetUnitById(ctx, domainUuid, unitUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if unitExist == nil {
		return response.ServiceUnavailableMsg("unit is not exist")
	}

	roleGroupExist, err := repository.RoleGroupRepo.GetRoleGroupById(ctx, domainUuid, roleGroupUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if roleGroupExist == nil {
		return response.ServiceUnavailableMsg("role group is not exist")
	}

	userExist, err := repository.UserCrmRepo.GetUserCrmById(ctx, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if userExist == nil {
		return response.ServiceUnavailableMsg("user is not exist")
	}

	if err := repository.UserCrmRepo.PatchUserCrm(ctx, domainUuid, id, unitUuid, roleGroupUuid); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"user_uuid": id,
	})
}
