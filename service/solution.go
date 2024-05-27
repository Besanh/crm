package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type (
	ISolution interface {
		PostSolution(ctx context.Context, domainId, userUuid string, solutionPost model.SolutionPost) (int, any)
		GetSolutions(ctx context.Context, domainUuid, userUuid string, filter model.SolutionFilter, limit, offset int) (int, any)
		GetSolutionById(ctx context.Context, domainUuid, id string) (int, any)
		PutSolutionById(ctx context.Context, domainUuid, userUuid, id string, solution model.Solution) (int, any)
		DeleteSolutionById(ctx context.Context, domainUuid, userUuid, id string) (int, any)
		ExportSolutions(ctx context.Context, domainUuid, userUuid, fileType string, filter model.SolutionFilter) (string, error)
	}
	Solution struct{}
)

func NewSolution() ISolution {
	return &Solution{}
}
func (s *Solution) PostSolution(ctx context.Context, domainUuid, userUuid string, solutionPost model.SolutionPost) (int, any) {
	solution := model.Solution{
		DomainUuid:   domainUuid,
		SolutionUuid: uuid.NewString(),
		SolutionName: solutionPost.SolutionName,
		SolutionCode: solutionPost.SolutionCode,
		Status:       solutionPost.Status,
		UnitUuid:     solutionPost.UnitUuid,
		CreatedBy:    userUuid,
		CreatedAt:    time.Now(),
	}

	if err := repository.SolutionRepo.InsertSolution(ctx, solution); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.Created(map[string]any{
		"id": solution.SolutionUuid,
	})
}

func (s *Solution) GetSolutions(ctx context.Context, domainUuid, userUuid string, filter model.SolutionFilter, limit, offset int) (int, any) {
	total, solutions, err := repository.SolutionRepo.GetSolutions(ctx, domainUuid, limit, offset, filter)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.Pagination(solutions, total, limit, offset)
}

func (s *Solution) GetSolutionById(ctx context.Context, domainUuid, id string) (int, any) {
	solution, err := repository.SolutionRepo.GetSolutionById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(map[string]any{
		"solution": solution,
	})
}

func (s *Solution) PutSolutionById(ctx context.Context, domainUuid, userUuid, id string, solution model.Solution) (int, any) {
	solutionExist, err := repository.SolutionRepo.GetSolutionById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(solutionExist.SolutionUuid) < 1 {
		return response.ServiceUnavailableMsg("solution not found")
	}

	solutionExist.SolutionName = solution.SolutionName
	solutionExist.SolutionCode = solution.SolutionCode
	solutionExist.Status = solution.Status
	solutionExist.UnitUuid = solution.UnitUuid
	solutionExist.UpdatedBy = userUuid
	solutionExist.UpdatedAt = time.Now()

	if err := repository.SolutionRepo.PutSolutionById(ctx, domainUuid, solution); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{})
}

func (s *Solution) DeleteSolutionById(ctx context.Context, domainUuid, userUuid, id string) (int, any) {
	solutionExist, err := repository.SolutionRepo.GetSolutionById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(solutionExist.SolutionUuid) < 1 {
		return response.ServiceUnavailableMsg("solution not found")
	}

	if err := repository.SolutionRepo.DeleteSolutionById(ctx, domainUuid, userUuid, id); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"id": id,
	})
}
