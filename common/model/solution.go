package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Solution struct {
	bun.BaseModel `bun:"solution,alias:solution"`
	DomainUuid    string    `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	SolutionUuid  string    `json:"solution_uuid" bun:"solution_uuid,type: char(36),pk,notnull"`
	SolutionCode  string    `json:"solution_code" bun:"solution_code,type:text,unique,notnull"`
	SolutionName  string    `json:"solution_name" bun:"solution_name,type:text"`
	Status        bool      `json:"status" bun:"status,type:boolean,default:'true'"`
	UnitUuid      string    `json:"unit_uuid" bun:"unit_uuid,type:text,notnull"`
	Unit          *Unit     `json:"unit" bun:"rel:has-one,join:unit_uuid=unit_uuid"`
	CreatedBy     string    `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy     string    `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
}

type SolutionPost struct {
	SolutionCode string `json:"solution_code" bun:"solution_code"`
	SolutionName string `json:"solution_name" bun:"solution_name"`
	Status       bool   `json:"status" bun:"status"`
	UnitUuid     string `json:"unit_uuid" bun:"unit_uuid"`
}
