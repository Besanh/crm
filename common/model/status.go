package model

import (
	"time"

	"github.com/uptrace/bun"
)

type CategoryStatus struct {
	bun.BaseModel      `bun:"category_status"`
	DomainUuid         string    `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull,unique:unq_category_code_domain"`
	Id                 string    `json:"category_status_uuid" bun:"category_status_uuid,pk,type:char(36)"`
	CategoryStatusName string    `json:"category_status_name" bun:"category_status_name,type:varchar(50),notnull"`
	CategoryStatusCode string    `json:"category_status_code" bun:"category_status_code,type:varchar(20),notnull,unique:unq_category_code_domain"`
	Description        string    `json:"description" bun:"description,type:text,nullzero"`
	CreatedAt          time.Time `json:"created_at" bun:"created_at,type:timestamp,nullzero,notnull,default:current_timestamp"`
	UpdatedAt          time.Time `json:"updated_at" bun:"updated_at,type:timestamp,nullzero,notnull,default:current_timestamp"`
}

type CategoryStatusView struct {
	bun.BaseModel      `bun:"category_status"`
	DomainUuid         string    `json:"domain_uuid" bun:"domain_uuid"`
	Id                 string    `json:"category_status_uuid" bun:"category_status_uuid,pk"`
	CategoryStatusName string    `json:"category_status_name" bun:"category_status_name"`
	CategoryStatusCode string    `json:"category_status_code" bun:"category_status_code"`
	Description        string    `json:"description" bun:"description"`
	CreatedAt          time.Time `json:"created_at" bun:"created_at"`
}

type Status struct {
	bun.BaseModel      `bun:"status,alias:status"`
	DomainUuid         string    `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull,unique:unq_status_code_domain_campaign"`
	CampaignUuid       string    `json:"campaign_uuid" bun:"campaign_uuid,type:char(36),notnull,unique:unq_status_code_domain_campaign"`
	Id                 string    `json:"status_uuid" bun:"status_uuid,pk,type:char(36),pk"`
	StatusCode         string    `json:"status_code" bun:"status_code,type:varchar(20),notnull,unique:unq_status_code_domain_campaign"`
	StatusName         string    `json:"status_name" bun:"status_name,type:varchar(200),notnull"`
	CategoryStatusUuid string    `json:"category_status_uuid" bun:"category_status_uuid,type:char(36),notnull,unique:unq_status_code_domain_campaign"`
	Selectable         bool      `json:"selectable" bun:"selectable,type:bool,default:'0',notnull"`
	Sale               bool      `json:"sale" bun:"sale,type:bool,default:'0',notnull"`
	Dnc                bool      `json:"dnc" bun:"dnc,type:bool,default:'0',notnull"`
	NotInterested      bool      `json:"not_interested" bun:"not_interested,type:bool,default:'0',notnull"`
	ScheduledCallback  bool      `json:"scheduled_callback" bun:"scheduled_callback,type:bool,default:'0',notnull"`
	CreatedAt          time.Time `json:"created_at" bun:"created_at,type:timestamp,nullzero,notnull,default:current_timestamp"`
	UpdatedAt          time.Time `json:"updated_at" bun:"updated_at,type:timestamp,nullzero,notnull,default:current_timestamp"`
}

type StatusUpdate struct {
	bun.BaseModel      `bun:"status,alias:status"`
	DomainUuid         string    `json:"domain_uuid" bun:"domain_uuid"`
	CampaignUuid       string    `json:"campaign_uuid" bun:"campaign_uuid"`
	Id                 string    `json:"status_uuid" bun:"status_uuid,pk,type:char(36),pk"`
	StatusCode         string    `json:"status_code" bun:"status_code,type:varchar(20)"`
	StatusName         string    `json:"status_name" bun:"status_name,type:varchar(200)"`
	CategoryStatusUuid string    `json:"category_status_uuid" bun:"category_status_uuid"`
	Selectable         bool      `json:"selectable" bun:"selectable"`
	Sale               bool      `json:"sale" bun:"sale"`
	Dnc                bool      `json:"dnc" bun:"dnc"`
	NotInterested      bool      `json:"not_interested" bun:"not_interested"`
	ScheduledCallback  bool      `json:"scheduled_callback" bun:"scheduled_callback"`
	UpdatedAt          time.Time `json:"updated_at" bun:"updated_at"`
}

type StatusView struct {
	bun.BaseModel      `bun:"status,alias:status"`
	DomainUuid         string `json:"domain_uuid" bun:"domain_uuid"`
	CampaignUuid       string `json:"campaign_uuid" bun:"campaign_uuid"`
	Id                 string `json:"status_uuid" bun:"status_uuid,pk"`
	StatusCode         string `json:"status_code" bun:"status_code"`
	StatusName         string `json:"status_name" bun:"status_name"`
	CategoryStatusUuid string `json:"category_status_uuid" bun:"category_status_uuid"`
	Selectable         bool   `json:"selectable" bun:"selectable"`
	Sale               bool   `json:"sale" bun:"sale"`
	Dnc                bool   `json:"dnc" bun:"dnc"`
	NotInterested      bool   `json:"not_interested" bun:"not_interested"`
	ScheduledCallback  bool   `json:"scheduled_callback" bun:"scheduled_callback"`
	CategoryStatusName string `json:"category_status_name" bun:"-"`
	CampaignName       string `json:"campaign_name" bun:"-"`
}

type StatusCampaignView struct {
	bun.BaseModel      `bun:"status,alias:s"`
	CampaignUuid       string `json:"campaign_uuid" bun:"campaign_uuid"`
	Id                 string `json:"status_uuid" bun:"status_uuid,pk"`
	StatusCode         string `json:"status_code" bun:"status_code"`
	StatusName         string `json:"status_name" bun:"status_name"`
	CategoryStatusUuid string `json:"category_status_uuid" bun:"category_status_uuid"`
	Selectable         bool   `json:"selectable" bun:"selectable"`
	Sale               bool   `json:"sale" bun:"sale"`
	Dnc                bool   `json:"dnc" bun:"dnc"`
	NotInterested      bool   `json:"not_interested" bun:"not_interested"`
	ScheduledCallback  bool   `json:"scheduled_callback" bun:"scheduled_callback"`
	CategoryStatusName string `json:"category_status_name" bun:"category_status_name"`
	CategoryStatusCode string `json:"category_status_code" bun:"category_status_code"`
}
