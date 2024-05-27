package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IContactGroupUser interface {
	InsertContactGroupUser(ctx context.Context, contactGroupUser model.ContactGroupUser) error
	GetContactGroupUserByContactGroupUuid(ctx context.Context, domainUuid, contactGroupUuid string) (*[]model.ContactGroupUser, error)
	DeleteContactGroupUserByContactGroupUuid(ctx context.Context, domainUuid, contactGroupUuid string) error
}

var ContactGroupUserRepo IContactGroupUser
