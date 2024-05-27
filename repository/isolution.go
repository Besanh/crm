package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type ISolution interface {
	InsertSolution(ctx context.Context, solution model.Solution) error
	GetSolutions(ctx context.Context, domainUuid string, limit, offset int, filter model.SolutionFilter) (int, *[]model.Solution, error)
	GetSolutionById(ctx context.Context, domainUuid, id string) (*model.Solution, error)
	PutSolutionById(ctx context.Context, domainUuid string, solution model.Solution) error
	DeleteSolutionById(ctx context.Context, domainUuid, userUuid, id string) error
}

var SolutionRepo ISolution
