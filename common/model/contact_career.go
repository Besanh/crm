package model

import (
	"time"

	"github.com/uptrace/bun"
)

type ContactCareer struct {
	bun.BaseModel      `bun:"contact_career,alias:cc"`
	DomainUuid         string               `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	ContactCareerUuid  string               `json:"contact_career_uuid" bun:"contact_career_uuid,type: char(36),pk,notnull"`
	CareerType         string               `json:"career_type" bun:"career_type,type:text,notnull"` // career
	CareerName         string               `json:"career_name" bun:"career_name,type:text,notnull"`
	Description        string               `json:"description" bun:"description,type:text"`
	Status             bool                 `json:"status" bun:"status,type:bool,nullzero,default:false"`
	CreatedBy          string               `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy          string               `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt          time.Time            `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt          time.Time            `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
	ContactcareerUsers []*ContactCareerUser `json:"contact_career_users" bun:"rel:has-many,join:contact_career_uuid=contact_career_uuid"`
	Career             []string             `json:"career" bun:"-"`
}

type ContactCareerUser struct {
	bun.BaseModel         `bun:"contact_career_user,alias:ccu"`
	DomainUuid            string `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	ContactCareerUserUuid string `json:"contact_career_user_uuid" bun:"contact_career_user_uuid,type: char(36),pk,notnull"`
	ContactCareerUuid     string `json:"contact_career_uuid" bun:"contact_career_uuid,type: char(36),notnull"`
	EntityType            string `json:"entity_type" bun:"entity_type,type:text,notnull"`
	UserUuid              string `json:"user_uuid" bun:"user_uuid,type: char(36),notnull"`
	User                  *User  `json:"user" bun:"rel:has-one,join:user_uuid=user_uuid"`
}
