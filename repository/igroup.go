package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IGroup interface {
	SelectGroupsWithTotalUser(ctx context.Context, domainUuid string, filter model.GroupFilter, limit, offset int) (*[]model.GroupViewWithTotalUser, int, error)
	SelectGroupsOfDomain(ctx context.Context, domainUuid string, filter model.GroupFilter, limit, offset int) ([]model.Group, int, error)
	SelectGroupById(ctx context.Context, domainUuid string, groupUuid string) (group *model.Group, err error)
	SelectGroupInfoById(ctx context.Context, domainUuid string, groupUuid string) (*model.GroupView, error)
	SelectGroupByName(ctx context.Context, domainUuid string, groupName string) (group *model.Group, err error)
	InsertGroup(ctx context.Context, group *model.Group) error
	InsertGroupUser(ctx context.Context, groupUser *model.GroupUser) error
	UpdateGroup(ctx context.Context, group *model.Group) error
	// UpdateGroupFullTransaction(ctx context.Context, group *model.Group, campaignGroups *[]model.CampaignGroup, groupUsers *[]model.GroupUser, groupKpi *model.GroupKPI) error
	// DeleteGroupTransaction(ctx context.Context, domainUuid string, groupUuid string) error
}

var GroupRepo IGroup
