package model

import (
	"contactcenter-api/common/model/omni"
	"time"

	"github.com/uptrace/bun"
)

type Unit struct {
	bun.BaseModel       `bun:"unit,alias:unit"`
	DomainUuid          string      `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	ParentUnitUuid      string      `json:"parent_unit_uuid" bun:"parent_unit_uuid,type: text"`
	UnitUuid            string      `json:"unit_uuid" bun:"unit_uuid,pk,type: char(36),notnull"`
	UnitName            string      `json:"unit_name" bun:"unit_name,type:text,notnull"`
	UnitCode            string      `json:"unit_code" bun:"unit_code,type:text,unique"`
	UnitLeader          string      `json:"unit_leader" bun:"unit_leader,type:text"`
	UnitBasis           bool        `json:"unit_basis" bun:"unit_basis,type:boolean,nullzero,default:false"`
	Status              bool        `json:"status" bun:"status,type:boolean"`
	Level               string      `json:"level" bun:"level,type:text"`
	BrandName           string      `json:"brand_name" bun:"brand_name,type:text"`
	PhoneNumber         string      `json:"phone_number" bun:"phone_number,type:text"`
	Email               string      `json:"email" bun:"email,type:text"`
	Country             string      `json:"country" bun:"country,type:text"`
	City                string      `json:"city" bun:"city,type:text"`
	TaxCode             string      `json:"tax_code" bun:"tax_code,type:text"`
	DateOfIncorporation string      `json:"date_of_incorporation" bun:"date_of_incorporation,type:text"`
	Address             string      `json:"address" bun:"address,type:text"`
	BankName            string      `json:"bank_name" bun:"bank_name,type:text"`
	AccountNumber       string      `json:"account_number" bun:"account_number,type:text"`
	Note                string      `json:"note" bun:"note,type:text"`
	UnitSocial          omni.Social `json:"unit_social" bun:"unit_social,type:text"`
	UnitConfig          UnitConfig  `json:"unit_config" bun:"unit_config,type:text"`
	CreatedBy           string      `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy           string      `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt           time.Time   `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt           time.Time   `json:"updated_at" bun:"updated_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
}

type UnitInfo struct {
	bun.BaseModel       `bun:"unit,alias:unit"`
	DomainUuid          string      `json:"domain_uuid" bun:"domain_uuid"`
	ParentUnitUuid      string      `json:"parent_unit_uuid" bun:"parent_unit_uuid"`
	ParentUnitCode      string      `json:"parent_unit_code" bun:"-"`
	UnitUuid            string      `json:"unit_uuid" bun:"unit_uuid,pk"`
	UnitName            string      `json:"unit_name" bun:"unit_name"`
	UnitCode            string      `json:"unit_code" bun:"unit_code"`
	UnitLeader          string      `json:"unit_leader" bun:"unit_leader"`
	UnitBasis           bool        `json:"unit_basis" bun:"unit_basis"`
	Status              bool        `json:"status" bun:"status"`
	Level               string      `json:"level" bun:"level"`
	BrandName           string      `json:"brand_name" bun:"brand_name"`
	PhoneNumber         string      `json:"phone_number" bun:"phone_number"`
	Email               string      `json:"email" bun:"email"`
	Country             string      `json:"country" bun:"country"`
	City                string      `json:"city" bun:"city"`
	TaxCode             string      `json:"tax_code" bun:"tax_code"`
	DateOfIncorporation string      `json:"date_of_incorporation" bun:"date_of_incorporation"`
	Address             string      `json:"address" bun:"address"`
	BankName            string      `json:"bank_name" bun:"bank_name"`
	AccountNumber       string      `json:"account_number" bun:"account_number"`
	Note                string      `json:"note" bun:"note"`
	UnitSocial          omni.Social `json:"unit_social" bun:"unit_social"`
	UnitConfig          UnitConfig  `json:"unit_config" bun:"unit_config"`
	Quantity            int         `json:"quantity" bun:"-"`
	CreatedBy           string      `json:"created_by" bun:"created_by"`
	UpdatedBy           string      `json:"updated_by" bun:"updated_by"`
	CreatedAt           time.Time   `json:"created_at" bun:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at" bun:"updated_at"`
	Users               []*User     `json:"users" bun:"rel:has-many,join:unit_uuid=unit_uuid"`
}

type UnitConfig struct {
	Version      string `json:"version"`
	Logo         string `json:"logo"`
	ApiUrl       string `json:"api_url"`
	ZaloApiKey   string `json:"zalo_api_key"`
	GoogleApiKey string `json:"google_api_key"`
	Partner      string `json:"partner"`
	IsRecording  bool   `json:"is_recording"`
	PaymentCost  string `json:"payment_cost"`
}
