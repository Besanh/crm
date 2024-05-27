package model

import (
	"time"

	"github.com/uptrace/bun"
)

type TicketOT struct {
	bun.BaseModel `bun:"ticket_ot,alias:ticket_ot"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid,type:uuid,notnull"`
	TicketUuid    string `json:"ticket_uuid" bun:"ticket_uuid,type:uuid,pk,notnull"`
	CustomerId    string `json:"customer_id" bun:"customer_id,type:text"`
	CustomerPhone string `json:"customer_phone" bun:"customer_phone,type:text"`
	AssigneeId    string `json:"assignee_id" bun:"assignee_id,type:uuid,default:NULL"`

	DepartmentId       string `json:"department_id" bun:"department_id,type:uuid,default:NULL"`
	Subject            string `json:"subject" bun:"subject"`
	Content            string `json:"content" bun:"content"`
	SolvingContent     string `json:"solving_content" bun:"solving_content"`
	Priority           string `json:"priority" bun:"priority"`
	Channel            string `json:"channel" bun:"channel"`
	SolutionUuid       string `json:"solution_uuid" bun:"solution_uuid,type:uuid,default:NULL"`
	Status             string `json:"status" bun:"status"`
	TicketCategoryUuid string `json:"ticket_category_uuid" bun:"ticket_category_uuid,type:uuid,default:NULL"`

	CreatedAt  time.Time `json:"-" bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt  time.Time `json:"-" bun:"updated_at,nullzero,default:NULL"`
	CreatedBy  string    `json:"created_by" bun:"created_by,type:uuid,nullzero"`
	UpdatedBy  string    `json:"updated_by" bun:"updated_by,type:uuid,nullzero"`
	TicketCode string    `json:"ticket_code" bun:"ticket_code,type:text,default:NULL"`
	SlaStatus  string    `json:"sla_status" bun:"-"`

	IsProcessed bool   `json:"is_processed" bun:"is_processed,type:boolean,default:false"`
	SenderId    string `json:"sender_id" bun:"sender_id,type:text"`
	FullName    string `json:"full_name" bun:"full_name,type:text"`
	PageId      string `json:"page_id" bun:"page_id,type:text"`
	PageName    string `json:"page_name" bun:"page_name,type:text"`
	ShopId      string `json:"shop_id" bun:"shop_id,type:text"`
	Token       string `json:"token" bun:"token,type:text"`
	OwnerId     string `json:"owner_id" bun:"-"`
}

type TicketOtUser struct {
	DomainUuid string `json:"domain_uuid"`
	UserUuid   string `json:"user_uuid"`
	Email      string `json:"email"`
	ZaloId     string `json:"zalo_id"` //social_id
	FacebookId string `json:"facebook_id"`
}
