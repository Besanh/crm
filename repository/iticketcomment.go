package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type ITicketComment interface {
	InsertTicketComment(ctx context.Context, domainUuid string, ticketComment model.TicketComment) error
	GetCommentByTicketId(ctx context.Context, domainUuid, ticketUuid string) (*[]model.TicketComment, error)
	DeleteTicketComment(ctx context.Context, domainUuid, ticketUuid string) error
}

var TicketCommentRepo ITicketComment
