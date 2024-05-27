package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type Lead struct{}

func NewLead() repository.ILeadRepo {
	repo := &Lead{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *Lead) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.Lead)(nil)); err != nil {
		panic(err)
	}
}

func (repo *Lead) InitColumn() {

}

func (repo *Lead) InitIndex() {

}

func (repo *Lead) GetLeadById(ctx context.Context, domainUuid, leadUuid string) (*model.Lead, error) {
	lead := new(model.Lead)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(lead).
		Where("lead_uuid = ?", leadUuid).
		Where("domain_uuid = ?", domainUuid).
		Limit(1)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return lead, nil
	}
}

func (repo *Lead) UpdateLead(ctx context.Context, lead *model.Lead) error {
	query := repository.FusionSqlClient.GetDB().NewUpdate().
		Model(lead).
		WherePK()
	resp, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected != 1 {
		return errors.New("update lead failed")
	}
	return nil
}
