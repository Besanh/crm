package model

import (
	"time"

	"github.com/uptrace/bun"
)

type ClassifyTag struct {
	bun.BaseModel   `bun:"classify_tag,alias:ct"`
	DomainUuid      string    `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	ClassifyTagUuid string    `json:"classify_tag_uuid" bun:"classify_tag_uuid,type: char(36),pk,notnull"`
	TagName         string    `json:"tag_name" bun:"tag_name,type:text,notnull"`
	LimitedFunction []string  `json:"limited_function" bun:"limited_function,type:text[]"`
	Description     string    `json:"description" bun:"description,type:text"`
	Status          bool      `json:"status" bun:"status,type:bool,nullzero,default:false"`
	CreatedBy       string    `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy       string    `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt       time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt       time.Time `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
	Member          []string  `json:"member" bun:"-"`
	Staff           []string  `json:"staff" bun:"-"`
}
