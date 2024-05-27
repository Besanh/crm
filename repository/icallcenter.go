package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type ICallCenter interface {
	GetCallCenterQueues(ctx context.Context, domainUuid string, filter model.CallCenterQueueFilter, limit, offset int) ([]model.CallCenterQueue, int, error)
	GetCallCenterQueueById(ctx context.Context, callCenterQueueUuid string) (*model.CallCenterQueue, error)
	GetCallCenterQueueByExtension(ctx context.Context, domainUuid string, extension string) (*model.CallCenterQueue, error)
	GetCallCenterAgentById(ctx context.Context, callCenterAgentUuid string) (*model.CallCenterAgent, error)
	GetCallCenterAgentsOfCallCenterQueues(ctx context.Context, callCenterQueueUuid ...string) (*[]model.CallCenterAgent, error)
	GetCallCenterAgentByUserId(ctx context.Context, domainUuid, userUuid string) (*model.CallCenterAgent, error)
	UpdateCallcenterAgentStatus(ctx context.Context, domainUuid, agentId, status string) error
	DeleteCallCenterTiersOfAgent(ctx context.Context, domainUuid, agentUuid string) error
	GetCallCenterTiersReady(ctx context.Context, callCenterQueueUuid string) (*[]model.CallCenterTierWithExtension, error)
	GetCallCenterTierOfAgent(ctx context.Context, domainUuid, agentUuid string) (*[]model.CallCenterTier, error)
	InsertCallCenterTier(ctx context.Context, callCenterTier *model.CallCenterTier) error
	UpdateCallCenterQueueStrategy(ctx context.Context, callCenterQueueUuid, strategy string) error
	DeleteCallCenterQueue(ctx context.Context, callCenterQueueUuid string) error
}

var CallCenterRepo ICallCenter
