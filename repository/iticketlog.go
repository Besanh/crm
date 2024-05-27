package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type ITicketLog interface {
	InsertNewTicketLog(ctx context.Context, ticketLog model.TicketLog) error
	GetTicketLogs(ctx context.Context, domainUuid, ticketUuid, ticketLogType, status string, offset, limit int) (*[]model.TicketLog, int, error)
	GetTicketLogsByType(ctx context.Context, domainUuid string, ticketUuid string, ticketLogType string) (*[]model.TicketLog, error)
	DeleteTicketLog(ctx context.Context, domainUuid, ticketUuid string) error
}

var TicketLogRepo ITicketLog
