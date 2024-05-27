package model

import (
	"time"

	"github.com/uptrace/bun"
)

type EventCalendarCategory struct {
	bun.BaseModel `bun:"event_calendar_category,alias:ecc"`
	DomainUuid    string    `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull"`
	EccUuid       string    `json:"ecc_uuid" bun:"ecc_uuid,type:char(36),pk,notnull"`
	Title         string    `json:"title" bun:"title,type:text,notnull"`
	Color         string    `json:"color" bun:"color,type:text,notnull"`
	Status        bool      `json:"status" bun:"status,,type:boolean,notnull,default:'false'"`
	CreatedBy     string    `json:"created_by" bun:"created_by,type:char(200)"`
	UpdatedBy     string    `json:"updated_by" bun:"updated_by,type:char(200)"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
}
