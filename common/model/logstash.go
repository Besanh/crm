package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Logstash struct {
	bun.BaseModel  `bun:"logstash,alias:logstash"`
	DomainUuid     string    `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	LogstashUuid   string    `json:"logstash_uuid" bun:"logstash_uuid,pk,type: char(36),notnull"`
	EntityUuid     string    `json:"entity_uuid" bun:"entity_uuid,type:text"`
	Entity         string    `json:"entity" bun:"entity,type:text"`
	EntityStatus   string    `json:"entity_status" bun:"entity_status,type:text"`
	LogstashAction string    `json:"logstash_action" bun:"logstash_action,type:text"` //add update delete import export
	LogstashType   string    `json:"logstash_type" bun:"logstash_type,type:text"`
	OldData        any       `json:"old_data" bun:"old_data,type:text"`
	NewData        any       `json:"new_data" bun:"new_data,type:text"`
	CreatedBy      string    `json:"created_by" bun:"created_by,type:text"`
	CreatedAt      time.Time `json:"created_at" bun:"created_at,type:timestamp,nullzero,notnull,default:current_timestamp"`
}
