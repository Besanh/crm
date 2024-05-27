package model

import (
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

type TicketCategory struct {
	bun.BaseModel            `bun:"ticket_category,alias:tc"`
	DomainUuid               string           `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull,unique:unq_ticket_category_domain"`
	TicketCategoryUuid       string           `json:"ticket_category_uuid" bun:"ticket_category_uuid,type: char(36),pk,notnull"`
	TicketCategoryCode       string           `json:"ticket_category_code" bun:"ticket_category_code,type:text,unique,notnull,unique:unq_ticket_category_domain"`
	TicketCategoryName       string           `json:"ticket_category_name" bun:"ticket_category_name,type:text"`
	ParentTicketCategoryUuid string           `json:"parent_ticket_category_uuid" bun:"parent_ticket_category_uuid,type:uuid,default:NULL"`
	Description              string           `json:"description" bun:"description,type:text"`
	Status                   bool             `json:"status" bun:"status,type:boolean"`
	Sla                      float64          `json:"sla" bun:"sla,type:numeric"`
	Level                    int              `json:"level" bun:"level,type:numeric"`
	UnitUuid                 string           `json:"unit_uuid" bun:"unit_uuid,type: char(36),nullzero"`
	CreatedBy                string           `json:"created_by" bun:"created_by,type: char(36),nullzero"`
	UpdatedBy                string           `json:"updated_by" bun:"updated_by,type: char(36),nullzero"`
	CreatedAt                time.Time        `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt                time.Time        `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
	SLAPolicies              []*SlaPolicyInfo `json:"sla_policies" bun:"rel:has-many"`
}

type TicketCategoryInfo struct {
	bun.BaseModel            `bun:"ticket_category,alias:tc"`
	DomainUuid               string           `json:"domain_uuid" bun:"domain_uuid"`
	TicketCategoryUuid       string           `json:"ticket_category_uuid" bun:"ticket_category_uuid,pk"`
	TicketCategoryCode       string           `json:"ticket_category_code" bun:"ticket_category_code"`
	TicketCategoryName       string           `json:"ticket_category_name" bun:"ticket_category_name"`
	ParentTicketCategoryUuid string           `json:"parent_ticket_category_uuid" bun:"parent_ticket_category_uuid"`
	ParentTicketCategoryName string           `json:"parent_ticket_category_name" bun:"parent_ticket_category_name"`
	Description              string           `json:"description" bun:"description"`
	Status                   bool             `json:"status" bun:"status"`
	Sla                      float64          `json:"sla" bun:"sla"`
	Level                    string           `json:"level" bun:"level"`
	UnitUuid                 string           `json:"unit_uuid" bun:"unit_uuid"`
	CreatedBy                string           `json:"created_by" bun:"created_by"`
	UpdatedBy                string           `json:"updated_by" bun:"updated_by"`
	CreatedAt                time.Time        `json:"created_at" bun:"created_at"`
	UpdatedAt                time.Time        `json:"updated_at" bun:"updated_at"`
	IsParent                 bool             `json:"is_parent" bun:"-"`
	SLAPolicies              []*SlaPolicyInfo `json:"sla_policies" bun:"rel:has-many"`
}

type TicketCategoryPost struct {
	DomainUuid               string      `json:"domain_uuid"`
	TicketCategoryUuid       string      `json:"ticket_category_uuid"`
	TicketCategoryCode       string      `json:"ticket_category_code"`
	TicketCategoryName       string      `json:"ticket_category_name" `
	ParentTicketCategoryUuid string      `json:"parent_ticket_category_uuid"`
	ParentTicketCategoryName string      `json:"parent_ticket_category_name" bun:"parent_ticket_category_name"`
	Description              string      `json:"description" bun:"description"`
	Status                   bool        `json:"status" bun:"status"`
	CreatedBy                string      `json:"created_by" bun:"created_by"`
	UpdatedBy                string      `json:"updated_by" bun:"updated_by"`
	CreatedAt                time.Time   `json:"created_at" bun:"created_at"`
	UpdatedAt                time.Time   `json:"updated_at" bun:"updated_at"`
	IsParent                 bool        `json:"is_parent" bun:"is_parent"`
	SLAPolicies              []SlaPolicy `json:"sla_policies"`
	UnitUuid                 string      `json:"unit_uuid" bun:"unit_uuid"`
}

type TicketCategoryFilter struct {
	DomainUuid               string       `json:"domain_uuid"`
	ParentTicketCategoryUuid string       `json:"parent_ticket_category_uuid"`
	TicketCategoryUuid       string       `json:"ticket_category_uuid"`
	TicketCategoryName       string       `json:"ticket_category_name"`
	TicketCategoryCode       string       `json:"ticket_category_code"`
	Active                   sql.NullBool `json:"active"`
	IsParent                 sql.NullBool `json:"is_parent"`
	Level                    string       `json:"level"`
}

type TicketCategoryHierarchy struct {
	TicketCategoryUuid string `json:"ticket_category_uuid"`
}
