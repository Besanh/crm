package model

import (
	"time"

	"github.com/uptrace/bun"
)

type TicketSLA struct {
	bun.BaseModel      `bun:"ticket_sla,alias:ticket_sla"`
	DomainUuid         string `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	TicketSlaUuid      string `json:"ticket_sla_uuid" bun:"ticket_sla_uuid,type: char(36),pk,notnull,default:uuid_generate_v4()"`
	TicketCategoryUuid string `json:"ticket_category_uuid" bun:"ticket_category_uuid,type:text,notnull"`
	// TicketCategory *TicketCategory `json:"ticket_category_uuid" bun:"rel:belongs-to,join:ticket_category_uuid=ticket_category_uuid"`
	TicketUuid   string    `json:"ticket_uuid" bun:"ticket_uuid,type: char(36),notnull"`
	TicketStatus string    `json:"ticket_status" bun:"ticket_status,type:text,notnull"`
	Status       string    `json:"status" bun:"status,type:text"`
	Stage        string    `json:"stage" bun:"stage,type:text,notnull"`
	StartTime    time.Time `json:"start_time" bun:"start_time,nullzero,notnull,default:current_timestamp"`
	EndTime      time.Time `json:"end_time" bun:"end_time,nullzero,notnull,default:current_timestamp"`
	CloseTime    time.Time `json:"close_time" bun:"close_time,nullzero,notnull,default:current_timestamp"`
	CreatedAt    time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
}
