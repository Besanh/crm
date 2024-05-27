package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Contact struct {
	bun.BaseModel `bun:"contact,alias:c"`
	DomainUuid    string     `json:"domain_uuid" bun:"domain_uuid,type:uuid,notnull"`
	ContactUuid   string     `json:"contact_uuid" bun:"contact_uuid,type:uuid,pk,notnull"`
	SourceUuid    string     `json:"source_uuid" bun:"source_uuid,type:text"`
	SourceName    string     `json:"source_name" bun:"source_name,type:varchar(255),default:NULL"`
	Status        bool       `json:"status" bun:"status"`
	UnitUuid      string     `json:"unit_uuid" bun:"unit_uuid,type: char(36)"`
	Unit          *Unit      `json:"unit" bun:"rel:belongs-to,join:unit_uuid=unit_uuid"`
	ContactType   string     `json:"contact_type" bun:"contact_type,type:varchar(255)"`
	ContactName   string     `json:"contact_name" bun:"contact_name,type:varchar(255)"`
	Profiles      []*Profile `json:"profiles" bun:"rel:has-many,join:contact_uuid=contact_uuid"`
	CreatedBy     string     `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy     string     `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt     time.Time  `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time  `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
}

type ContactView struct {
	bun.BaseModel `bun:"contact,alias:c"`
	DomainUuid    string        `json:"domain_uuid" bun:"domain_uuid"`
	ContactUuid   string        `json:"contact_uuid" bun:"contact_uuid,pk"`
	SourceUuid    string        `json:"source_uuid" bun:"source_uuid"`
	SourceName    string        `json:"source_name" bun:"source_name"`
	Status        bool          `json:"status" bun:"status"`
	UnitUuid      string        `json:"unit_uuid" bun:"unit_uuid"`
	Unit          *Unit         `json:"unit" bun:"rel:belongs-to,join:unit_uuid=unit_uuid"`
	ContactType   string        `json:"contact_type" bun:"contact_type"`
	ContactName   string        `json:"contact_name" bun:"contact_name"`
	Profiles      []ProfileView `json:"profiles" bun:"rel:has-many,join:contact_uuid=contact_uuid"`
	CreatedBy     string        `json:"created_by" bun:"created_by"`
	UpdatedBy     string        `json:"updated_by" bun:"updated_by"`
	CreatedAt     time.Time     `json:"created_at" bun:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" bun:"updated_at"`
}

type ContactInfo struct {
	ContactUuid string        `json:"contact_uuid"`
	Status      bool          `json:"status"`
	Profiles    []ProfileView `json:"profiles"`
	UnitUuid    string        `json:"unit_uuid"`
	ContactType string        `json:"contact_type"`
	ContactName string        `json:"contact_name"`
}

type ContactNotePost struct {
	Phone       string `json:"phone"`
	ContactUuid string `json:"contact_uuid"`
	Content     string `json:"content"`
	Type        string `json:"type"`
	NoteBy      string `json:"note_by"`
}

type ContactImportResult struct {
	TotalSuccess int `json:"total_success"`
	TotalFail    int `json:"total_fail"`
	Total        int `json:"total"`
}

type ContactOwner struct {
	bun.BaseModel `bun:"contact_owner,alias:co"`
	ContactUuid   string    `json:"contact_uuid" bun:"contact_uuid,type:char(36),unique:unqx_contact_owner"`
	DomainUuid    string    `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull"`
	UserUuid      string    `json:"user_uuid" bun:"user_uuid,type:char(36),notnull,unique:unqx_contact_owner"`
	Username      string    `json:"username" bun:"username,type:text,notnull"`
	Type          string    `json:"type" bun:"type,type:text,notnull"`
	AssignedAt    time.Time `json:"assigned_at" bun:"assigned_at,nullzero"`
}

type ContactChannel struct {
	bun.BaseModel `bun:"contact_channel,alias:cc"`
	ContactUuid   string `json:"contact_uuid" bun:"contact_uuid,type:char(36)"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull"`
	Channel       string `json:"channel" bun:"channel,type:text,notnull"`
	Vendor        string `json:"vendor" bun:"vendor,type:text"`
	Data          string `json:"data" bun:"data,type:text"`
	PageId        string `json:"page_id" bun:"page_id,type:text"`
}

type ContactNote struct {
	bun.BaseModel `bun:"contact_note,alias:cn"`
	ContactUuid   string    `json:"contact_uuid" bun:"contact_uuid,type:char(36)"`
	DomainUuid    string    `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull"`
	Type          string    `json:"type" bun:"type,type:text,notnull"`
	Status        bool      `json:"status" bun:"status,type:boolean,nullzero,default:true"`
	Content       string    `json:"content" bun:"content,type:text"`
	UserUuid      string    `json:"user_uuid" bun:"user_uuid,type:char(36),notnull"`
	CreatedBy     string    `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy     string    `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
}

type ContactPhone struct {
	bun.BaseModel `bun:"contact_phone,alias:cp"`
	ContactUuid   string `json:"contact_uuid" bun:"contact_uuid,type:char(36),unique:unqx_contact_phone"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull"`
	Data          string `json:"data" bun:"data,type:text,unique:unqx_contact_phone"`
	Type          string `json:"type" bun:"type,type:text,notnull"`
}

type ContactEmail struct {
	bun.BaseModel `bun:"contact_email,alias:ce"`
	ContactUuid   string `json:"contact_uuid" bun:"contact_uuid,type:char(36),unique:unqx_contact_email"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull"`
	Data          string `json:"data" bun:"data,type:text,unique:unqx_contact_email"`
	EmailType     string `json:"email_type" bun:"email_type,type:text,notnull"`
	EmailAddress  string `json:"email_address" bun:"email_address,type:text"`
	EmailPrimary  int32  `json:"email_primary" bun:"email_primary,type:integer,notnull"`
}

type ContactPost struct {
	SourceUuid  string `json:"source_uuid"`
	SourceName  string `json:"source_name"`
	Status      bool   `json:"status"`
	UnitUuid    string `json:"unit_uuid"`
	ContactType string `json:"contact_type"`
	ContactName string `json:"contact_name"`
}

type ContactMapData struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

type ContactNoteData struct {
	bun.BaseModel `bun:"contact_note,alias:cn"`
	ContactUuid   string    `json:"-" bun:"contact_uuid"`
	DomainUuid    string    `json:"-" bun:"domain_uuid"`
	Type          string    `json:"type" bun:"type"`
	Status        bool      `json:"status" bun:"status"`
	Content       string    `json:"content" bun:"content"`
	UserUuid      string    `json:"user_uuid" bun:"user_uuid"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at"`
	Username      string    `json:"username" bun:"username"`
	LastName      string    `json:"last_name" bun:"last_name"`
	FirstName     string    `json:"first_name" bun:"first_name"`
	MiddleName    string    `json:"middle_name" bun:"middle_name"`
}

type ContactMoreInformationData struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}
