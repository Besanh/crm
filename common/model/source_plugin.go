package model

import (
	"time"

	"github.com/uptrace/bun"
)

type SourcePlugin struct {
	bun.BaseModel `bun:"source_plugin,alias:sp"`
	DomainUuid    string    `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	SourceUuid    string    `json:"source_uuid" bun:"source_uuid,type: char(36),pk,notnull"`
	SourceName    string    `json:"source_name" bun:"source_name,type:text"`
	Status        bool      `json:"status" bun:"status,type:boolean"`
	CreatedBy     string    `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy     string    `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
}
