package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type ITicket interface {
	InsertTicket(ctx context.Context, ticket *model.Ticket, ticketSla *model.TicketSLA, ticketLog *model.TicketLog, ticketComment *model.TicketComment, checkExistTicketComment bool) error
	PutTicket(ctx context.Context, domainUuid string, ticket *model.Ticket, ticketSlaUpdate *model.TicketSLA, ticketSlaNew *model.TicketSLA,
		ticketLog *model.TicketLog) error
	GetTicketById(ctx context.Context, domainUuid, ticketId string) (*model.Ticket, error)
	GetTicketsInfo(ctx context.Context, domainUuid string, limit, offset int, ticketFilter *model.TicketFilter) (*[]model.Ticket, int, error)
	GetLatestTicket(ctx context.Context, domainUuid string) (*model.Ticket, error)
	UpdateConversationIdTicketByTicketId(ctx context.Context, domainUuid, ticketUuid, conversationId string) error
	GetTicketsExport(ctx context.Context, domainUuid string, limit int, offset int, ticketFilter *model.TicketFilter) (*[]model.TicketExport, int, error)
	DeleteTicket(ctx context.Context, domainUuid, ticketUuid string) error
	PatchTicketAttachment(ctx context.Context, domainUuid string, ticket model.Ticket) error
	GetLatestTicketEmail(ctx context.Context, domainUuid string) (*model.Ticket, error)
	GetTicketByProfileUuids(ctx context.Context, domainUuid string, profileUuids []string) (result *[]model.Ticket, err error)
}

var TicketRepo ITicket
