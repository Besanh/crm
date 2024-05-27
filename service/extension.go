package service

import (
	"contactcenter-api/common/constants"
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	qrcodeUtil "contactcenter-api/common/qrcode"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	"contactcenter-api/internal/freeswitch"
	"contactcenter-api/repository"
	"contactcenter-api/service/common"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type IExtension interface {
	CheckExtensionIsLoggedIn(ctx context.Context, domainUuid, extension string) (int, any)
	GetExtensions(ctx context.Context, domainUuid string, userUuid string, filter model.ExtensionFilter, limit, offset int) (int, any)
	GetExtensionsInUnit(ctx context.Context, domainUuid string, userUuid string, filter model.ExtensionFilter, limit, offset int) (int, any)
	GetExtensionQrCode(ctx context.Context, domainUuid, userUuid, app string) (int, any)
	PostExtension(ctx context.Context, domainUuid, userUuid string, extensionPost model.ExtensionPost) (int, any)
	PutExtension(ctx context.Context, domainUuid, userUuid string, extensionUuid string, extensionPost model.ExtensionPost) (int, any)
	DeleteExtension(ctx context.Context, domainUuid, userUuid string, extensionUuid string) (int, any)
	GetExtensionByIdOrExten(ctx context.Context, domainUuid string, userUuid string, id string) (int, any)
	PatchExtension(ctx context.Context, domainUuid, userUuid string, extensionUuid string, status bool) (int, any)
}

type Extension struct {
}

func NewExtension() IExtension {
	return &Extension{}
}

func (service *Extension) CheckExtensionIsLoggedIn(ctx context.Context, domainUuid, extension string) (int, any) {
	domain, err := repository.DomainRepo.GetDomainById(ctx, domainUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if domain == nil {
		return response.NotFoundMsg("domain is not existed")
	}
	if extensionEntry, err := repository.ExtensionRepo.GetExtensionByExten(ctx, domainUuid, extension); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if extensionEntry == nil {
		return response.NotFoundMsg("extension is not existed")
	}
	extensionRegisRes, err := repository.SipRegistrationRepo.GetSipRegistrationOfExtension(ctx, domain.DomainName, extension)
	if err != nil {
		return response.ServiceUnavailableMsg(err.Error())
	} else if extensionRegisRes == nil {
		return response.OK(map[string]any{
			"extension":    extension,
			"is_logged_in": false,
		})
	}
	return response.OK(map[string]any{
		"extension":    extension,
		"is_logged_in": true,
	})
}
func (s *Extension) GetExtensions(ctx context.Context, domainUuid string, userUuid string, filter model.ExtensionFilter, limit, offset int) (int, any) {
	user, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		return response.ServiceUnavailableMsg(err.Error())
	}

	if user.Level == constants.MANAGER {
		filter.ManageExtensionUuids = make([]string, 0)
		for _, e := range user.ManageExtensions {
			filter.ManageExtensionUuids = append(filter.ManageExtensionUuids, e.ExtensionUuid)
		}
	} else if user.Level == constants.LEADER {
		filter.ManageExtensionUuids = make([]string, 0)
		for _, e := range user.ManageExtensions {
			filter.ManageExtensionUuids = append(filter.ManageExtensionUuids, e.ExtensionUuid)
		}
	}
	agentExtensions, total, err := repository.ExtensionRepo.GetExtensionsInfo(ctx, domainUuid, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	// data := *agentExtensions
	// for i := 0; i < len(data); i++ {
	// 	if len(data[i].FollowMeUuid) > 0 {
	// 		if followMeDestination, err := repository.FollowMeRepo.GetFollowMeDestinationByFollowMeUuid(ctx, domainUuid, data[i].FollowMeUuid); err != nil {
	// 			log.Error(err)
	// 		} else if followMeDestination != nil {
	// 			data[i].FollowMe = followMeDestination.FollowMeDestination
	// 		}
	// 	}
	// }
	return response.Pagination(agentExtensions, total, limit, offset)
}

func (s *Extension) GetExtensionsInUnit(ctx context.Context, domainUuid string, userUuid string, filter model.ExtensionFilter, limit, offset int) (int, any) {
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
							extension, err := repository.ExtensionRepo.GetExtensionByUserUuid(ctx, domainUuid, v.UserUuid)
							if err != nil {
								log.Error(err)
								continue
							}
							if extension == nil {
								log.Infof("Cannt get extension of user %s", v.UserUuid)
								continue
							}
							filter.ManageExtensionUuids = append(filter.ManageExtensionUuids, extension.ExtensionUuid)
						}
					}
				}
			}
		}
	}
	if len(filter.ManageExtensionUuids) > 0 {
		filter.UnitUuid = ""
	}
	agentExtensions, total, err := repository.ExtensionRepo.GetExtensionsInfo(ctx, domainUuid, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	// data := *agentExtensions
	// for i := 0; i < len(data); i++ {
	// 	if len(data[i].FollowMeUuid) > 0 {
	// 		if followMeDestination, err := repository.FollowMeRepo.GetFollowMeDestinationByFollowMeUuid(ctx, domainUuid, data[i].FollowMeUuid); err != nil {
	// 			log.Error(err)
	// 		} else if followMeDestination != nil {
	// 			data[i].FollowMe = followMeDestination.FollowMeDestination
	// 		}
	// 	}
	// }
	return response.Pagination(agentExtensions, total, limit, offset)
}
func (s *Extension) GetExtensionByIdOrExten(ctx context.Context, domainUuid string, userUuid string, id string) (int, any) {
	extension, err := repository.ExtensionRepo.GetExtensionInfoByExtensionUuid(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if extension == nil {
		return response.NotFoundMsg("extension is not existed")
	}
	// if len(extension.FollowMeUuid) > 0 {
	// 	if followMeDestination, err := repository.FollowMeRepo.GetFollowMeDestinationByFollowMeUuid(ctx, domainUuid, extension.FollowMeUuid); err != nil {
	// 		log.Error(err)
	// 	} else if followMeDestination != nil {
	// 		extension.FollowMe = followMeDestination.FollowMeDestination
	// 	}
	// }
	result := map[string]any{
		"extension_uuid":      extension.ExtensionUuid,
		"extension":           extension.Extension,
		"user_uuid":           extension.UserUuid,
		"username":            extension.Username,
		"enabled":             extension.Enabled,
		"domain_uuid":         extension.DomainUuid,
		"domain_name":         extension.DomainName,
		"follow_me":           extension.FollowMe,
		"is_link_call_center": extension.IsLinkCallCenter,
	}
	return response.OK(result)
}

func (s *Extension) GetExtensionQrCode(ctx context.Context, domainUuid, userUuid, app string) (int, any) {
	extension, err := repository.ExtensionRepo.GetExtensionInfoByExtensionUuid(ctx, domainUuid, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if extension == nil {
		return response.NotFound()
	}
	if app == "pitel_connect" {
		value := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>`+
			`<AccountConfig version="1">`+
			`<Account>`+
			`<APIDomain>%s</APIDomain>`+
			`<WssUrl>%s</WssUrl>`+
			`<UserName>%s</UserName>`+
			`<RegisterServer>%s</RegisterServer>`+
			`<OutboundServer>%s</OutboundServer>`+
			`<UserID>%s</UserID>`+
			`<AuthID>%s</AuthID>`+
			`<AuthPass>%s</AuthPass>`+
			`<AccountName>%s</AccountName>`+
			`<DisplayName>%s</DisplayName>`+
			`<Dialplan />`+
			`<RandomPort />`+
			`<SecOutboundServer />`+
			`<Voicemail />`+
			`</Account>`+
			`</AccountConfig>`,
			PBXInfo.APIDomain,
			PBXInfo.PBXWss,
			fmt.Sprintf("%s@%s", extension.Username, extension.DomainName),
			extension.DomainName,
			PBXInfo.PBXOutboundProxy,
			extension.Extension,
			extension.Extension,
			extension.Password,
			extension.Extension,
			extension.Extension+"@"+extension.DomainName,
		)
		bytes, err := qrcodeUtil.GenerateQrCode(value)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error)
		}
		result := qrcodeUtil.ByteToString(bytes)
		return response.Data(http.StatusOK, result)
	} else if app == "zopier" {
		value := fmt.Sprintf("https://oem.zoiper.com/qr.php?provider_id=c663604a0e85b7c48948e3d88094ec37&u=%s&h=%s&p=%s&o=%s&t=&x=&a=%s&tr=0",
			extension.Extension,
			extension.DomainName,
			extension.Password,
			PBXInfo.PBXOutboundProxy,
			extension.Extension,
		)
		client := resty.New()
		resp, err := client.R().Get(value)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error)
		}
		result := qrcodeUtil.ByteToString(resp.Body())
		return response.Data(http.StatusOK, result)
	} else {
		return response.BadRequestMsg("app is invalid")
	}
}

func (s *Extension) PostExtension(ctx context.Context, domainUuid, userUuid string, extensionPost model.ExtensionPost) (int, any) {
	if ext, err := repository.ExtensionRepo.GetExtensionByExten(ctx, domainUuid, extensionPost.Extension); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if ext != nil {
		return response.BadRequestMsg("extension is existed")
	}
	domain, err := repository.DomainRepo.GetDomainById(ctx, domainUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if domain == nil {
		return response.BadRequestMsg("domain is invalid")
	}
	extensionUser := new(model.ExtensionUser)
	extension := model.Extension{
		Extension:             extensionPost.Extension,
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
		Enabled:               fmt.Sprintf("%t", extensionPost.Enabled),
		Password:              extensionPost.Password,
	}
	if len(extensionPost.UserUuid) > 0 {
		user, err := repository.UserCrmRepo.GetUserCrmById(ctx, extensionPost.UserUuid)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if user == nil {
			return response.BadRequestMsg("user_uuid is invalid")
		}
		extensionUser = &model.ExtensionUser{
			ExtensionUserUuid: uuid.NewString(),
			DomainUuid:        domainUuid,
			UserUuid:          extensionPost.UserUuid,
			ExtensionUuid:     extension.ExtensionUuid,
		}
	} else {
		extensionUser = nil
	}

	if err := repository.ExtensionRepo.InsertExtensionTransaction(ctx, &extension, extensionUser); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	result := map[string]any{
		"message":        "success",
		"extension_uuid": extension.ExtensionUuid,
		"extension":      extension.Extension,
		"enabled":        extension.Enabled == "true",
	}
	// if extensionPost.IsFollowMe {
	// 	if _, err := handleExtensionFollowMe(ctx, *domain, extension, extensionPost.FollowMe); err != nil {
	// 		log.Error(err)
	// 		return response.ServiceUnavailableMsg(err.Error())
	// 	}
	// 	result["follow_me"] = extensionPost.FollowMe
	// }
	if extensionPost.IsRingGroup && len(extensionPost.RingGroup.Main) > 0 {
		if err := handleExtensionRingGroup(ctx, *domain, extension.Extension, extensionPost.RingGroup); err != nil {
			log.Error(err)
		} else {
			freeswitch.ESLClient.ClearCacheDialplan(domain.DomainName)
			if err := freeswitch.ESLClient.SendToAllHost("reloadxml"); err != nil {
				log.Error(err)
			}
		}
	}
	freeswitch.ESLClient.ClearCacheOfDomain(domain.DomainName)
	go AddTransaction(domainUuid, userUuid, "extension", extension.ExtensionUuid, "add", "done", "", map[string]any{
		"extension":      extension,
		"extension_user": extensionUser,
	}, nil)
	return response.Created(result)
}

func (s *Extension) PutExtension(ctx context.Context, domainUuid, userUuid string, extensionUuid string, extensionPost model.ExtensionPost) (int, any) {
	extension, err := repository.ExtensionRepo.GetExtensionByIdOrExten(ctx, domainUuid, extensionUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if extension == nil {
		return response.BadRequestMsg("extension is invalid")
	}
	oldData := *extension
	domain, err := repository.DomainRepo.GetDomainById(ctx, domainUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if domain == nil {
		return response.BadRequestMsg("domain is invalid")
	}
	extensionUser := new(model.ExtensionUser)
	if extensionPost.Password != "deFauLt" {
		extension.Password = extensionPost.Password
	}
	extension.Enabled = fmt.Sprintf("%t", extensionPost.Enabled)
	if len(extensionPost.UserUuid) > 0 {
		if extTmp, err := repository.ExtensionRepo.GetExtensionByUserUuid(ctx, domainUuid, extensionPost.UserUuid); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if extTmp != nil && extTmp.Extension != extension.Extension {
			return response.BadRequestMsg("user is already mapping")
		}
		user, err := repository.UserCrmRepo.GetUserCrmById(ctx, extensionPost.UserUuid)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if user == nil {
			return response.BadRequestMsg("user_uuid is invalid")
		}
		extensionUser = &model.ExtensionUser{
			ExtensionUserUuid: uuid.NewString(),
			DomainUuid:        domainUuid,
			UserUuid:          extensionPost.UserUuid,
			ExtensionUuid:     extension.ExtensionUuid,
		}
	} else {
		extensionUser = nil
	}
	if err := repository.ExtensionRepo.UpdateExtensionTransaction(ctx, extension, extensionUser); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	result := map[string]any{
		"message":        "success",
		"extension_uuid": extension.ExtensionUuid,
		"extension":      extension.Extension,
		"enabled":        extension.Enabled == "true",
	}
	// if extensionPost.IsFollowMe {
	// 	if _, err := handleExtensionFollowMe(ctx, *domain, *extension, extensionPost.FollowMe); err != nil {
	// 		log.Error(err)
	// 		return response.ServiceUnavailableMsg(err.Error())
	// 	}
	// 	result["follow_me"] = extensionPost.FollowMe
	// } else if extensionPost.IsDeleteFollowMe && len(extension.FollowMeUuid) > 0 {
	// 	if err := handleRemoveFollowMe(ctx, extension.ExtensionUuid, extension.FollowMeUuid); err != nil {
	// 		log.Error(err)
	// 		return response.ServiceUnavailableMsg(err.Error())
	// 	}
	// }
	freeswitch.ESLClient.ClearCacheOfDomain(domain.DomainName)
	go AddTransaction(domainUuid, userUuid, "extension", extension.ExtensionUuid, "update", "done", "", map[string]any{
		"extension": oldData,
	}, map[string]any{
		"extension":      extension,
		"extension_user": extensionUser,
	})
	return response.OK(result)
}

func (s *Extension) DeleteExtension(ctx context.Context, domainUuid, userUuid string, extensionUuid string) (int, any) {
	extension, err := repository.ExtensionRepo.GetExtensionByIdOrExten(ctx, domainUuid, extensionUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if extension == nil {
		return response.BadRequestMsg("extension is invalid")
	}
	oldData := *extension
	domain, err := repository.DomainRepo.GetDomainById(ctx, domainUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if domain == nil {
		return response.BadRequestMsg("domain is invalid")
	}
	if err := repository.ExtensionRepo.DeleteExtensionTransaction(ctx, domainUuid, extensionUuid); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	freeswitch.ESLClient.ClearCacheOfDomain(domain.DomainName)
	go AddTransaction(domainUuid, userUuid, "extension", extension.ExtensionUuid, "delete", "done", "", map[string]any{
		"extension": oldData,
	}, nil)
	return response.OK(map[string]any{
		"message":        "success",
		"extension_uuid": extension.ExtensionUuid,
		"extension":      extension.Extension,
	})
}

func (s *Extension) PatchExtension(ctx context.Context, domainUuid, userUuid string, extensionUuid string, status bool) (int, any) {
	extension, err := repository.ExtensionRepo.GetExtensionInfoByExtensionUuid(ctx, domainUuid, extensionUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if extension == nil {
		return response.NotFoundMsg("extension is not existed")
	}
	if err := repository.ExtensionRepo.UpdateExtensionEnabled(ctx, domainUuid, extension.Extension, status); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(map[string]any{
		"message":        "success",
		"extension_uuid": extension.ExtensionUuid,
	})
}
func handleExtensionFollowMe(ctx context.Context, domain model.Domain, extension model.Extension, followMeConfig model.FollowMeConfig) (bool, error) {
	followMe := model.FollowMe{
		DomainUuid:         domain.DomainUuid,
		FollowMeEnabled:    true,
		DialString:         "",
		FollowMeIgnoreBusy: false,
	}
	if len(extension.FollowMeUuid) < 1 {
		// create
		followMe.FollowMeUuid = uuid.NewString()
	} else {
		followMe.FollowMeUuid = extension.FollowMeUuid
	}
	extension.DialDomain = domain.DomainName
	extension.DoNotDisturb = "false"
	extension.ForwardAllEnabled = "false"
	extension.ForwardBusyEnabled = "false"
	extension.ForwardNoAnswerEnabled = "false"
	extension.ForwardUserNotRegisteredEnabled = "false"
	extension.FollowMeUuid = followMe.FollowMeUuid
	followMeDestinationUuid := uuid.NewString()
	followMeDestination, err := repository.FollowMeRepo.GetFollowMeDestinationByFollowMeUuid(ctx, domain.DomainUuid, followMe.FollowMeUuid)
	if err != nil {
		log.Error(err)
		return false, err
	} else if followMeDestination != nil {
		followMeDestinationUuid = followMeDestination.FollowMeDestinationUuid
	}
	followMePrompt := 0
	if followMeConfig.Confirm {
		followMePrompt = 1
	}
	followMeDestination = &model.FollowMeDestination{
		DomainUuid:              domain.DomainUuid,
		FollowMeUuid:            followMe.FollowMeUuid,
		FollowMeDestination:     followMeConfig.Destination,
		FollowMeDestinationUuid: followMeDestinationUuid,
		FollowMeDelay:           followMeConfig.Delay,
		FollowMeTimeout:         followMeConfig.Timeout,
		FollowMePrompt:          followMePrompt,
		FollowMeOrder:           0,
	}
	variables := []string{}
	presenceId := extension.NumberAlias
	if util.CheckPatternNumeric(extension.Extension) {
		presenceId = extension.Extension
	}
	variables = append(variables, fmt.Sprintf("presence_id=%s@%s", presenceId, domain.DomainName))
	variables = append(variables, "origination_caller_id_number=${cond(${from_user_exists} == true ? : ${origination_caller_id_number})}")
	variables = append(variables, "effective_caller_id_number=${cond(${from_user_exists} == true ? : ${effective_caller_id_number})}")
	variables = append(variables, "origination_caller_id_name=${cond(${from_user_exists} == true ? : ${origination_caller_id_name})}")
	variables = append(variables, "effective_caller_id_name=${cond(${from_user_exists} == true ? : ${effective_caller_id_name})}")
	variables = append(variables, "sip_h_X-accountcode=${accountcode}")
	variables = append(variables, "fail_on_single_reject=USER_BUSY")
	if followMeConfig.Confirm {
		variables = append(variables, "group_confirm_key=exec")
		variables = append(variables, "group_confirm_file=lua confirm.lua")
		variables = append(variables, "confirm=true")
	}
	variables = append(variables, "instant_ringback=true")
	variables = append(variables, "ignore_early_media=true")
	variables = append(variables, fmt.Sprintf("domain_uuid=%s", domain.DomainUuid))
	variables = append(variables, fmt.Sprintf("sip_invite_domain=%s", domain.DomainName))
	variables = append(variables, fmt.Sprintf("domain_name=%s", domain.DomainName))
	variables = append(variables, fmt.Sprintf("extension_uuid=%s", extension.ExtensionUuid))
	variables = append(variables, fmt.Sprintf("leg_delay_start=%d", followMeConfig.Delay))
	variables = append(variables, fmt.Sprintf("originate_delay_start=%d", followMeConfig.Delay))
	variables = append(variables, fmt.Sprintf("sleep=%d", followMeConfig.Delay*1000))
	variables = append(variables, fmt.Sprintf("leg_timeout=%d", followMeConfig.Timeout))
	variables = append(variables, "is_follow_me_loopback=true")
	variablesStr := strings.Join(variables, "\\,export:")
	dialString := fmt.Sprintf("{ignore_early_media=true}loopback/export:%s\\,transfer:%s/%s/inline", variablesStr, followMeConfig.Destination, domain.DomainName)
	extension.DialString = dialString
	followMe.DialString = dialString
	if err := repository.ExtensionRepo.UpdateExtensionFollowMeTransaction(ctx, &extension, &followMe, followMeDestination); err != nil {
		log.Error(err)
		return false, err
	}
	return true, nil
}

func handleRemoveFollowMe(ctx context.Context, extensionUuid string, followMeUuid string) error {
	if err := repository.ExtensionRepo.DeleteExtensionFollowMeTransaction(ctx, extensionUuid, followMeUuid); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
