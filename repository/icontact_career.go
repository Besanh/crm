package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IContactCareer interface {
	InsertContactCareer(ctx context.Context, contactCareer *model.ContactCareer) error
	InsertContactCareerTransaction(ctx context.Context, contactCareer *model.ContactCareer, contactCareerUsers *[]model.ContactCareerUser) error
	GetContactCareerById(ctx context.Context, domainUuid, id string) (*model.ContactCareer, error)
	GetContactCareers(ctx context.Context, domainUuid string, limit, offset int, filter model.ContactCareerFilter) (int, *[]model.ContactCareer, error)
	PutContactCareer(ctx context.Context, contactCareer *model.ContactCareer, contactCareerUsers *[]model.ContactCareerUser) error
	DeleteContactCareer(ctx context.Context, domainUuid string, contactCareer *model.ContactCareer) error
	GetContacCareerByCareerName(ctx context.Context, domainUuid, careerName string) (*model.ContactCareer, error)
}

var ContactCareerRepo IContactCareer
