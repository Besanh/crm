package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"time"
)

type ContactToGroupRepo struct{}

func NewContactToGroupRepo() repository.IContactToGroup {
	repo := &ContactToGroupRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *ContactToGroupRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ContactToGroup)(nil)); err != nil {
		panic(err)
	}
}

func (repo *ContactToGroupRepo) InitColumn() {
}

func (repo *ContactToGroupRepo) InitIndex() {

}
