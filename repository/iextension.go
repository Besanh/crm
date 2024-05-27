package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IExtension interface {
	GetExtensionByExten(ctx context.Context, domainUuid, exten string) (*model.ExtensionView, error)
	GetExtensionByUserUuid(ctx context.Context, domainUuid, userUuid string) (*model.ExtensionView, error)
	GetExtensionInfoByUserUuid(ctx context.Context, domainUuid, userUuid string) (*model.ExtensionInfoWithPassword, error)
	UpdateExtensionEnabled(ctx context.Context, domainUuid, extension string, isEnabled bool) error
	GetExtensionsOfUserUuids(ctx context.Context, domainUuid string, userUuids []string) (*[]model.ExtensionData, error)
	GetExtensionsInfo(ctx context.Context, domainUuid string, filter model.ExtensionFilter, limit, offset int) (*[]model.ExtensionInfo, int, error)
	GetExtensionInfoByExtensionUuid(ctx context.Context, domainUuid, extensionUuid string) (*model.ExtensionInfoWithPassword, error)
	InsertExtensionTransaction(ctx context.Context, extension *model.Extension, extensionUser *model.ExtensionUser) error
	GetExtensionByIdOrExten(ctx context.Context, domainUuid string, extensionUuid string) (*model.Extension, error)
	UpdateExtensionTransaction(ctx context.Context, extension *model.Extension, extensionUser *model.ExtensionUser) error
	DeleteExtensionTransaction(ctx context.Context, domainUuid string, extensionUuid string) error
	UpdateExtensionFollowMeTransaction(ctx context.Context, extension *model.Extension, followMe *model.FollowMe, followMeDestination *model.FollowMeDestination) error
	DeleteExtensionFollowMeTransaction(ctx context.Context, extensionUuid string, followMeUuid string) error
	UpdateExtensionRingGroupTransaction(ctx context.Context, dialplans []model.Dialplan, dialplanDetailsMap map[string][]model.DialplanDetail, ringGroups []model.RingGroup, ringGroupDestinations []model.RingGroupDestination) error
	FindExtensionDataOfExten(ctx context.Context, domainUuid string, exten string) (*model.ExtensionData, error)
}

var ExtensionRepo IExtension
