package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IProfile interface {
	GetProfileInfoById(ctx context.Context, domainUuid, profileId string) (*model.ProfileView, error)
	GetProfileById(ctx context.Context, domainUuid, profileId string) (*model.Profile, error)
	GetProfilesInfo(ctx context.Context, domainUuid string, filter model.ProfileFilter, limit, offset int) (int, *[]model.ProfileView, error)
	GetProfileByPhoneNumber(ctx context.Context, domainUuid string, phoneNumber ...string) (*model.ProfileView, error)
	InsertProfileTransaction(ctx context.Context, profile *model.Profile, profilePhones []model.ProfilePhone, profileEmails []model.ProfileEmail, profileOwners []model.ProfileOwner, profileNotes []model.ProfileNote) error
	UpdateProfileTransaction(ctx context.Context, profileUuid string, profile *model.Profile, contactPhones []model.ProfilePhone, contactEmails []model.ProfileEmail, profileOwners []model.ProfileOwner) error
	UpdateProfile(ctx context.Context, profile *model.Profile) error
	DeleteProfileTransaction(ctx context.Context, profileUuid string) error
	UpdateProfileByField(ctx context.Context, domainUuid string, profile model.Profile) error
	GetManageProfiles(ctx context.Context, domainUuid string, filter model.ProfileFilter, limit, offset int) (int, *[]model.ProfileManageView, error)
	DeleteProfile(ctx context.Context, profile []model.Profile) error
	DeleteProfileWithTicketTransaction(ctx context.Context, domainUuid string, profiles []model.Profile, tickets []model.Ticket) error
}

var ProfileRepo IProfile
