package model

import (
	"time"

	"contactcenter-api/common/model/omni"

	"github.com/uptrace/bun"
)

type Ticket struct {
	bun.BaseModel           `bun:"ticket,alias:ticket"`
	DomainUuid              string           `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	TicketUuid              string           `json:"ticket_uuid" bun:"ticket_uuid,type: char(36),pk,notnull"`
	TicketCategoryUuid      string           `json:"ticket_category_uuid" bun:"ticket_category_uuid,type: char(36),notnull"`
	TicketCategoryHierarchy map[int]string   `json:"ticket_category_hierarchy" bun:"ticket_category_hierarchy,type:text,nullzero"`
	OmniUuid                string           `json:"omni_uuid" bun:"omni_uuid,type:text"`
	TicketCode              string           `json:"ticket_code" bun:"ticket_code,type:text"`
	CustomerId              string           `json:"customer_id" bun:"customer_id,type:text"`
	FullName                string           `json:"full_name" bun:"full_name,type:text"`
	CustomerPhone           string           `json:"customer_phone" bun:"customer_phone,type:text"`
	AssigneeUuid            string           `json:"assignee_uuid" bun:"assignee_uuid,type: char(36)"`
	UnitUuid                string           `json:"unit_uuid" bun:"unit_uuid,type: char(36)"`
	SolutionUuid            string           `json:"solution_uuid" bun:"solution_uuid,type: char(36)"`
	Subject                 string           `json:"subject" bun:"subject,type:text"`
	Priority                string           `json:"priority" bun:"priority,type:text"`
	Channel                 string           `json:"channel" bun:"channel,type:text"`
	Status                  string           `json:"status" bun:"status,type:text"`
	Content                 string           `json:"content" bun:"content,type:text"`
	SolvingContent          string           `json:"solving_content" bun:"solving_content,type:text"`
	Attachment              []string         `json:"attachment" bun:"attachment,type:text[]"`
	OmniInfo                string           `json:"omni_info" bun:"omni_info,type:text"`
	Omni                    *omni.Omni       `bun:"rel:has-one,join:omni_uuid=omni_uuid"`
	Unit                    *Unit            `json:"unit" bun:"rel:has-one,join:unit_uuid=unit_uuid"`
	SolutionItem            *Solution        `json:"solution" bun:"rel:has-one,join:solution_uuid=solution_uuid"`
	Assignee                *UserView        `json:"assignee" bun:"-"`
	TicketSla               []*TicketSLA     `json:"ticket_sla" bun:"rel:has-many,join:ticket_uuid=ticket_uuid"`
	TicketComment           []*TicketComment `json:"ticket_comment" bun:"rel:has-many,join:ticket_uuid=ticket_uuid"`
	TicketCategory          *TicketCategory  `json:"ticket_category" bun:"rel:has-one,join:ticket_category_uuid=ticket_category_uuid"`
	User                    *UserBasicInfo   `json:"user" bun:"rel:belongs-to,join:updated_by=user_uuid"`
	SlaStatus               string           `json:"sla_status" bun:"-"`
	CreatedBy               string           `json:"created_by" bun:"created_by,type: char(36),nullzero"`
	UpdatedBy               string           `json:"updated_by" bun:"updated_by,type: char(36),nullzero"`
	CreatedAt               time.Time        `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt               time.Time        `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
}

type TicketPost struct {
	bun.BaseModel           `bun:"ticket,alias:ticket"`
	DomainUuid              string   `json:"domain_uuid" bun:"domain_uuid"`
	TicketUuid              string   `json:"ticket_uuid" bun:"ticket_uuid"`
	TicketCategoryUuid      string   `json:"ticket_category_uuid" bun:"ticket_category_uuid"`
	TicketCategoryHierarchy string   `json:"ticket_category_hierarchy" bun:"ticket_category_hierarchy"`
	OmniUuid                string   `json:"omni_uuid" bun:"omni_uuid"`
	TicketCode              string   `json:"ticket_code" bun:"ticket_code"`
	CustomerId              string   `json:"customer_id" bun:"customer_id"`
	FullName                string   `json:"full_name" bun:"full_name"`
	CustomerPhone           string   `json:"customer_phone" bun:"customer_phone"`
	AssigneeUuid            string   `json:"assignee_uuid" bun:"assignee_uuid"`
	UnitUuid                string   `json:"unit_uuid" bun:"unit_uuid"`
	SolutionUuid            string   `json:"solution_uuid" bun:"solution_uuid"`
	Subject                 string   `json:"subject" bun:"subject"`
	Priority                string   `json:"priority" bun:"priority"`
	Channel                 string   `json:"channel" bun:"channel"`
	Status                  string   `json:"status" bun:"status"`
	Content                 string   `json:"content" bun:"content"`
	SolvingContent          string   `json:"solving_content" bun:"solving_content"`
	Attachment              []string `json:"attachment" bun:"attachment"`
}

type TicketFilter struct {
	Common       string `json:"common" form:"common"`
	Subject      string `json:"subject" form:"subject"`
	CustomerId   string `json:"customer_id" form:"customer_id"`
	AssigneeUuid string `json:"assignee_uuid" form:"assignee_uuid"`
	Status       string `json:"status" form:"status"`
	CategoryUuid string `json:"ticket_category_uuid" form:"ticket_category_uuid"`
	FromDate     string `json:"from_date" form:"from_date"`
	ToDate       string `json:"to_date" form:"to_date"`
	PhoneNumber  string `json:"phone_number" form:"phone_number"`
	Priority     string `json:"priority" form:"priority"`
	SLA          string `json:"sla" form:"sla"`
	TicketCode   string `json:"ticket_code" form:"ticket_code"`
	Content      string `json:"content" form:"content"`
	SenderId     string `json:"sender_id" form:"sender_id"`
	FullName     string `json:"full_name" form:"full_name"`
	Channel      string `json:"channel" form:"channel"`
	CreatedBy    string `json:"created_by" form:"created_by"`
	FileType     string `json:"file_type"`
	Limit        int
	Offset       int
}

type TicketExport struct {
	bun.BaseModel      `bun:"ticket,alias:ticket"`
	DomainUuid         string    `bun:"domain_uuid"`
	TicketUuid         string    `bun:"ticket_uuid"`
	Subject            string    `bun:"subject"`
	CustomerId         string    `bun:"customer_id"`
	CustomerPhone      string    `bun:"customer_phone"`
	TicketCategoryName string    `bun:"ticket_category_name"`
	Channel            string    `bun:"channel"`
	Priority           string    `bun:"priority"`
	Content            string    `bun:"content"`
	SolvingContent     string    `bun:"-"`
	TicketStatus       string    `bun:"status"`
	Solution           string    `bun:"solution_name"`
	CreatedBy          string    `bun:"user_created"`
	UpdatedBy          string    `bun:"user_updated"`
	AssigneeUserName   string    `bun:"user_assignee"`
	SlaResult          string    `bun:"-"`
	TicketCode         string    `bun:"ticket_code"`
	SlaStatus          string    `bun:"sla_status"`
	FullName           string    `bun:"full_name"`
	CreatedAt          time.Time `bun:"created_at"`
	UpdatedAt          time.Time `bun:"updated_at"`
}

type OmniInfo struct {
	Fchat *omni.FChatInfo `json:"fchat"`
}
