package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IContactTag interface {
	InsertContactTag(ctx context.Context, contactTag *model.ContactTag) error
	InsertContactTagTransaction(ctx context.Context, contactTag *model.ContactTag, contactTagUsers *[]model.ContactTagUser) error
	GetContactTagById(ctx context.Context, domainUuid, id string) (*model.ContactTag, error)
	GetContactTags(ctx context.Context, domainUuid string, limit, offset int, filter model.ContactTagFilter) (int, *[]model.ContactTag, error)
	PutContactTagTransaction(ctx context.Context, contactTag *model.ContactTag, contactTagUsers *[]model.ContactTagUser) error
	DeleteContactTag(ctx context.Context, domainUuid string, contactTag *model.ContactTag) error
	GetContacTagByTagName(ctx context.Context, domainUuid, tagName string) (*model.ContactTag, error)
}

var ContactTagRepo IContactTag
