package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IContactTagUser interface {
	InsertContactTagUser(ctx context.Context, contactTagUser model.ContactTagUser) error
	GetContactTagUserByContactTagUuid(ctx context.Context, domainUuid, contactTagUuid string) (*[]model.ContactTagUser, error)
	DeleteContactTagUserByContactTagUuid(ctx context.Context, domainUuid, contactTagUuid string) error
}

var ContactTagUserRepo IContactTagUser
