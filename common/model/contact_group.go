package model

import (
	"time"

	"github.com/uptrace/bun"
)

type ContactGroup struct {
	bun.BaseModel     `bun:"contact_group,alias:cg"`
	DomainUuid        string              `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	ContactGroupUuid  string              `json:"contact_group_uuid" bun:"contact_group_uuid,type: char(36),pk,notnull"`
	GroupType         string              `json:"group_type" bun:"group_type,type:text,notnull"` // customer, member
	GroupName         string              `json:"group_name" bun:"group_name,type:text,notnull"`
	Description       string              `json:"description" bun:"description,type:text"`
	Status            bool                `json:"status" bun:"status,type:bool,nullzero,default:false"`
	CreatedBy         string              `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy         string              `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt         time.Time           `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt         time.Time           `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
	ContactGroupUsers []*ContactGroupUser `json:"contact_group_users" bun:"rel:has-many,join:contact_group_uuid=contact_group_uuid"`
	Member            []string            `json:"member" bun:"-"`
	Staff             []string            `json:"staff" bun:"-"`
}

type ContactGroupUser struct {
	bun.BaseModel        `bun:"contact_group_user,alias:cgu"`
	DomainUuid           string `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	ContactGroupUserUuid string `json:"contact_group_user_uuid" bun:"contact_group_user_uuid,type: char(36),pk,notnull"`
	ContactGroupUuid     string `json:"contact_group_uuid" bun:"contact_group_uuid,type: char(36),notnull"`
	EntityType           string `json:"entity_type" bun:"entity_type,type:text,notnull"`
	UserUuid             string `json:"user_uuid" bun:"user_uuid,type: char(36),notnull"`
	User                 *User  `json:"user" bun:"rel:has-one,join:user_uuid=user_uuid"`
}
