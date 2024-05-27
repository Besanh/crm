package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type ContactCareerUserRepo struct{}

func NewContactCareerUserRepo() repository.IContactCareerUser {
	repo := &ContactCareerUserRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()
	return repo
}

func (repo *ContactCareerUserRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ContactCareerUser)(nil)); err != nil {
		panic(err)
	}
}

func (repo *ContactCareerUserRepo) InitColumn() {

}

func (repo *ContactCareerUserRepo) InitIndex() {

}

func (repo *ContactCareerUserRepo) InsertContactCareerUser(ctx context.Context, ContactCareerUser model.ContactCareerUser) error {
	resp, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&ContactCareerUser).Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("insert contact career user failed")
	}

	return nil
}

func (repo *ContactCareerUserRepo) GetContactCareerUserByContactCareerUuid(ctx context.Context, domainUuid, contactCareerUuid string) (*[]model.ContactCareerUser, error) {
	result := new([]model.ContactCareerUser)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("domain_uuid = ?", domainUuid)
	query.Where("contact_career_uuid = ?", contactCareerUuid)

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, errors.New("contact career user not found")
	}

	return result, nil
}

func (repo *ContactCareerUserRepo) DeleteContactCareerUserByContactCareerUuid(ctx context.Context, domainUuid, contactCareerUuids string) error {
	query := repository.FusionSqlClient.GetDB().NewDelete().Model((*model.ContactCareerUser)(nil)).
		Where("domain_uuid = ?", domainUuid).
		Where("contact_career_uuid = ?", contactCareerUuids)

	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("delete contact career user failed")
	}

	return nil
}
