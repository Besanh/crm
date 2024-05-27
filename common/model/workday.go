package model

import (
	"time"

	"github.com/uptrace/bun"
)

type WorkDay struct {
	bun.BaseModel `bun:"workday,alias:wd"`
	DomainUuid    string        `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	WorkdayUuid   string        `json:"workday_uuid" bun:"workday_uuid,type: char(36),pk,notnull"`
	WorkDayId     string        `json:"workday_id" bun:"workday_id,type: text,notnull"`
	WorkDayName   string        `json:"workday_name" bun:"workday_name,type:text"`
	Status        bool          `json:"status" bun:"status,type:boolean,nullzero,default:'false'"`
	Day           string        `json:"day" bun:"day,type:text"`
	StartTime     time.Duration `json:"start_time" bun:"start_time,type:bigint"`
	EndTime       time.Duration `json:"end_time" bun:"end_time,type:bigint"`
	WorkdayType   string        `json:"workday_type" bun:"workday_type,type:workday_type,default:'work'"`
	IsWork        bool          `json:"is_work" bun:"is_work,type:boolean,nullzero,default:'false'"`
	Offset        int           `json:"offset" bun:"offset,type:integer,default:0"`
	Description   string        `json:"description" bun:"description,type:text"`
	UnitUuid      string        `json:"unit_uuid" bun:"unit_uuid,type: char(36)"`
	Unit          *Unit         `json:"unit" bun:"rel:has-one,join:unit_uuid=unit_uuid"`
	CreatedBy     string        `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy     string        `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt     time.Time     `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time     `json:"updated_at" bun:"updated_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
}

type Holiday struct {
	DomainUuid  string    `json:"domain_uuid" bun:"domain_uuid"`
	WorkdayUuid string    `json:"workday_uuid" bun:"workday_uuid"`
	Status      bool      `json:"status" bun:"status"`
	Day         string    `json:"day" bun:"day"`
	StartTime   time.Time `json:"start_time" bun:"start_time"`
	EndTime     time.Time `json:"end_time" bun:"end_time"`
	WorkdayType string    `json:"workday_type" bun:"workday_type"`
	IsWork      bool      `json:"is_work" bun:"is_work"`
	Offset      int       `json:"offset" bun:"offset"`
	Description string    `json:"description" bun:"description"`
	UnitUuid    string    `json:"unit_uuid" bun:"unit_uuid"`
	CreatedBy   string    `json:"created_by" bun:"created_by"`
	UpdatedBy   string    `json:"updated_by" bun:"updated_by"`
	CreatedAt   time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bun:"updated_at"`
}
