package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type ISLAPolicy interface {
	GetSLAPolicyOfFilter(ctx context.Context, ticketCategoryUuid, priority, status string) (*model.SlaPolicy, error)
	InsertSLAPolicies(ctx context.Context, slaPolicies []model.SlaPolicy) error
	GetSLAPolicyByInfo(ctx context.Context, domainUuid, ticketCategoryUuid, status, priority string) (*model.SlaPolicyInfo, error)
	GetSLAPolicyById(ctx context.Context, domainUuid, ticketCategoryUuid string) (*[]model.SlaPolicy, error)
	GetSLAPolicyInfoById(ctx context.Context, domainUuid, ticketCategoryUuid string) (*[]model.SlaPolicyInfo, error)
}

var SLAPolicyRepo ISLAPolicy
