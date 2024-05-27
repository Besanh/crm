package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"time"
)

type ContactToCareerRepo struct{}

func NewContactToCareerRepo() repository.IContactToCareer {
	repo := &ContactToCareerRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *ContactToCareerRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ContactToCareer)(nil)); err != nil {
		panic(err)
	}
}

func (repo *ContactToCareerRepo) InitColumn() {
}

func (repo *ContactToCareerRepo) InitIndex() {

}
