package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type ITicketOT interface {
	InsertTicketOTs(ctx context.Context, ticketOT *model.TicketOT) error
	GetTicketOTs(ctx context.Context) (*[]model.TicketOT, error)
	GetTicketOTInDay(ctx context.Context, channel []string) (*[]model.TicketOT, error)
	UpdateTicketOT(ctx context.Context, ticketOT *model.TicketOT) error
	GetTicketOTById(ctx context.Context, ticketOTUuid string) (*model.TicketOT, error)
}

var TicketOT ITicketOT
