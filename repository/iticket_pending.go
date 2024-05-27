package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type ITicketPending interface {
	InsertTicketPending(ctx context.Context, TicketPending *model.TicketPending) error
	GetTicketPendings(ctx context.Context, domainUuid string, limit, offset int, filter model.TicketFilter) (*[]model.TicketPending, int, error)
	GetTicketPendingInDay(ctx context.Context, channel []string) (*[]model.TicketPending, error)
	UpdateTicketPending(ctx context.Context, TicketPending *model.TicketPending) error
	GetTicketPendingById(ctx context.Context, TicketPendingUuid string) (*model.TicketPending, error)
	DeleteTicketPendingById(ctx context.Context, ticketPendingUuid string) error
}

var TicketPendingRepo ITicketPending
