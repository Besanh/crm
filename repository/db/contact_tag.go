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

type ContactTagRepo struct{}

func NewContactTagRepo() repository.IContactTag {
	repo := &ContactTagRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *ContactTagRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ContactTag)(nil)); err != nil {
		panic(err)
	}
}

func (repo *ContactTagRepo) InitColumn() {
}

func (repo *ContactTagRepo) InitIndex() {

}

func (repo *ContactTagRepo) InsertContactTag(ctx context.Context, contactTag *model.ContactTag) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(contactTag)
	resp, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected != 1 {
		return errors.New("insert contact tag failed")
	}
	return nil
}

func (repo *ContactTagRepo) InsertContactTagTransaction(ctx context.Context, contactTag *model.ContactTag, contactTagUsers *[]model.ContactTagUser) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(contactTag).Exec(ctx); err != nil {
			return err
		}
		contactTagUsersExist, err := repository.ContactTagUserRepo.GetContactTagUserByContactTagUuid(ctx, contactTag.DomainUuid, contactTag.ContactTagUuid)
		if err != nil {
			return err
		} else if len(*contactTagUsersExist) > 0 {
			if _, err := tx.NewDelete().Model(&model.ContactTagUser{}).Where("contact_tag_uuid = ?", contactTag.ContactTagUuid).Exec(ctx); err != nil {
				return err
			}
		}
		if len(*contactTagUsers) > 0 {
			if _, err := tx.NewInsert().Model(contactTagUsers).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *ContactTagRepo) GetContactTags(ctx context.Context, domainUuid string, limit, offset int, filter model.ContactTagFilter) (int, *[]model.ContactTag, error) {
	result := new([]model.ContactTag)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Relation("ContactTagUsers", func(q *bun.SelectQuery) *bun.SelectQuery {
			q.Relation("User")
			return q
		}).
		Where("domain_uuid = ?", domainUuid)
	if len(filter.TagName) > 0 {
		query.Where("? ILIKE ?", bun.Ident("tag_name"), "%"+filter.TagName+"%")
	}

	if filter.Status.Valid {
		query.Where("status = ?", filter.Status.Bool)
		if !filter.Status.Bool {
			query.WhereOr("status IS NULL")
		}
	}
	if len(filter.TagType) > 0 {
		query.Where("? = ANY(limited_function)", filter.TagType)
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
		return 0, result, errors.New("contact Tag not found")
	}
	return total, result, nil
}

func (repo *ContactTagRepo) GetContactTagById(ctx context.Context, domainUuid, id string) (*model.ContactTag, error) {
	result := new(model.ContactTag)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Relation("ContactTagUsers", func(q *bun.SelectQuery) *bun.SelectQuery {
			q.Relation("User")
			return q
		}).
		Where("ct.domain_uuid = ?", domainUuid).
		Where("ct.contact_tag_uuid = ?", id)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *ContactTagRepo) PutContactTagTransaction(ctx context.Context, contactTag *model.ContactTag, contactTagUsers *[]model.ContactTagUser) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(contactTag).WherePK("contact_tag_uuid").Exec(ctx); err != nil {
			return err
		}
		contactTagUsersExist, err := repository.ContactTagUserRepo.GetContactTagUserByContactTagUuid(ctx, contactTag.DomainUuid, contactTag.ContactTagUuid)
		if err != nil {
			return err
		} else if len(*contactTagUsersExist) > 0 {
			if _, err := tx.NewDelete().Model(&model.ContactTagUser{}).Where("contact_tag_uuid = ?", contactTag.ContactTagUuid).Exec(ctx); err != nil {
				return err
			}
		}
		if len(*contactTagUsers) > 0 {
			if _, err := tx.NewInsert().Model(contactTagUsers).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *ContactTagRepo) DeleteContactTag(ctx context.Context, domainUuid string, contactTag *model.ContactTag) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().Model(contactTag).WherePK("contact_tag_uuid").Exec(ctx); err != nil {
			return err
		}
		contactTagUsersExist, err := repository.ContactTagUserRepo.GetContactTagUserByContactTagUuid(ctx, contactTag.DomainUuid, contactTag.ContactTagUuid)
		if err != nil {
			return err
		} else if len(*contactTagUsersExist) > 0 {
			if _, err := tx.NewDelete().Model(&model.ContactTagUser{}).Where("contact_tag_uuid = ?", contactTag.ContactTagUuid).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

/**
* Require tag unique in domain
 */
func (repo *ContactTagRepo) GetContacTagByTagName(ctx context.Context, domainUuid, tagName string) (*model.ContactTag, error) {
	result := new(model.ContactTag)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("domain_uuid = ?", domainUuid).
		Where("tag_name = ?", tagName)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result, nil
}
