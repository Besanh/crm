package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IContact interface {
	// Get
	GetContactsInfo(ctx context.Context, domainUuid string, filter model.ContactFilter, limit, offset int) (int, *[]model.ContactView, error)
	GetContactInfo(ctx context.Context, domainUuid string, filter model.ContactFilter) (*model.ContactView, error)
	GetContactInfoById(ctx context.Context, domainUuid string, contactId string) (*model.ContactView, error)
	GetContactById(ctx context.Context, domainUuid string, contactId string) (*model.Contact, error)
	GetContactByPhoneNumber(ctx context.Context, domainUuid string, phoneNumber ...string) (*model.Contact, error)
	GetContactNotes(ctx context.Context, domainUuid string, contactId string, limit, offset int) (*[]model.ContactNoteData, int, error)
	GetContactInfoByCallId(ctx context.Context, domainUuid string, callId string) (*model.ContactView, error)

	// Put
	UpdateContact(ctx context.Context, contact *model.Contact) error
	UpdateContactTransaction(ctx context.Context, contactUuid string, contact *model.Contact, contactPhones []model.ContactPhone, contactEmails []model.ContactEmail, contactOwners []model.ProfileOwner, contactToTags []model.ContactToTag, contactToGroups []model.ContactToGroup, contactToCareers []model.ContactToCareer) error

	// Insert
	InsertContact(ctx context.Context, contact *model.Contact) error
	InsertContactPhone(ctx context.Context, contact *model.ContactPhone) error
	InsertContactEmail(ctx context.Context, contact *model.ContactEmail) error
	InsertListContact(ctx context.Context, contact *[]model.Contact) error
	InsertContactTransaction(ctx context.Context, contact *model.Contact, contactPhones []model.ContactPhone, contactEmails []model.ContactEmail, contactOwners []model.ProfileOwner, contactNotes []model.ContactNote, contactToTags []model.ContactToTag, contactToGroups []model.ContactToGroup, contactToCareers []model.ContactToCareer) error
	InsertContactsTransaction(ctx context.Context, contacts []model.Contact, contactPhones []model.ContactPhone, contactEmails []model.ContactEmail) error
	InsertContactNote(ctx context.Context, entry *model.ContactNote) error

	// Delete
	DeleteContactTransaction(ctx context.Context, contactUuid string) error
}

var ContactRepo IContact
