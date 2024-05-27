package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IContactCareerUser interface {
	InsertContactCareerUser(ctx context.Context, contactCareerUser model.ContactCareerUser) error
	GetContactCareerUserByContactCareerUuid(ctx context.Context, domainUuid, contactCareerUuid string) (*[]model.ContactCareerUser, error)
	DeleteContactCareerUserByContactCareerUuid(ctx context.Context, domainUuid, contactCareerUuid string) error
}

var ContactCareerUserRepo IContactCareerUser
