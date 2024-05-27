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

type ContactCareerRepo struct{}

func NewContactCareerRepo() repository.IContactCareer {
	repo := &ContactCareerRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *ContactCareerRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ContactCareer)(nil)); err != nil {
		panic(err)
	}
}

func (repo *ContactCareerRepo) InitColumn() {
}

func (repo *ContactCareerRepo) InitIndex() {

}

func (repo *ContactCareerRepo) InsertContactCareer(ctx context.Context, contactCareer *model.ContactCareer) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(contactCareer)
	resp, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected != 1 {
		return errors.New("insert contact career failed")
	}
	return nil
}

func (repo *ContactCareerRepo) InsertContactCareerTransaction(ctx context.Context, contactCareer *model.ContactCareer, contactCareerUsers *[]model.ContactCareerUser) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(contactCareer).Exec(ctx); err != nil {
			return err
		}
		contactCareerUsersExist, err := repository.ContactCareerUserRepo.GetContactCareerUserByContactCareerUuid(ctx, contactCareer.DomainUuid, contactCareer.ContactCareerUuid)
		if err != nil {
			return err
		} else if len(*contactCareerUsersExist) > 0 {
			if _, err := tx.NewDelete().Model(&model.ContactCareerUser{}).Where("contact_career_uuid = ?", contactCareer.ContactCareerUuid).Exec(ctx); err != nil {
				return err
			}
		}
		if len(*contactCareerUsers) > 0 {
			if _, err := tx.NewInsert().Model(contactCareerUsers).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *ContactCareerRepo) GetContactCareers(ctx context.Context, domainUuid string, limit, offset int, filter model.ContactCareerFilter) (int, *[]model.ContactCareer, error) {
	result := new([]model.ContactCareer)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Relation("ContactcareerUsers", func(q *bun.SelectQuery) *bun.SelectQuery {
			q.Relation("User")
			return q
		}).
		Where("domain_uuid = ?", domainUuid)
	if len(filter.CareerName) > 0 {
		query.Where("? ILIKE ?", bun.Ident("career_name"), "%"+filter.CareerName+"%")
	}
	if len(filter.CareerType) > 0 {
		query.Where("career_type = ?", filter.CareerType)
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
		return 0, result, errors.New("contact career not found")
	}
	return total, result, nil
}

func (repo *ContactCareerRepo) GetContactCareerById(ctx context.Context, domainUuid, id string) (*model.ContactCareer, error) {
	result := new(model.ContactCareer)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Relation("ContactcareerUsers", func(q *bun.SelectQuery) *bun.SelectQuery {
			q.Relation("User")
			return q
		}).
		Where("cc.domain_uuid = ?", domainUuid).
		Where("cc.contact_career_uuid = ?", id)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *ContactCareerRepo) PutContactCareer(ctx context.Context, contactCareer *model.ContactCareer, contactCareerUsers *[]model.ContactCareerUser) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(contactCareer).
			WherePK().Exec(ctx); err != nil {
			return err
		}
		contactCareerUsersExist, err := repository.ContactCareerUserRepo.GetContactCareerUserByContactCareerUuid(ctx, contactCareer.DomainUuid, contactCareer.ContactCareerUuid)
		if err != nil {
			return err
		} else if len(*contactCareerUsersExist) > 0 {
			if _, err := tx.NewDelete().Model(&model.ContactCareerUser{}).Where("contact_career_uuid = ?", contactCareer.ContactCareerUuid).Exec(ctx); err != nil {
				return err
			}
		}
		if len(*contactCareerUsers) > 0 {
			if _, err := tx.NewInsert().Model(contactCareerUsers).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo *ContactCareerRepo) DeleteContactCareer(ctx context.Context, domainUuid string, contactCareer *model.ContactCareer) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().Model(contactCareer).WherePK("contact_career_uuid").Exec(ctx); err != nil {
			return err
		}
		contactCareerUsersExist, err := repository.ContactCareerUserRepo.GetContactCareerUserByContactCareerUuid(ctx, contactCareer.DomainUuid, contactCareer.ContactCareerUuid)
		if err != nil {
			return err
		} else if len(*contactCareerUsersExist) > 0 {
			if _, err := tx.NewDelete().Model(&model.ContactCareerUser{}).Where("contact_career_uuid = ?", contactCareer.ContactCareerUuid).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

/**
* Require group unique in domain
 */
func (repo *ContactCareerRepo) GetContacCareerByCareerName(ctx context.Context, domainUuid, careerName string) (*model.ContactCareer, error) {
	result := new(model.ContactCareer)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("domain_uuid = ?", domainUuid).
		Where("career_name = ?", careerName)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result, nil
}
