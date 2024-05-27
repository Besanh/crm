package model

import (
	"time"

	"github.com/uptrace/bun"
)

type EventCalendarTodo struct {
	bun.BaseModel `bun:"event_calendar_todo,alias:ect"`
	DomainUuid    string    `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull"`
	EcUuid        string    `json:"ec_uuid" bun:"ec_uuid,type:char(36),notnull"`
	EctUuid       string    `json:"ect_uuid" bun:"ect_uuid,type:char(36),pk,notnull"`
	Content       string    `json:"content" bun:"content,type:text,notnull"`
	IsDone        bool      `json:"is_done" bun:"is_done,type:boolean,default:false"`
	CreatedBy     string    `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy     string    `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
}
