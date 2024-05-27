package model

import (
	"time"

	"github.com/uptrace/bun"
)

type EventCalendarAttachment struct {
	bun.BaseModel `bun:"event_calendar_attachment,alias:eca"`
	DomainUuid    string    `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull"`
	EcUuid        string    `json:"ec_uuid" bun:"ec_uuid,type:char(36),notnull"`
	EcaUuid       string    `json:"eca_uuid" bun:"eca_uuid,type:char(36),pk,notnull"`
	FileName      string    `json:"file_name" bun:"file_name,type:text,notnull"`
	PathFile      string    `json:"path_file" bun:"path_file,type:text,notnull"`
	CreatedBy     string    `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy     string    `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
}
