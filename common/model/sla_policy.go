package model

import (
	"time"

	"github.com/uptrace/bun"
)

type SlaPolicy struct {
	bun.BaseModel      `bun:"sla_policy,alias:sp"`
	DomainUuid         string    `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	TicketCategoryUuid string    `json:"ticket_category_uuid" bun:"ticket_category_uuid,type: char(36),notnull"`
	SlaPolicyUuid      string    `json:"sla_policy_uuid" bun:"sla_policy_uuid,type: char(36),pk,notnull"`
	Priority           string    `json:"priority" bun:"priority,type:sla_priority,notnull,default:'normal'"`
	Status             string    `json:"status" bun:"status,type:text,notnull"`
	ResponseTime       int       `json:"response_time" bun:"response_time,type:numeric,notnull"`
	ResponseType       string    `json:"response_type" bun:"response_type,type:text"`
	CreatedBy          string    `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy          string    `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt          time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt          time.Time `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
}

type SlaPolicyInfo struct {
	bun.BaseModel      `bun:"sla_policy,alias:sp"`
	DomainUuid         string    `json:"domain_uuid" bun:"domain_uuid"`
	SLAPolicyUuid      string    `json:"sla_policy_uuid" bun:"sla_policy_uuid"`
	TicketCategoryUuid string    `json:"ticket_category_uuid" bun:"ticket_category_uuid"`
	TicketCategoryName string    `json:"ticket_category_name" bun:"-"`
	Priority           string    `json:"priority" bun:"priority,type:sla_priority"`
	Status             string    `json:"status" bun:"status"`
	ResponseTime       int       `json:"response_time" bun:"response_time"`
	ResponseType       string    `json:"response_type" bun:"response_type"`
	CreatedBy          string    `json:"created_by" bun:"created_by"`
	UpdatedBy          string    `json:"updated_by" bun:"updated_by"`
	CreatedAt          time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" bun:"updated_at"`
}
