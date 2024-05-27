package model

import (
	"time"

	"github.com/uptrace/bun"
)

type ManagementScope struct {
	bun.BaseModel       `bun:"management_scope,alias:ms"`
	DomainUuid          string    `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	ManagementScopeUuid string    `json:"management_scope_uuid" bun:"management_scope_uuid,pk,type: char(36),notnull"`
	ManagementScopeName string    `json:"management_scope_name" bun:"management_scope_name,type:text,notnull"`
	Status              bool      `json:"status" bun:"status,type:boolean"`
	CreatedBy           string    `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy           string    `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt           time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt           time.Time `json:"updated_at" bun:"updated_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
}
