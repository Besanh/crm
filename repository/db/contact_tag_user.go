package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type ContactTagUserRepo struct{}

func NewContactTagUserRepo() repository.IContactTagUser {
	repo := &ContactTagUserRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()
	return repo
}

func (repo *ContactTagUserRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ContactTagUser)(nil)); err != nil {
		panic(err)
	}
}

func (repo *ContactTagUserRepo) InitColumn() {

}

func (repo *ContactTagUserRepo) InitIndex() {

}

func (repo *ContactTagUserRepo) InsertContactTagUser(ctx context.Context, ContactTagUser model.ContactTagUser) error {
	resp, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&ContactTagUser).Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("insert contact tag user failed")
	}

	return nil
}

func (repo *ContactTagUserRepo) GetContactTagUserByContactTagUuid(ctx context.Context, domainUuid, contactTagUuid string) (*[]model.ContactTagUser, error) {
	result := new([]model.ContactTagUser)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("domain_uuid = ?", domainUuid)
	query.Where("contact_tag_uuid = ?", contactTagUuid)

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, errors.New("contact tag user not found")
	}

	return result, nil
}

func (repo *ContactTagUserRepo) DeleteContactTagUserByContactTagUuid(ctx context.Context, domainUuid, contactTagUuids string) error {
	query := repository.FusionSqlClient.GetDB().NewDelete().Model((*model.ContactTagUser)(nil)).
		Where("domain_uuid = ?", domainUuid).
		Where("contact_tag_uuid = ?", contactTagUuids)

	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("delete contact tag user failed")
	}

	return nil
}
