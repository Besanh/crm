package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type ITicketSLA interface {
	UpdateTicketSLA(ctx context.Context, ticketSLA *model.TicketSLA) error
	GetTicketSLAByTicketId(ctx context.Context, domainUuid, ticketUuid string, state string) (*model.TicketSLA, error)
	UpdateTicketSLAByTicketIds(ctx context.Context, ticketUuids []string) error
	GetTicketSLAs(ctx context.Context, domainUuid, ticketUuid string) (*[]model.TicketSLA, error)
	GetTicketSLAById(ctx context.Context, domainUuid, ticketSlaUuid string) (*model.TicketSLA, error)
	GetTicketSLALatestByTicketId(ctx context.Context, domainUuid, ticketUuid string) (*model.TicketSLA, error)
	GetTicketSlaByStatus(ctx context.Context, domainUuid, ticketSlaUuid, status, stage string) (*model.TicketSLA, error)
	GetTicketSlaLatestByStatusTicketAndId(ctx context.Context, domainUuid, ticketSlaUuid, status string) (*model.TicketSLA, error)
	DeleteTicketSla(ctx context.Context, domainUuid, ticketUuid string) error
}

var TicketSLARepo ITicketSLA
