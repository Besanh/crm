package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type ICampaign interface {
	GetCampaignsOptionView(ctx context.Context, domainUuid string, filter model.CampaignFilter, limit, offset int) (*[]model.CampaignOptionView, int, error)
	GetCampaigns(ctx context.Context, domainUuid string, filter model.GeneralFilter, limit, offset int) (*[]model.Campaign, int, error)
	GetCampaignsInfo(ctx context.Context, domainUuid string, filter model.CampaignFilter, limit, offset int) (*[]model.CampaignView, int, error)
	GetCampaignInfoById(ctx context.Context, domainUuid, campaignUuid string) (*model.CampaignView, error)
	GetCampaignById(ctx context.Context, domainUuid string, campaignUuid string) (campaign *model.Campaign, err error)
	GetCampaignByName(ctx context.Context, domainUuid string, campaignName string) (*model.Campaign, error)
	GetCampaignsActive(ctx context.Context) (*[]model.CampaignView, error)
	GetCampaignsActiveContainHopper(ctx context.Context) (*[]model.CampaignView, error)
	InsertCampaign(ctx context.Context, campaign *model.Campaign) error
	InsertCampaignAutodialerTransaction(ctx context.Context, campaign *model.Campaign, callCenterQueue *model.CallCenterQueue, dialplan *model.Dialplan) error
	InsertCampaignUsers(ctx context.Context, campaignUuid string, campaignUses *[]model.CampaignUser) error
	InsertCampaignGroups(ctx context.Context, campaignUuid string, campaignGroups *[]model.CampaignGroup) error
	GetCampaignUserByUserUuid(ctx context.Context, campaignUuid string, userUuid string) (*model.CampaignUser, error)
	GetCampaignUsersByCampaignUuid(ctx context.Context, campaignUuid string) (*[]model.CampaignUser, error)
	GetCampaignsStatusInfo(ctx context.Context, domainUuid string, limit, offset int) (*[]model.CampaignStatusView, int, error)
	GetInboundCampaignByQueueExtension(ctx context.Context, domainUuid, queueExtension string) (*model.Campaign, error)
	UpdateCampaignTransaction(ctx context.Context, campaign *model.Campaign, campaignUsers *[]model.CampaignUser, statuses []string) error
	UpdateCampaign(ctx context.Context, domainUuid string, campaign *model.Campaign) error
	DeleteCampaign(ctx context.Context, domainUuid string, campaignUuid string) error
	DeleteCampaignUser(ctx context.Context, domainUuid string, campaignUuid string) error
	DeleteCampaignGroup(ctx context.Context, domainUuid string, campaignUuid string) error
	InsertCampaignSchedules(ctx context.Context, campaignUuid string, schedules ...model.CampaignSchedule) error
	GetCampaignUuidsOfUsers(ctx context.Context, domainUuid string, userUuids []string) ([]string, error)
	AppendCampaignGroups(ctx context.Context, campaignUuid string, campaignGroups ...model.CampaignGroup) error
	AppendCampaignUsers(ctx context.Context, campaignUuid string, campaignUsers ...model.CampaignUser) error
	GetCampaignsByCarrierId(ctx context.Context, domainUuid, carrierUuid string) (*[]model.Campaign, error)
	GetCampaignCustomDatasByCampaignUuid(ctx context.Context, campaignUuid string) (*[]model.CampaignCustomData, error)
	InsertCampaignCustomData(ctx context.Context, campaignCustomDatas ...model.CampaignCustomData) error
}

var CampaignRepo ICampaign
