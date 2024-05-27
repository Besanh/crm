package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type ContactGroupUserRepo struct{}

func NewContactGroupUserRepo() repository.IContactGroupUser {
	repo := &ContactGroupUserRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()
	return repo
}

func (repo *ContactGroupUserRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ContactGroupUser)(nil)); err != nil {
		panic(err)
	}
}

func (repo *ContactGroupUserRepo) InitColumn() {

}

func (repo *ContactGroupUserRepo) InitIndex() {

}

func (repo *ContactGroupUserRepo) InsertContactGroupUser(ctx context.Context, ContactGroupUser model.ContactGroupUser) error {
	resp, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&ContactGroupUser).Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("insert contact group user failed")
	}

	return nil
}

func (repo *ContactGroupUserRepo) GetContactGroupUserByContactGroupUuid(ctx context.Context, domainUuid, contactGroupUuid string) (*[]model.ContactGroupUser, error) {
	result := new([]model.ContactGroupUser)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("domain_uuid = ?", domainUuid)
	query.Where("contact_group_uuid = ?", contactGroupUuid)

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, errors.New("contact group user not found")
	}

	return result, nil
}

func (repo *ContactGroupUserRepo) DeleteContactGroupUserByContactGroupUuid(ctx context.Context, domainUuid, contactGroupUuids string) error {
	query := repository.FusionSqlClient.GetDB().NewDelete().Model((*model.ContactGroupUser)(nil)).
		Where("domain_uuid = ?", domainUuid).
		Where("contact_group_uuid = ?", contactGroupUuids)

	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("delete contact group user failed")
	}

	return nil
}
