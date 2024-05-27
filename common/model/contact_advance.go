package model

import (
	"time"

	"github.com/uptrace/bun"
)

// Contact relation to tag
type ContactToTag struct {
	bun.BaseModel    `bun:"contact_to_tag,alias:ctt"`
	DomainUuid       string      `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	ContactToTagUuid string      `json:"contact_to_tag_uuid" bun:"contact_to_tag_uuid,type: char(36),pk,notnull"`
	ContactUuid      string      `json:"contact_uuid" bun:"contact_uuid,type: char(36),notnull"`
	ContactTagUuid   string      `json:"contact_tag_uuid" bun:"contact_tag_uuid,type: char(36),notnull"`
	ContactTag       *ContactTag `json:"contact_tag" bun:"rel:has-one,join:contact_tag_uuid=contact_tag_uuid"`
	CreatedAt        time.Time   `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt        time.Time   `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
}

// Contact relation to group
type ContactToGroup struct {
	bun.BaseModel      `bun:"contact_to_group,alias:ctg"`
	DomainUuid         string        `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	ContactToGroupUuid string        `json:"contact_to_group_uuid" bun:"contact_to_group_uuid,type: char(36),pk,notnull"`
	ContactUuid        string        `json:"contact_uuid" bun:"contact_uuid,type: char(36),notnull"`
	ContactGroupUuid   string        `json:"contact_group_uuid" bun:"contact_group_uuid,type: char(36),notnull"`
	ContactGroup       *ContactGroup `json:"contact_group" bun:"rel:has-one,join:contact_group_uuid=contact_group_uuid"`
	CreatedAt          time.Time     `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt          time.Time     `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
}

// Contact relation to career
type ContactToCareer struct {
	bun.BaseModel       `bun:"contact_to_career,alias:ctc"`
	DomainUuid          string         `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	ContactToCareerUuid string         `json:"contact_to_career_uuid" bun:"contact_to_career_uuid,type: char(36),pk,notnull"`
	ContactUuid         string         `json:"contact_uuid" bun:"contact_uuid,type: char(36),notnull"`
	ContactCareerUuid   string         `json:"contact_career_uuid" bun:"contact_career_uuid,type: char(36),notnull"`
	ContactCareer       *ContactCareer `json:"contact_career" bun:"rel:has-one,join:contact_career_uuid=contact_career_uuid"`
	CreatedAt           time.Time      `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt           time.Time      `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
}
