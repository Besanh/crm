package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	goauth "contactcenter-api/middleware/auth/goauth"
	"contactcenter-api/repository"
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/jinzhu/copier"
)

type (
	IAuth interface {
		GenerateTokenByApiKey(ctx context.Context, apiKey string) (int, any)
		SigninContactCenter(ctx context.Context, domain, username, password string) (int, any)
	}
	Auth struct {
	}
)

var AuthServiceGlobal IAuth

func NewAuth() IAuth {
	return &Auth{}
}

func (service *Auth) SigninContactCenter(ctx context.Context, domainName, username, password string) (int, any) {
	user, err := repository.UserRepo.GetUserByUsernameAndDomainName(ctx, domainName, username)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if user == nil {
		return response.ServiceUnavailableMsg("username or password is not valid")
	}
	passCurrent := user.Password
	hash := md5.New()
	salt := user.Salt
	tmp := salt + password
	_, err = hash.Write([]byte(tmp))
	if err != nil {
		log.Error(err)
		return response.BadRequestMsg("username or password is not valid")
	}
	passEncrypted := string(hex.EncodeToString(hash.Sum(nil)))
	if passEncrypted != passCurrent {
		return response.BadRequestMsg("username or password is not valid")
	}
	// Update time login
	userInfo, err := repository.UserCrmRepo.GetUserCrmById(ctx, user.UserUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if userInfo == nil {
		return response.BadRequestMsg("user info error")
	}
	if err := repository.UserRepo.UpdateLoginTimeUser(ctx, user.DomainUuid, user.UserUuid); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	res := map[string]any{
		"user_uuid":       user.UserUuid,
		"enable":          user.UserEnabled == "true",
		"domain_uuid":     user.DomainUuid,
		"level":           user.Level,
		"api_key":         user.ApiKey,
		"username":        user.Username,
		"fullname":        userInfo.FirstName + " " + userInfo.MiddleName + " " + userInfo.LastName,
		"enable_webrtc":   userInfo.EnableWebrtc,
		"last_login_date": time.Now(),
	}
	clientAuth := goauth.AuthClient{
		ClienId: user.ApiKey,
		UserId:  user.UserUuid,
		UserData: map[string]any{
			"domain_uuid": user.DomainUuid,
			"domain_name": user.DomainName,
			"username":    user.Username,
			"level":       user.Level,
		},
	}
	client, err := goauth.GoAuthClient.ClientCredential(ctx, clientAuth, false)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	token := map[string]any{
		"client_id":     client.ClienId,
		"user_uuid":     client.UserId,
		"token":         client.Token,
		"refresh_token": client.RefreshToken,
		"expired_in":    client.ExpiredIn,
		"token_type":    client.TokenType,
		"domain_uuid":   user.DomainUuid,
		"unit_uuid":     userInfo.UnitUuid,
	}
	res["token"] = token
	if extension, err := repository.ExtensionRepo.GetExtensionByUserUuid(ctx, user.DomainUuid, user.UserUuid); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if extension != nil {
		res["extension"] = extension.Extension
	}
	if len(userInfo.RoleUuid) > 0 {
		if roleGroup, err := repository.RoleGroupRepo.GetRoleGroupById(ctx, user.DomainUuid, userInfo.RoleUuid); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if roleGroup != nil {
			permissionMainView := model.PermissionView{}
			permissionAdvanceView := model.PermissionAdvanceView{}
			copier.Copy(&permissionMainView, roleGroup.PermissionMain)
			copier.Copy(&permissionAdvanceView, roleGroup.PermissionAdvance)
			res["privilege"] = permissionMainView
			res["privilege_advance"] = permissionAdvanceView
		}
	}
	if domainConfig, err := repository.DomainRepo.GetDomainConfigById(ctx, userInfo.DomainUuid); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if domainConfig != nil {
		res["domain_config"] = domainConfig
	}
	return response.Data(http.StatusOK, res)
}

func (service *Auth) GenerateTokenByApiKey(ctx context.Context, apiKey string) (int, any) {
	return 0, nil
	// user, err := repository.UserRepo.GetUserByAPIKey(ctx, apiKey)
	// if err != nil {
	// 	log.Error(err)
	// 	return response.ServiceUnavailableMsg(err.Error())
	// } else if user == nil {
	// 	log.Error(err)
	// 	return response.NotFound()
	// }

	// clientAuth := goauth.AuthClient{
	// 	ClienId: apiKey,
	// 	UserId:  user.UserUuid,
	// 	UserData: map[string]any{
	// 		"domain_uuid": user.DomainUuid,
	// 		"domain_name": user.DomainName,
	// 		"username":    user.Username,
	// 		"level":       user.Level,
	// 	},
	// }
	// client, err := goauth.GoAuthClient.ClientCredential(ctx, clientAuth, false)
	// if err != nil {
	// 	log.Error(err)
	// 	return response.ServiceUnavailableMsg(err.Error())
	// }
	// token := map[string]any{
	// 	"client_id":     client.ClienId,
	// 	"user_id":       client.UserId,
	// 	"token":         client.Token,
	// 	"refresh_token": client.RefreshToken,
	// 	"expired_in":    client.ExpiredIn,
	// 	"token_type":    client.TokenType,
	// 	"domain_uuid":   user.DomainUuid,
	// }
	// return response.OK(token)
}
