package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"time"

	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type (
	Contact struct {
	}
)

func NewContact() repository.IContact {
	repo := &Contact{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *Contact) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// repository.FusionSqlClient.GetDB().RegisterModel((*model.GroupUser)(nil))

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.Contact)(nil)); err != nil {
		panic(err)
	}
	// if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ContactChannel)(nil)); err != nil {
	// 	panic(err)
	// }
	// if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ContactPhone)(nil)); err != nil {
	// 	panic(err)
	// }
	// if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ContactEmail)(nil)); err != nil {
	// 	panic(err)
	// }
	// if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ContactOwner)(nil)); err != nil {
	// 	panic(err)
	// }
	// if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ContactNote)(nil)); err != nil {
	// 	panic(err)
	// }
}

func (repo *Contact) InitColumn() {

}

func (repo *Contact) InitIndex() {

}

func (repo *Contact) GetContactsInfo(ctx context.Context, domainUuid string, filter model.ContactFilter, limit, offset int) (int, *[]model.ContactView, error) {
	contacts := new([]model.ContactView)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(contacts).
		Relation("Profiles").
		Relation("Profiles.Emails").
		Relation("Profiles.Phones").
		Relation("Profiles.UserOwners").
		Relation("Profiles.ListRelatedProfile").
		Where("c.domain_uuid = ?", domainUuid).
		Order("c.created_at DESC")

	if !filter.StartTime.IsZero() {
		query = query.Where("c.created_at >= ?", filter.StartTime)
	}
	if !filter.EndTime.IsZero() {
		query = query.Where("c.created_at <= ?", filter.EndTime)
	}
	if len(filter.ContactType) > 0 {
		query = query.Where("c.contact_type = ?", filter.ContactType)
	}
	if len(filter.ContactName) > 0 {
		query = query.Where("c.contact_name LIKE ?", "%"+filter.ContactName+"%")
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, contacts, nil
	}
	return total, contacts, err
}

func (repo *Contact) GetContactById(ctx context.Context, domainUuid string, contactId string) (*model.Contact, error) {
	contact := new(model.Contact)
	err := repository.FusionSqlClient.GetDB().NewSelect().Model(contact).
		Where("domain_uuid = ?", domainUuid).
		Where("contact_uuid = ?", contactId).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return contact, err
}

func (repo *Contact) GetContactInfoById(ctx context.Context, domainUuid string, contactId string) (*model.ContactView, error) {
	contact := new(model.ContactView)
	err := repository.FusionSqlClient.GetDB().NewSelect().Model(contact).
		Relation("Profiles").
		Relation("Profiles.Emails").
		Relation("Profiles.Phones").
		Relation("Profiles.UserOwners").
		Relation("Profiles.ListRelatedProfile").
		Where("c.domain_uuid = ?", domainUuid).
		Where("c.contact_uuid = ?", contactId).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return contact, err
}

func (repo *Contact) UpdateContact(ctx context.Context, contact *model.Contact) error {
	query := repository.FusionSqlClient.GetDB().NewUpdate().Model(contact).
		Where("contact_uuid = ?", contact.ContactUuid)
	_, err := query.Exec(ctx)
	return err
}

func (repo *Contact) InsertContact(ctx context.Context, contact *model.Contact) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(contact)
	_, err := query.Exec(ctx)
	return err
}

func (repo *Contact) InsertListContact(ctx context.Context, contact *[]model.Contact) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(contact)
	_, err := query.Exec(ctx)
	return err
}

func (repo *Contact) GetContactByPhoneNumber(ctx context.Context, domainUuid string, phoneNumber ...string) (*model.Contact, error) {
	contact := new(model.Contact)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(contact).
		// Join("LEFT JOIN contact_phone cp ON c.contact_uuid = cp.contact_uuid").
		Where("c.domain_uuid = ?", domainUuid).
		// Where("cp.data IN (?)", bun.In(phoneNumber)).
		Where("c.phone_number IN (?)", bun.In(phoneNumber)).
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return contact, err
}

func (repo *Contact) GetContactInfo(ctx context.Context, domainUuid string, filter model.ContactFilter) (*model.ContactView, error) {
	contact := new(model.ContactView)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(contact).
		Relation("LocationProvince").
		Relation("LocationDistrict").
		Relation("LocationWard").
		Relation("Emails").
		Relation("Phones").
		Relation("UserOwners").
		Join("LEFT JOIN contact_phone cp ON c.contact_uuid = cp.contact_uuid").
		Where("c.domain_uuid = ?", domainUuid)
	if len(filter.ContactType) > 0 {
		query = query.Where("c.contact_type = ?", filter.ContactType)
	}
	if len(filter.ContactName) > 0 {
		query = query.Where("c.contact_name LIKE ?", "%"+filter.ContactName+"%")
	}
	query.Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return contact, err
}

func (repo *Contact) InsertContactTransaction(ctx context.Context, contact *model.Contact, contactPhones []model.ContactPhone, contactEmails []model.ContactEmail, contactOwners []model.ProfileOwner, contactNotes []model.ContactNote, contactToTags []model.ContactToTag, contactToGroups []model.ContactToGroup, contactToCareers []model.ContactToCareer) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(contact).Exec(ctx); err != nil {
			return err
		}
		if len(contactPhones) > 0 {
			if _, err := tx.NewInsert().Model(&contactPhones).Exec(ctx); err != nil {
				return err
			}
		}
		if len(contactEmails) > 0 {
			if _, err := tx.NewInsert().Model(&contactEmails).Exec(ctx); err != nil {
				return err
			}
		}
		if len(contactOwners) > 0 {
			if _, err := tx.NewInsert().Model(&contactOwners).Exec(ctx); err != nil {
				return err
			}
		}
		if len(contactNotes) > 0 {
			if _, err := tx.NewInsert().Model(&contactNotes).Exec(ctx); err != nil {
				return err
			}
		}
		if len(contactToTags) > 0 {
			if _, err := tx.NewInsert().Model(&contactToTags).Exec(ctx); err != nil {
				return err
			}
		}
		if len(contactToGroups) > 0 {
			if _, err := tx.NewInsert().Model(&contactToGroups).Exec(ctx); err != nil {
				return err
			}
		}
		if len(contactToCareers) > 0 {
			if _, err := tx.NewInsert().Model(&contactToCareers).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *Contact) UpdateContactTransaction(ctx context.Context, contactUuid string, contact *model.Contact, contactPhones []model.ContactPhone, contactEmails []model.ContactEmail, contactOwners []model.ProfileOwner, contactToTags []model.ContactToTag, contactToGroups []model.ContactToGroup, contactToCareers []model.ContactToCareer) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(contact).Where("contact_uuid = ?", contactUuid).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.ContactPhone)(nil)).Where("contact_uuid = ?", contactUuid).Exec(ctx); err != nil {
			return err
		}
		if len(contactPhones) > 0 {
			if _, err := tx.NewInsert().Model(&contactPhones).Exec(ctx); err != nil {
				return err
			}
		}
		if _, err := tx.NewDelete().Model((*model.ContactEmail)(nil)).Where("contact_uuid = ?", contactUuid).Exec(ctx); err != nil {
			return err
		}
		if len(contactEmails) > 0 {
			if _, err := tx.NewInsert().Model(&contactEmails).Exec(ctx); err != nil {
				return err
			}
		}
		if _, err := tx.NewDelete().Model((*model.ProfileOwner)(nil)).Where("contact_uuid = ?", contactUuid).Exec(ctx); err != nil {
			return err
		}
		if len(contactOwners) > 0 {
			if _, err := tx.NewInsert().Model(&contactOwners).Exec(ctx); err != nil {
				return err
			}
		}
		if _, err := tx.NewDelete().Model((*model.ContactToTag)(nil)).Where("contact_uuid = ?", contactUuid).Exec(ctx); err != nil {
			return err
		}
		if len(contactToTags) > 0 {
			if _, err := tx.NewInsert().Model(&contactToTags).Exec(ctx); err != nil {
				return err
			}
		}
		if _, err := tx.NewDelete().Model((*model.ContactToGroup)(nil)).Where("contact_uuid = ?", contactUuid).Exec(ctx); err != nil {
			return err
		}
		if len(contactToGroups) > 0 {
			if _, err := tx.NewInsert().Model(&contactToGroups).Exec(ctx); err != nil {
				return err
			}
		}
		if _, err := tx.NewDelete().Model((*model.ContactToCareer)(nil)).Where("contact_uuid = ?", contactUuid).Exec(ctx); err != nil {
			return err
		}
		if len(contactToCareers) > 0 {
			if _, err := tx.NewInsert().Model(&contactToCareers).Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *Contact) GetContactNotes(ctx context.Context, domainUuid string, contactId string, limit, offset int) (*[]model.ContactNoteData, int, error) {
	notes := new([]model.ContactNoteData)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(notes).
		Column("cn.*").
		Where("cn.domain_uuid = ?", domainUuid).
		Where("cn.contact_uuid = ?", contactId).
		Order("created_at DESC")
	subQuery := repository.FusionSqlClient.GetDB().NewSelect().TableExpr("v_users u").
		Column("u.username").
		ColumnExpr("c.contact_name_given as first_name, c.contact_name_middle as middle_name, c.contact_name_family as last_name").
		Join("LEFT JOIN v_contacts c ON c.contact_uuid = u.contact_uuid").
		Where("u.user_uuid = cn.user_uuid")
	query = query.ColumnExpr("tmpu.*").Join("LEFT JOIN lateral (?) tmpu ON true", subQuery)
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return nil, 0, nil
	}
	return notes, total, err
}

func (repo *Contact) InsertContactNote(ctx context.Context, entry *model.ContactNote) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(entry)
	_, err := query.Exec(ctx)
	return err
}

func (repo *Contact) InsertContactsTransaction(ctx context.Context, contacts []model.Contact, contactPhones []model.ContactPhone, contactEmails []model.ContactEmail) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if len(contacts) > 0 {
			if _, err := tx.NewInsert().Model(&contacts).Exec(ctx); err != nil {
				return err
			}
		}
		if len(contactPhones) > 0 {
			if _, err := tx.NewInsert().Model(&contactPhones).Exec(ctx); err != nil {
				return err
			}
		}
		if len(contactEmails) > 0 {
			if _, err := tx.NewInsert().Model(&contactEmails).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *Contact) DeleteContactTransaction(ctx context.Context, contactUuid string) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().Model((*model.Contact)(nil)).Where("contact_uuid = ?", contactUuid).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.ContactPhone)(nil)).Where("contact_uuid = ?", contactUuid).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.ContactEmail)(nil)).Where("contact_uuid = ?", contactUuid).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.ContactOwner)(nil)).Where("contact_uuid = ?", contactUuid).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.ContactNote)(nil)).Where("contact_uuid = ?", contactUuid).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.ContactToTag)(nil)).Where("contact_uuid = ?", contactUuid).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.ContactToGroup)(nil)).Where("contact_uuid = ?", contactUuid).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.ContactToCareer)(nil)).Where("contact_uuid = ?", contactUuid).Exec(ctx); err != nil {
			return err
		}
		return nil
	})
}

func (repo *Contact) GetContactInfoByCallId(ctx context.Context, domainUuid string, callId string) (*model.ContactView, error) {
	result := new(model.ContactView)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("domain_uuid = ?", domainUuid).
		Where("call_id = ?", callId)

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}
