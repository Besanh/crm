package model

import (
	"time"

	"github.com/uptrace/bun"
)

type TicketLog struct {
	bun.BaseModel `bun:"ticket_log,alias:ticket_log"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	TicketLogUuid string `json:"ticket_log_id" bun:"ticket_log_id,type: char(36),pk"`
	TicketUuid    string `json:"ticket_uuid" bun:"ticket_uuid,type: char(36),notnull"`
	Status        string `json:"status" bun:"status,type:text"`

	CreatedAt time.Time `json:"created_at" bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
	CreatedBy string    `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy string    `json:"updated_by" bun:"updated_by,type:text"`

	Content      string `json:"content" bun:"content,type:text"`
	Type         string `json:"type" bun:"type,type:text"`
	AssigneeFrom string `json:"assignee_from" bun:"assignee_from,type:text"`
	AssigneeTo   string `json:"assignee_to" bun:"assignee_to,type:text"`
}

type TicketLogFilter struct {
	TicketUuid   string `json:"ticket_uuid" form:"ticket_uuid"`
	TicketStatus string `json:"status" form:"status"`
	Type         string `json:"type" form:"type"`
}

type TicketInfo struct {
	bun.BaseModel      `bun:"ticket,alias:ticket"`
	DomainUuid         string            `json:"domain_uuid" bun:"domain_uuid"`
	TicketUuid         string            `json:"ticket_uuid" bun:"ticket_uuid,type: char(36),pk"`
	CustomerId         string            `json:"customer_id" bun:"customer_id"`
	AssigneeUuid       string            `json:"assignee_uuid" bun:"assignee_uuid"`
	DepartmentId       string            `json:"department_id" bun:"department_id"`
	Subject            string            `json:"subject" bun:"subject"`
	Content            string            `json:"content" bun:"content"`
	Priority           string            `json:"priority" bun:"priority"`
	Channel            string            `json:"channel" bun:"channel"`
	SolutionUuid       string            `json:"solution_uuid" bun:"solution_uuid"`
	Status             string            `json:"status" bun:"status"`
	TicketCategoryUuid string            `json:"ticket_category_uuid" bun:"ticket_category_uuid"`
	TicketCategory     *TicketCategory   `json:"ticket_category" bun:"rel:has-one,join:ticket_category_uuid=ticket_category_uuid"`
	TicketSla          *TicketSLA        `json:"ticket_sla" bun:"-"`
	Blacklist          map[string]string `json:"-" bun:"-"`
	TicketCode         string            `json:"ticket_code" bun:"ticket_code"`
	FullName           string            `json:"full_name" bun:"full_name"`
	SenderId           string            `json:"sender_id" bun:"sender_id"`
	PageId             string            `json:"page_id" bun:"page_id"`
	ShopId             string            `json:"shop_id" bun:"shop_id"`
	ConversationId     string            `json:"conversation_id" bun:"conversation_id"`
	CreatedAt          time.Time         `json:"created_at" bun:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at" bun:"updated_at"`
	CreatedBy          string            `json:"created_by" bun:"created_by"`
	UpdatedBy          string            `json:"updated_by" bun:"updated_by"`
}
