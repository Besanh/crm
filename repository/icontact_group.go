package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IContactGroup interface {
	InsertContactGroup(ctx context.Context, contactGroup *model.ContactGroup) error
	InsertContactGroupTransaction(ctx context.Context, contactGroup *model.ContactGroup, contactGroupUsers *[]model.ContactGroupUser) error
	GetContactGroupById(ctx context.Context, domainUuid, id string) (*model.ContactGroup, error)
	GetContactGroups(ctx context.Context, domainUuid string, limit, offset int, filter model.ContactGroupFilter) (int, *[]model.ContactGroup, error)
	PutContactGroup(ctx context.Context, contactGroup *model.ContactGroup, contactGroupUsers *[]model.ContactGroupUser) error
	DeleteContactGroup(ctx context.Context, domainUuid string, contactGroup *model.ContactGroup) error
	GetContactGroupByGroupName(ctx context.Context, domainUuid, groupName string) (*model.ContactGroup, error)
}

var ContactGroupRepo IContactGroup
