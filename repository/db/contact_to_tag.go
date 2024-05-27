package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"time"
)

type ContactToTagRepo struct{}

func NewContactToTagRepo() repository.IContactToTag {
	repo := &ContactToTagRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *ContactToTagRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ContactToTag)(nil)); err != nil {
		panic(err)
	}
}

func (repo *ContactToTagRepo) InitColumn() {
}

func (repo *ContactToTagRepo) InitIndex() {

}
