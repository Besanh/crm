package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"time"
)

type SourcePluginRepo struct{}

func NewSourcePluginRepo() repository.ISourcePlugin {
	repo := &SourcePluginRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *SourcePluginRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	repository.FusionSqlClient.GetDB().RegisterModel((*model.SourcePlugin)(nil))
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.SourcePlugin)(nil)); err != nil {
		panic(err)
	}
}

func (repo *SourcePluginRepo) InitColumn() {

}

func (repo *SourcePluginRepo) InitIndex() {
}
