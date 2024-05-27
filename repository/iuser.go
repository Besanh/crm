package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IUser interface {
	GetUserByUsername(ctx context.Context, domainUuid string, userName string) (*model.User, error)
	GetUserByIdOrUsername(ctx context.Context, domainUuid string, value string) (*model.User, error)
	GetUsersInfo(ctx context.Context, domainUuid string, limit, offset int, filter model.UserFilter) (int, *[]model.UserView, error)
	GetUserById(ctx context.Context, domainUuid string, userUuid string) (*model.User, error)
	GetUserOfDomainByName(ctx context.Context, domainUuid string, userName string) (user *model.UserView, err error)
	InsertUser(ctx context.Context, user *model.User) error
	InsertUserTransaction(ctx context.Context, user *model.User, contact *model.VContact, groupUser *[]model.GroupUser, contactEmail *model.VContactEmail, roleGroup *model.RoleGroup) error
	GetAllUserUuidOfGroup(ctx context.Context, domainUuid, groupUuid string) ([]string, error)
	GetLiveUsersBackground(ctx context.Context) (*[]model.UserLive, error)
	GetContactById(ctx context.Context, domainUuid, contactUuid string) (*model.VContact, error)
	UpdateUserTransaction(ctx context.Context, user *model.User, contact *model.VContact, contactEmail *model.VContactEmail) error
	GetUserPasswordInfo(ctx context.Context, domainUuid, userUuid string) (*model.User, error)
	UpdateUserPassword(ctx context.Context, userUuid string, salt, passwordEncrypted string) error
	GetUserInfoById(ctx context.Context, userUuid string) (*model.UserView, error)
	GetUserViewById(ctx context.Context, domainUuid, userUuid string) (*model.UserView, error)
	GetUserViewByIdOrUsername(ctx context.Context, domainUuid, id string) (*model.UserView, error)
	GetUserByExtension(ctx context.Context, domainUuid, extension string) (*model.UserExtensionView, error)
	GetLiveUsers(ctx context.Context, domainUuid string, filter model.MonitorFilter) (*[]model.UserLiveView, error)
	GetUserByUsernameAndDomainName(ctx context.Context, domainName string, username string) (*model.UserAuth, error)
	GetUserByAPIKey(ctx context.Context, apiKey string) (*model.UserAuth, error)
	UpdateLoginTimeUser(ctx context.Context, domainUuid string, userUuid string) error
	GetUserDataById(ctx context.Context, userUuid string) (*model.UserData, error)
	GetUsersInfoOfGroupUsers(ctx context.Context, domainUuid string, userUuid, userlevel string, groupUuids []string, level []string, isExtension bool) (*[]model.UserInfoData, error)
	UpdateUserAndGroupTransaction(ctx context.Context, user *model.User, contact *model.VContact, groupUsers *[]model.GroupUser, contactEmail *model.VContactEmail) error
	UpdateUserActive(ctx context.Context, domainUuid string, userUuid, extensionUuid string, enable bool) error
	DeleteUserTransaction(ctx context.Context, domainUuid, userUuid, contactUuid string) error
	GetUserLiveByExtension(ctx context.Context, domainUuid, extension string) (*model.UserLive, error)
	GetUserCustomDatasByUserUuid(ctx context.Context, userUuid string) (*[]model.UserCustomData, error)
	InsertUserCustomData(ctx context.Context, userCustomDatas ...model.UserCustomData) error
	GetUserBasicInfoById(ctx context.Context, domainUuid, userId string) (*model.UserBasicInfo, error)
	SelectUserLive(ctx context.Context) (*[]model.UserLive, error)
}

var UserRepo IUser
