package db

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type SolutionRepo struct{}

func NewSolution() repository.ISolution {
	repo := &SolutionRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *SolutionRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.Solution)(nil)); err != nil {
		panic(err)
	}
}

func (repo *SolutionRepo) InitColumn() {

}

func (repo *SolutionRepo) InitIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Solution)(nil)).IfNotExists().Index("idx_solution_domain_uuid").Column("domain_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Solution)(nil)).IfNotExists().Index("idx_solution_uuid").Column("solution_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Solution)(nil)).IfNotExists().Index("idx_solution_code").Column("solution_code").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Solution)(nil)).IfNotExists().Index("idx_solution_name").Column("solution_name").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Solution)(nil)).IfNotExists().Index("idx_solution_status").Column("status").Exec(ctx); err != nil {
		log.Error(err)
	}
}

func (repo *SolutionRepo) InsertSolution(ctx context.Context, solution model.Solution) error {
	resp, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&solution).Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("insert solution failed")
	}

	return nil
}

func (repo *SolutionRepo) GetSolutions(ctx context.Context, domainUuid string, limit, offset int, filter model.SolutionFilter) (int, *[]model.Solution, error) {
	solutions := new([]model.Solution)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(solutions).
		Relation("Unit").
		Where("solution.domain_uuid = ?", domainUuid)
	if len(filter.SolutionName) > 0 {
		query = query.Where("solution_name = ?", filter.SolutionName)
	}
	if len(filter.SolutionCode) > 0 {
		query = query.Where("solution_code = ?", filter.SolutionCode)
	}
	if filter.Status.Valid {
		query = query.Where("solution.status = ?", filter.Status.Bool)
	}
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	query.Order("solution.created_at DESC")
	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, nil, nil
	} else if err != nil {
		return 0, nil, err
	}

	return total, solutions, nil
}

func (repo *SolutionRepo) GetSolutionById(ctx context.Context, domainUuid, solutionUuid string) (*model.Solution, error) {
	solution := new(model.Solution)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(solution).
		Relation("Unit").
		Where("solution.domain_uuid = ?", domainUuid).
		Where("solution_uuid = ?", solutionUuid)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return solution, nil
}

func (repo *SolutionRepo) PutSolutionById(ctx context.Context, domainUuid string, solution model.Solution) error {
	resp, err := repository.FusionSqlClient.GetDB().NewUpdate().Model(&solution).
		Where("domain_uuid = ?", domainUuid).
		Where("solution_uuid = ?", solution.SolutionUuid).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("update solution failed")
	}

	return nil
}

func (repo *SolutionRepo) DeleteSolutionById(ctx context.Context, domainUuid, userUuid, id string) error {
	res, err := repository.FusionSqlClient.GetDB().NewDelete().Model(&model.Solution{}).
		Where("domain_uuid = ?", domainUuid).
		Where("solution_uuid = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("delete solution failed")
	}

	return nil
}
