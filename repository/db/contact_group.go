package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/uptrace/bun"
)

type ContactGroupRepo struct{}

func NewContactGroupRepo() repository.IContactGroup {
	repo := &ContactGroupRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *ContactGroupRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ContactGroup)(nil)); err != nil {
		panic(err)
	}
}

func (repo *ContactGroupRepo) InitColumn() {
}

func (repo *ContactGroupRepo) InitIndex() {

}

func (repo *ContactGroupRepo) InsertContactGroup(ctx context.Context, contactGroup *model.ContactGroup) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(contactGroup)
	resp, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected != 1 {
		return errors.New("insert contact group failed")
	}
	return nil
}

func (repo *ContactGroupRepo) InsertContactGroupTransaction(ctx context.Context, contactGroup *model.ContactGroup, contactGroupUsers *[]model.ContactGroupUser) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(contactGroup).Exec(ctx); err != nil {
			return err
		}
		contactGroupUsersExist, err := repository.ContactGroupUserRepo.GetContactGroupUserByContactGroupUuid(ctx, contactGroup.DomainUuid, contactGroup.ContactGroupUuid)
		if err != nil {
			return err
		} else if len(*contactGroupUsersExist) > 0 {
			if _, err := tx.NewDelete().Model(&model.ContactGroupUser{}).Where("contact_group_uuid = ?", contactGroup.ContactGroupUuid).Exec(ctx); err != nil {
				return err
			}
		}
		if len(*contactGroupUsers) > 0 {
			if _, err := tx.NewInsert().Model(contactGroupUsers).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *ContactGroupRepo) GetContactGroups(ctx context.Context, domainUuid string, limit, offset int, filter model.ContactGroupFilter) (int, *[]model.ContactGroup, error) {
	result := new([]model.ContactGroup)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Relation("ContactGroupUsers", func(q *bun.SelectQuery) *bun.SelectQuery {
			q.Relation("User")
			return q
		}).
		Where("domain_uuid = ?", domainUuid)
	if len(filter.GroupName) > 0 {
		query.Where("? ILIKE ?", bun.Ident("group_name"), "%"+filter.GroupName+"%")
	}
	if len(filter.GroupType) > 0 {
		query.Where("group_type = ?", filter.GroupType)
	}

	if filter.Status.Valid {
		query.Where("status = ?", filter.Status.Bool)
		if !filter.Status.Bool {
			query.WhereOr("status IS NULL")
		}
	}
	if len(filter.StartTime) > 0 {
		query.Where("created_at >= ?", filter.StartTime)
	}
	if len(filter.EndTime) > 0 {
		query.Where("created_at <= ?", filter.EndTime)
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	query.Order("created_at DESC")

	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return 0, result, err
	} else if err == sql.ErrNoRows {
		return 0, result, errors.New("contact group not found")
	}
	return total, result, nil
}

func (repo *ContactGroupRepo) GetContactGroupById(ctx context.Context, domainUuid, id string) (*model.ContactGroup, error) {
	result := new(model.ContactGroup)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Relation("ContactGroupUsers", func(q *bun.SelectQuery) *bun.SelectQuery {
			q.Relation("User")
			return q
		}).
		Where("cg.domain_uuid = ?", domainUuid).
		Where("cg.contact_group_uuid = ?", id)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *ContactGroupRepo) PutContactGroup(ctx context.Context, contactGroup *model.ContactGroup, contactGroupUsers *[]model.ContactGroupUser) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(contactGroup).WherePK("contact_group_uuid").Exec(ctx); err != nil {
			return err
		}
		contactGroupUsersExist, err := repository.ContactGroupUserRepo.GetContactGroupUserByContactGroupUuid(ctx, contactGroup.DomainUuid, contactGroup.ContactGroupUuid)
		if err != nil {
			return err
		} else if len(*contactGroupUsersExist) > 0 {
			if _, err := tx.NewDelete().Model(&model.ContactGroupUser{}).Where("contact_group_uuid = ?", contactGroup.ContactGroupUuid).Exec(ctx); err != nil {
				return err
			}
		}
		if len(*contactGroupUsers) > 0 {
			if _, err := tx.NewInsert().Model(contactGroupUsers).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *ContactGroupRepo) DeleteContactGroup(ctx context.Context, domainUuid string, contactGroup *model.ContactGroup) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().Model(contactGroup).WherePK("contact_group_uuid").Exec(ctx); err != nil {
			return err
		}
		contactGroupUsersExist, err := repository.ContactGroupUserRepo.GetContactGroupUserByContactGroupUuid(ctx, contactGroup.DomainUuid, contactGroup.ContactGroupUuid)
		if err != nil {
			return err
		} else if len(*contactGroupUsersExist) > 0 {
			if _, err := tx.NewDelete().Model(&model.ContactGroupUser{}).Where("contact_group_uuid = ?", contactGroup.ContactGroupUuid).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

/**
* Require group unique in domain
 */
func (repo *ContactGroupRepo) GetContactGroupByGroupName(ctx context.Context, domainUuid, groupName string) (*model.ContactGroup, error) {
	result := new(model.ContactGroup)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("domain_uuid = ?", domainUuid).
		Where("group_name = ?", groupName)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result, nil
}
