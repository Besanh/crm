package model

import (
	"time"

	"github.com/uptrace/bun"
)

type ClassifyGroup struct {
	bun.BaseModel     `bun:"classify_group,alias:cg"`
	DomainUuid        string    `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	ClassifyGroupUuid string    `json:"classify_group_uuid" bun:"classify_group_uuid,type: char(36),pk,notnull"`
	GroupType         string    `json:"group_type" bun:"group_type,type:text,notnull"` // customer, member
	GroupName         string    `json:"group_name" bun:"group_name,type:text,notnull"`
	Description       string    `json:"description" bun:"description,type:text"`
	Status            bool      `json:"status" bun:"status,type:bool,nullzero,default:false"`
	CreatedBy         string    `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy         string    `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt         time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt         time.Time `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
	Member            []string  `json:"member" bun:"-"`
	Staff             []string  `json:"staff" bun:"-"`
}
