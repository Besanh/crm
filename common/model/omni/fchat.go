package omni

import (
	"time"
)

type FChatInfo struct {
	SenderId       string `json:"sender_id" bun:"-"`
	PageId         string `json:"page_id" bun:"-"`
	PageName       string `json:"page_name" bun:"-"`
	ShopId         string `json:"shop_id" bun:"-"`
	Token          string `json:"token" bun:"-"`
	ConversationId string `json:"conversation_id" bun:"-"`
	Fullname       string `json:"fullname" bun:"-"`
}

type TicketPendingFchatPost struct {
	DomainUuid         string `json:"domain_uuid" bun:"domain_uuid"`
	TicketUuid         string `json:"ticket_uuid" bun:"ticket_uuid,"`
	CustomerId         string `json:"customer_id" bun:"customer_id"`
	OmniUuid           string `json:"omni_uuid" bun:"omni_uuid"`
	CustomerPhone      string `json:"customer_phone" bun:"customer_phone"`
	AssigneeUuid       string `json:"assignee_uuid" bun:"assignee_uuid"`
	UnitUuid           string `json:"unit_uuid" bun:"unit_uuid"`
	Subject            string `json:"subject" bun:"subject"`
	Content            string `json:"content" bun:"content"`
	SolvingContent     string `json:"solving_content" bun:"solving_content"`
	Priority           string `json:"priority" bun:"priority"`
	Channel            string `json:"channel" bun:"channel"`
	SolutionUuid       string `json:"solution_uuid" bun:"solution_uuid,type"`
	Status             string `json:"status" bun:"status"`
	TicketCategoryUuid string `json:"ticket_category_uuid" bun:"ticket_category_uuid"`

	CreatedAt  time.Time `json:"-" bun:"created_at"`
	UpdatedAt  time.Time `json:"-" bun:"updated_at"`
	CreatedBy  string    `json:"created_by" bun:"created_by"`
	UpdatedBy  string    `json:"updated_by" bun:"updated_by"`
	TicketCode string    `json:"ticket_code" bun:"ticket_code"`
	SlaStatus  string    `json:"sla_status" bun:"-"`

	IsProcessed bool   `json:"is_processed" bun:"is_processede"`
	SenderId    string `json:"sender_id" bun:"sender_id,type:text"`
	FullName    string `json:"full_name" bun:"full_name,type:text"`
	PageId      string `json:"page_id" bun:"page_id,type:text"`
	PageName    string `json:"page_name" bun:"page_name,type:text"`
	ShopId      string `json:"shop_id" bun:"shop_id,type:text"`
	Token       string `json:"token" bun:"token,type:text"`
	OwnerId     string `json:"owner_id" bun:"-"`
}
