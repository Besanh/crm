package model

import (
	"time"

	"github.com/uptrace/bun"
)

type ClassifyCareer struct {
	bun.BaseModel      `bun:"classify_career,alias:cc"`
	DomainUuid         string    `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	ClassifyCareerUuid string    `json:"classify_career_uuid" bun:"classify_career_uuid,type: char(36),pk,notnull"`
	CareerType         string    `json:"career_type" bun:"career_type,type:text,notnull"` // career
	CareerName         string    `json:"career_name" bun:"career_name,type:text,notnull"`
	Description        string    `json:"description" bun:"description,type:text"`
	Status             bool      `json:"status" bun:"status,type:bool,nullzero,default:false"`
	CreatedBy          string    `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy          string    `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt          time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt          time.Time `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
	Career             []string  `json:"career" bun:"-"`
}
