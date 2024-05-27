package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Lead struct {
	bun.BaseModel    `bun:"lead"`
	DomainUuid       string    `json:"-" bun:"domain_uuid,type:uuid,notnull"`
	LeadUuid         string    `json:"lead_uuid" bun:"lead_uuid,pk,type:uuid,notnull"`
	LeadName         string    `json:"lead_name" bun:"lead_name,nullzero,type:text"`
	LeadCode         string    `json:"lead_code" bun:"lead_code,type:text,nullzero"`
	ListUuid         string    `json:"list_uuid" bun:"list_uuid,type:uuid,notnull"`
	PhoneNumber      string    `json:"phone_number" bun:"phone_number,type:text,notnull"`
	ContractNumber   string    `json:"contract_number" bun:"contract_number,type:text"`
	Status           string    `json:"status" bun:"status,type:text,notnull,default:'NEW'"`
	AltStatus        string    `json:"alt_status" bun:"alt_status,type:text,notnull,default:'new'"`
	IdentityNumber   string    `json:"identity_number" bun:"identity_number,type:text"`
	IdentityIssuedAt string    `json:"identity_issued_at" bun:"identity_issued_at,type:text"`
	IdentityIssuedOn time.Time `json:"identity_issued_on" bun:"identity_issued_on,type:timestamp,nullzero"`
	Address          string    `json:"address" bun:"address,type:text"`
	Province         string    `json:"province" bun:"province,type:text"`
	District         string    `json:"district" bun:"district,type:text"`
	Ward             string    `json:"ward" bun:"ward,type:text"`
	DateOfBirth      string    `json:"date_of_birth" bun:"date_of_birth,type:text"`
	Additional       string    `json:"additional" bun:"additional,nullzero,type:text"`
	AfterDay         int       `json:"after_day" bun:"after_day,type:integer"`
	CalledCount      int       `json:"called_count" bun:"called_count,type:integer,notnull"`
	LastCallTime     time.Time `json:"last_call_time" bun:"last_call_time,type:timestamp,default:'1970-01-01 08:00:01',nullzero"`
	CreatedAt        time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt        time.Time `json:"updated_at" bun:"updated_at,type:timestamp,notnull,nullzero,default:current_timestamp"`

	// NEW
	Comment               string    `json:"comment" bun:"comment,type:text"`
	RefId                 string    `json:"ref_id" bun:"ref_id,type:text"`
	RefCode               string    `json:"ref_code" bun:"ref_code,type:text"`
	Gender                string    `json:"gender" bun:"gender,type:text"`
	OtherIdentityNumber   string    `json:"other_identity_number" bun:"other_identity_number,type:text"`
	OtherIdentityIssuedAt string    `json:"other_identity_issued_at" bun:"other_identity_issued_at,type:text"`
	OtherIdentityIssuedOn time.Time `json:"other_identity_issued_on" bun:"other_identity_issued_on,type:timestamp,nullzero"`

	AssignUserUuid   string `json:"assign_user_uuid" bun:"assign_user_uuid,type:uuid,nullzero,default:NULL"`
	AssignerUserUuid string `json:"assigner_user_uuid" bun:"assigner_user_uuid,type:uuid,nullzero,default:NULL"`

	Network string `json:"network" bun:"network,type:text"`

	ExpiredAt time.Time `json:"expired_at" bun:"expired_at,type:timestamp,nullzero,default:NULL"`
	IsExpired bool      `json:"is_expired" bun:"is_expired,type:bool,default:false"`

	// Crm
	RelatedProfileType string `json:"related_profile_type" bun:"related_profile_type,type:text,nullzero,default:NULL"`
	RelatedProfileUuid string `json:"related_profile_uuid" bun:"related_profile_uuid,type:uuid,nullzero,default:NULL"`
	Country            string `json:"country" bun:"country,type:text"`
	JobTitle           string `json:"job_title" bun:"job_title,type:text"`
	IsConvertProfile   bool   `json:"is_convert_profile" bun:"is_convert_profile,type:bool,nullzero,default:false"`
}
