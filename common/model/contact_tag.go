package model

import (
	"time"

	"github.com/uptrace/bun"
)

type ContactTag struct {
	bun.BaseModel   `bun:"contact_tag,alias:ct"`
	DomainUuid      string            `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	ContactTagUuid  string            `json:"contact_tag_uuid" bun:"contact_tag_uuid,type: char(36),pk,notnull"`
	TagName         string            `json:"tag_name" bun:"tag_name,type:text,notnull"`
	LimitedFunction []string          `json:"limited_function" bun:"limited_function,type:text[]"`
	Description     string            `json:"description" bun:"description,type:text"`
	ContactTagUsers []*ContactTagUser `json:"contact_tag_users" bun:"rel:has-many,join:contact_tag_uuid=contact_tag_uuid"`
	Status          bool              `json:"status" bun:"status,type:bool,nullzero,default:false"`
	CreatedBy       string            `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy       string            `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt       time.Time         `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt       time.Time         `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
	Member          []string          `json:"member" bun:"-"`
	Staff           []string          `json:"staff" bun:"-"`
}

type ContactTagUser struct {
	bun.BaseModel      `bun:"contact_tag_user,alias:ctu"`
	DomainUuid         string `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	ContactTagUserUuid string `json:"contact_tag_user_uuid" bun:"contact_tag_user_uuid,type: char(36),pk,notnull"`
	ContactTagUuid     string `json:"contact_tag_uuid" bun:"contact_tag_uuid,type: char(36),notnull"`
	EntityType         string `json:"entity_type" bun:"entity_type,type:text,notnull"`
	UserUuid           string `json:"user_uuid" bun:"user_uuid,type: char(36),notnull"`
	User               *User  `json:"user" bun:"rel:has-one,join:user_uuid=user_uuid"`
}
