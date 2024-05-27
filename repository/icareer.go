package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type ICareer interface {
	GetCareers(ctx context.Context, filter model.CareerFilter) (*[]model.CareerView, error)
}

var CareerRepo ICareer
