package model

import (
	"time"

	"github.com/uptrace/bun"
)

type TicketComment struct {
	bun.BaseModel     `bun:"ticket_comment,alias:ticket_comment"`
	DomainUuid        string    `json:"domain_uuid" bun:"domain_uuid,type:uuid,notnull"`
	TicketCommentUuid string    `json:"ticket_comment_uuid" bun:"ticket_comment_uuid,type: char(36),pk,notnull"`
	TicketUuid        string    `json:"ticket_uuid" bun:"ticket_uuid,type:uuid,notnull"`
	Message           string    `json:"message" bun:"message,type:text,default:null"`
	Status            bool      `json:"status" bun:"status,type:boolean,default:'true'"`
	User              *User     `json:"user" bun:"rel:belongs-to,join:created_by=user_uuid"`
	CreatedBy         string    `json:"created_by" bun:"created_by,type:uuid,nullzero"`
	UpdatedBy         string    `json:"updated_by" bun:"updated_by,type:uuid,nullzero"`
	CreatedAt         time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt         time.Time `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
}
