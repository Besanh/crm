package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type ILeadRepo interface {
	GetLeadById(ctx context.Context, domainUuid, leadUuid string) (*model.Lead, error)
	UpdateLead(ctx context.Context, lead *model.Lead) error
}

var LeadRepo ILeadRepo
