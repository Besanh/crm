package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type ClassifyCareer struct{}

func NewClassifyCareerRepo() repository.IClassifyCareerRepo {
	repo := &ClassifyCareer{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *ClassifyCareer) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	repository.FusionSqlClient.GetDB().RegisterModel((*model.ClassifyCareer)(nil))
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ClassifyCareer)(nil)); err != nil {
		panic(err)
	}
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ClassifyGroup)(nil)); err != nil {
		panic(err)
	}
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ClassifyTag)(nil)); err != nil {
		panic(err)
	}
}

func (repo *ClassifyCareer) InitColumn() {

}
func (repo *ClassifyCareer) InitIndex() {

}

func (repo *ClassifyCareer) InsertClassifyCareer(ctx context.Context, classifyCareer model.ClassifyCareer) error {
	resp, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&classifyCareer).Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("insert classify career failed")
	}

	return nil
}

func (repo *ClassifyCareer) GetClassifyCareers(ctx context.Context, domainUuid string, limit, offset int) (int, []model.ClassifyCareer, error) {
	var result []model.ClassifyCareer
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(&result).
		Distinct().
		Where("domain_uuid = ?", domainUuid)
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	query.Order("created_at DESC")
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return 0, result, err
	} else if err == sql.ErrNoRows {
		return 0, result, errors.New("classify career not found")
	}
	return total, result, nil
}

func (repo *ClassifyCareer) GetClassifyCareerById(ctx context.Context, domainUuid, id string) (*model.ClassifyCareer, error) {
	result := new(model.ClassifyCareer)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(&result).
		Where("domain_uuid = ?", domainUuid).
		Where("classify_career_uuid = ?", id)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *ClassifyCareer) PutClassifyCareerById(ctx context.Context, domainUuid string, classifyCareer model.ClassifyCareer) error {
	resp, err := repository.FusionSqlClient.GetDB().NewUpdate().Model(&classifyCareer).
		Where("domain_uuid = ?", domainUuid).
		Where("classify_career_uuid = ?", classifyCareer.ClassifyCareerUuid).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("update classify career failed")
	}

	return nil
}
