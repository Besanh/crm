package omni

import (
	"time"

	"github.com/uptrace/bun"
)

type Omni struct {
	bun.BaseModel `bun:"omni_tmp,alias:omni"`
	DomainUuid    string    `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	OmniUuid      string    `json:"omni_uuid" bun:"omni_uuid,type: char(36),pk,notnull"`
	OmniName      string    `json:"omni_name" bun:"omni_name,type:text,notnull"`
	OmniType      string    `json:"omni_type" bun:"omni_type,type:text,notnull"`
	Supplier      string    `json:"supplier" bun:"supplier,type:text,notnull"`
	Logo          string    `json:"logo" bun:"logo,type:text"`
	Status        bool      `json:"status" bun:"status,type:boolean,default:'true'"`
	Config        string    `json:"config" bun:"config,type:text"`
	CreatedBy     string    `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy     string    `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
}

type Social struct {
	Twitter   string `json:"twitter"`
	Linkedin  string `json:"linkedin"`
	Skype     string `json:"skype"`
	Google    string `json:"google"`
	Facebook  string `json:"facebook"`
	Zalo      string `json:"zalo"`
	Instagram string `json:"instagram"`
	Gmail     string `json:"gmail"`
	Youtube   string `json:"youtube"`
	Tiktok    string `json:"tiktok"`
	Telegram  string `json:"telegram"`
	Website   string `json:"website"`
}
