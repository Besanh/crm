package model

import (
	"contactcenter-api/common/model/omni"
	"time"

	"github.com/uptrace/bun"
)

type TicketPending struct {
	bun.BaseModel      `bun:"ticket_pending,alias:tp"`
	DomainUuid         string           `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	TicketUuid         string           `json:"ticket_uuid" bun:"ticket_uuid,type: char(36),pk,notnull"`
	TicketCategoryUuid string           `json:"ticket_category_uuid" bun:"ticket_category_uuid,type: char(36),notnull"`
	OmniUuid           string           `json:"omni_uuid" bun:"omni_uuid,type: char(36),notnull"`
	TicketCode         string           `json:"ticket_code" bun:"ticket_code,type:text"`
	CustomerId         string           `json:"customer_id" bun:"customer_id,type:text"`
	FullName           string           `json:"full_name" bun:"full_name,type:text"`
	CustomerPhone      string           `json:"customer_phone" bun:"customer_phone,type:text"`
	AssigneeUuid       string           `json:"assignee_uuid" bun:"assignee_uuid,type:text"`
	UnitUuid           string           `json:"unit_uuid" bun:"unit_uuid,type:text"`
	SolutionUuid       string           `json:"solution_uuid" bun:"solution_uuid,type:text"`
	Subject            string           `json:"subject" bun:"subject,type:text"`
	Priority           string           `json:"priority" bun:"priority,type:text"`
	Channel            string           `json:"channel" bun:"channel,type:text"`
	Status             string           `json:"status" bun:"status,type:text"`
	IsProcessed        bool             `json:"is_processed" bun:"is_processed,type:bool,default:false"`
	Content            string           `json:"content" bun:"content,type:text"`
	SolvingContent     string           `json:"solving_content" bun:"solving_content,type:text"`
	CreatedBy          string           `json:"created_by" bun:"created_by,type: char(36),nullzero"`
	UpdatedBy          string           `json:"updated_by" bun:"updated_by,type: char(36),nullzero"`
	CreatedAt          time.Time        `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt          time.Time        `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
	Unit               *Unit            `json:"unit" bun:"rel:has-one,join:unit_uuid=unit_uuid"`
	SolutionItem       *Solution        `json:"solution" bun:"rel:has-one,join:solution_uuid=solution_uuid"`
	TicketSla          []*TicketSLA     `json:"ticket_sla" bun:"rel:has-many,join:ticket_uuid=ticket_uuid"`
	TicketComment      []*TicketComment `json:"ticket_comment" bun:"rel:has-many,join:ticket_uuid=ticket_uuid"`
	TicketCategory     *TicketCategory  `json:"ticket_category" bun:"rel:has-one,join:ticket_category_uuid=ticket_category_uuid"`
	User               *UserBasicInfo   `json:"user" bun:"rel:belongs-to,join:updated_by=user_uuid"`
	SlaStatus          string           `json:"sla_status" bun:"-"`
	OmniInfo           *OmniInfo        `json:"omni_info" bun:"omni_info"`
	Omni               *omni.Omni       `bun:"rel:has-one,join:omni_uuid=omni_uuid"`
}

type TicketPendingPost struct {
	bun.BaseModel      `bun:"ticket_pending,alias:tp"`
	DomainUuid         string    `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	TicketUuid         string    `json:"ticket_uuid" bun:"ticket_uuid,type: char(36),pk,notnull"`
	CustomerId         string    `json:"customer_id" bun:"customer_id,type:text"`
	Fullname           string    `json:"full_name" bun:"full_name,type:text"`
	CustomerPhone      string    `json:"customer_phone" bun:"customer_phone,type:text"`
	OmniUuid           string    `json:"omni_uuid" bun:"omni_uuid,type: char(36),default:NULL"`
	AssigneeUuid       string    `json:"assignee_uuid" bun:"assignee_uuid,type: char(36),default:NULL"`
	UnitUuid           string    `json:"unit_uuid" bun:"unit_uuid,type: char(36),default:NULL"`
	Subject            string    `json:"subject" bun:"subject"`
	Priority           string    `json:"priority" bun:"priority"`
	Channel            string    `json:"channel" bun:"channel"`
	Status             string    `json:"status" bun:"status"`
	SolutionUuid       string    `json:"solution_uuid" bun:"solution_uuid,type: char(36),default:NULL"`
	TicketCategoryUuid string    `json:"ticket_category_uuid" bun:"ticket_category_uuid,type: char(36),default:NULL"`
	TicketCode         string    `json:"ticket_code" bun:"ticket_code,type:text,default:NULL"`
	SlaStatus          string    `json:"sla_status" bun:"-"`
	IsProcessed        bool      `json:"is_processed" bun:"is_processed,type:boolean,default:false"`
	SourcePluginUuid   string    `json:"source_plugin_uuid" bun:"source_plugin_uuid,type: char(36),notnull"`
	OmniInfo           *OmniInfo `json:"omni_info" bun:"omni_info"`
	Content            string    `json:"content" bun:"content"`
	SolvingContent     string    `json:"solving_content" bun:"solving_content"`
	CreatedBy          string    `json:"created_by" bun:"created_by,type: char(36),nullzero"`
	UpdatedBy          string    `json:"updated_by" bun:"updated_by,type: char(36),nullzero"`
	CreatedAt          time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt          time.Time `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
}

type TicketPendingUser struct {
	DomainUuid  string `json:"domain_uuid"`
	UserUuid    string `json:"user_uuid"`
	Email       string `json:"email"`
	ZaloId      string `json:"zalo_id"`
	FacebookId  string `json:"facebook_id"`
	ChatInAppId string `json:"chat_in_app_id"`
}
