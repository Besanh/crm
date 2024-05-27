package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Pbx struct {
	bun.BaseModel `bun:"pbx,alias:pbx"`
	DomainUuid    string    `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	PbxUuid       string    `json:"pbx_uuid" bun:"pbx_uuid,pk,type: char(36),pk,notnull"`
	Supplier      string    `json:"supplier" bun:"supplier,type:text,notnull"`
	PbxName       string    `json:"pbx_name" bun:"pbx_name,type:text,notnull"`
	Status        bool      `json:"status" bun:"status,type:boolean,nullzero,default:false"`
	Verified      bool      `json:"verified" bun:"verified,type:boolean,nullzero,default:false"`
	Domain        string    `json:"domain" bun:"domain,type:text"`
	ApiKey        string    `json:"api_key" bun:"api_key,type:text"`
	OutboundProxy string    `json:"outbound_proxy" bun:"outbound_proxy,type:text"`
	Wss           string    `json:"wss" bun:"wss,type:text"`
	Transport     string    `json:"transport" bun:"transport,type:text"`
	Port          string    `json:"port" bun:"port,type:text"`
	UrlCall       string    `json:"url_call" bun:"url_call,type:text"`
	Webhook       []string  `json:"webhook" bun:"webhook,type:text"`
	UnitUuid      string    `json:"unit_uuid" bun:"unit_uuid,type:text"`
	Unit          *Unit     `json:"unit" bun:"rel:has-one,join:unit_uuid=unit_uuid"`
	CreatedBy     string    `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy     string    `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
}
