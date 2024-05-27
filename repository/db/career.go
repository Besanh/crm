package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

type CareerRepo struct{}

func NewCareer() repository.ICareer {
	repo := &CareerRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *CareerRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	repository.FusionSqlClient.GetDB().RegisterModel((*model.Career)(nil))
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.Career)(nil)); err != nil {
		panic(err)
	}
}

func (repo *CareerRepo) InitColumn() {

}
func (repo *CareerRepo) InitIndex() {

}

func (repo *CareerRepo) GetCareers(ctx context.Context, filter model.CareerFilter) (*[]model.CareerView, error) {
	result := new([]model.CareerView)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result)
	if len(filter.CareerCode) > 0 {
		if filter.IsSearchExactly.Valid {
			query.Where("career_code = ?", filter.CareerCode)
		} else {
			query.Where("? ILIKE ?", bun.Ident("career_code"), "%"+filter.CareerCode+"%")
		}
	}
	if len(filter.CareerName) > 0 {
		if filter.IsSearchExactly.Valid {
			query.Where("career_name = ?", filter.CareerName)
		} else {
			query.Where("? ILIKE ?", bun.Ident("career_name"), "%"+filter.CareerName+"%")
		}
	}
	if len(filter.Source) > 0 {
		query.Where("source IN (?)", bun.In(filter.Source))
	}
	err := query.Scan(ctx)
	if err != sql.ErrNoRows {
		return result, nil
	} else if err != nil {
		return result, err
	}
	return result, nil
}
