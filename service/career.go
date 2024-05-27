package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/repository"
	"context"
)

type (
	ICareer interface {
		GetCareers(ctx context.Context, filter model.CareerFilter) (int, any)
	}
	Career struct{}
)

func NewCareer() ICareer {
	return &Career{}
}

func (service *Career) GetCareers(ctx context.Context, filter model.CareerFilter) (int, any) {
	careers, err := repository.CareerRepo.GetCareers(ctx, filter)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(careers)
}
