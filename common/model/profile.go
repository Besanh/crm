package model

import (
	"contactcenter-api/common/model/omni"
	"time"

	"github.com/uptrace/bun"
)

type Profile struct {
	bun.BaseModel        `bun:"profile,alias:p"`
	DomainUuid           string               `json:"domain_uuid" bun:"domain_uuid,type:uuid,notnull"`
	ProfileUuid          string               `json:"profile_uuid" bun:"profile_uuid,type:uuid,pk,notnull"`
	PhoneNumber          string               `json:"phone_number" bun:"phone_number,type:text,notnull"`
	Email                string               `json:"email" bun:"email,type:text"`
	AvatarUrl            string               `json:"avatar_url" bun:"avatar_url,type:text"`
	Fullname             string               `json:"fullname" bun:"fullname,type:text,notnull"`
	Gender               string               `json:"gender" bun:"gender,type:text"`
	Birthday             string               `json:"birthday" bun:"birthday,type:text"`
	JobTitle             string               `json:"job_title" bun:"job_title,type:text"`
	Address              string               `json:"address" bun:"address,type:text"`
	IdentityNumber       string               `json:"identity_number" bun:"identity_number,type:text"`
	IdentityIssueOn      string               `json:"identity_issue_on" bun:"identity_issue_on,type:text"`
	IdentityIssueAt      string               `json:"identity_issue_at" bun:"identity_issue_at,type:text"`
	Passport             string               `json:"passport" bun:"passport,type:text"`
	Description          string               `json:"description" bun:"description,type:text"`
	RefId                string               `json:"ref_id" bun:"ref_id,type:text"`
	RefCode              string               `json:"ref_code" bun:"ref_code,type:text"`
	MoreInformation      string               `json:"more_information" bun:"more_information,type:text"`
	Country              string               `json:"country" bun:"country,type:text"`
	Province             string               `json:"province" bun:"province,type:text"`
	District             string               `json:"district" bun:"district,type:text"`
	Ward                 string               `json:"ward" bun:"ward,type:text"`
	ContactSocial        *omni.Social         `json:"contact_social" bun:"contact_social,type:text"`
	Status               bool                 `json:"status" bun:"status,type:boolean"`
	CallId               string               `json:"call_id" bun:"call_id,type:text"`
	StatusCall           string               `json:"status_call" bun:"status_call,type:text"`
	AltStatusCall        string               `json:"alt_status_call" bun:"alt_status_call,type:text"`
	AdditionalInfo       string               `json:"additional_info" bun:"additional_info,type:text"`
	AssignnerUuid        string               `json:"assignner_uuid" bun:"assignner_uuid,type:text"`
	AssigneeUuid         string               `json:"assignee_uuid" bun:"assignee_uuid,type:text"`
	ContactUuid          string               `json:"contact_uuid" bun:"contact_uuid,type:uuid,notnull"`
	RelatedProfileUuid   string               `json:"related_profile_uuid" bun:"related_profile_uuid,type:uuid,default:NULL"`
	ProfileType          string               `json:"profile_type" bun:"profile_type,type:text,notnull"`
	ProfileName          string               `json:"profile_name" bun:"profile_name,type:text"`
	ProfileCode          string               `json:"profile_code" bun:"profile_code,type:text"`
	SocialMappingContact SocialMappingContact `json:"social_mapping_contact" bun:"social_mapping_contact,type:jsonb"`
	CreatedBy            string               `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy            string               `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt            time.Time            `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt            time.Time            `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
}

type ProfileView struct {
	bun.BaseModel        `bun:"profile,alias:p"`
	DomainUuid           string               `json:"domain_uuid" bun:"domain_uuid"`
	ProfileUuid          string               `json:"profile_uuid" bun:"profile_uuid,pk"`
	PhoneNumber          string               `json:"phone_number" bun:"phone_number"`
	Email                string               `json:"email" bun:"email"`
	AvatarUrl            string               `json:"avatar_url" bun:"avatar_url"`
	Fullname             string               `json:"fullname" bun:"fullname"`
	Gender               string               `json:"gender" bun:"gender"`
	Birthday             string               `json:"birthday" bun:"birthday"`
	JobTitle             string               `json:"job_title" bun:"job_title"`
	Address              string               `json:"address" bun:"address"`
	IdentityNumber       string               `json:"identity_number" bun:"identity_number"`
	IdentityIssueOn      string               `json:"identity_issue_on" bun:"identity_issue_on"`
	IdentityIssueAt      string               `json:"identity_issue_at" bun:"identity_issue_at"`
	Passport             string               `json:"passport" bun:"passport"`
	Description          string               `json:"description" bun:"description"`
	RefId                string               `json:"ref_id" bun:"ref_id"`
	RefCode              string               `json:"ref_code" bun:"ref_code"`
	MoreInformationStr   string               `json:"-" bun:"more_information"`
	MoreInformation      any                  `json:"more_information" bun:"-"`
	Country              string               `json:"country" bun:"country"`
	Province             string               `json:"province" bun:"province"`
	District             string               `json:"district" bun:"district"`
	Ward                 string               `json:"ward" bun:"ward"`
	Phones               []*ProfilePhone      `json:"-" bun:"rel:has-many,join:profile_uuid=profile_uuid"`
	Emails               []*ProfileEmail      `json:"-" bun:"rel:has-many,join:profile_uuid=profile_uuid"`
	ContactSocial        *omni.Social         `json:"contact_social" bun:"contact_social"`
	CallId               string               `json:"call_id" bun:"call_id"`
	StatusCall           string               `json:"status_call" bun:"status_call"`
	AltStatusCall        string               `json:"alt_status_call" bun:"alt_status_call"`
	AdditionalInfo       string               `json:"additional_info" bun:"additional_info"`
	AssignnerUuid        string               `json:"assignner_uuid" bun:"assignner_uuid"`
	AssigneeUuid         string               `json:"assignee_uuid" bun:"assignee_uuid"`
	ContactUuid          string               `json:"contact_uuid" bun:"contact_uuid"`
	UserOwners           []*ProfileOwner      `json:"-" bun:"rel:has-many,join:profile_uuid=profile_uuid"`
	RelatedProfileUuid   string               `json:"related_profile_uuid" bun:"related_profile_uuid,type:uuid,default:NULL"`
	ListRelatedProfile   []*ProfileView       `json:"list_related_profile" bun:"rel:has-many,join:profile_uuid=related_profile_uuid"`
	Status               bool                 `json:"status"`
	ProfileType          string               `json:"profile_type" bun:"profile_type"`
	ProfileName          string               `json:"profile_name" bun:"profile_name"`
	ProfileCode          string               `json:"profile_code"`
	SocialMappingContact SocialMappingContact `json:"social_mapping_contact" bun:"social_mapping_contact"`
	CreatedBy            string               `json:"created_by" bun:"created_by"`
	UpdatedBy            string               `json:"updated_by" bun:"updated_by"`
	CreatedAt            time.Time            `json:"created_at" bun:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at" bun:"updated_at"`
}

type ProfilePost struct {
	AvartarBase64      string                       `json:"avartar_base64"`
	PhoneNumber        string                       `json:"phone_number"`
	Email              string                       `json:"email"`
	MoreInformation    []ContactMoreInformationData `json:"more_information"`
	RefId              string                       `json:"ref_id"`
	RefCode            string                       `json:"ref_code"`
	JobTitle           string                       `json:"job_title"`
	Description        string                       `json:"description"`
	Note               string                       `json:"note"`
	Birthday           string                       `json:"birthday"`
	Gender             string                       `json:"gender"`
	Fullname           string                       `json:"fullname"`
	IdentityNumber     string                       `json:"identity_number"`
	IdentityIssueOn    string                       `json:"identity_issue_on"`
	IdentityIssueAt    string                       `json:"identity_issue_at"`
	Passport           string                       `json:"passport"`
	Address            string                       `json:"address"`
	Country            string                       `json:"country"`
	Province           string                       `json:"province"`
	District           string                       `json:"district"`
	Ward               string                       `json:"ward"`
	Emails             []ContactMapData             `json:"emails"`
	Phones             []ContactMapData             `json:"phones"`
	UserOwners         []ContactMapData             `json:"user_owners"`
	Status             bool                         `json:"status"`
	ContactSocial      omni.Social                  `json:"contact_social"`
	CallId             string                       `json:"call_id"`
	StatusCall         string                       `json:"status_call"`
	AltStatusCall      string                       `json:"alt_status_call"`
	AdditionalInfo     string                       `json:"additional_info"`
	AssignnerUuid      string                       `json:"assignner_uuid"`
	AssigneeUuid       string                       `json:"assignee_uuid"`
	ProfileType        string                       `json:"profile_type"`
	ProfileName        string                       `json:"profile_name"`
	ProfileCode        string                       `json:"profile_code"`
	ContactUuid        string                       `json:"contact_uuid"`
	RelatedProfileUuid string                       `json:"related_profile_uuid"`
	LeadUuid           string                       `json:"lead_uuid"`
}

type ProfileOwner struct {
	bun.BaseModel `bun:"profile_owner,alias:po"`
	ProfileUuid   string    `json:"profile_uuid" bun:"profile_uuid,type:char(36),unique:unqx_profile_owner"`
	DomainUuid    string    `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull"`
	UserUuid      string    `json:"user_uuid" bun:"user_uuid,type:char(36),notnull,unique:unqx_profile_owner"`
	Username      string    `json:"username" bun:"username,type:text,notnull"`
	Type          string    `json:"type" bun:"type,type:text,notnull"`
	AssignedAt    time.Time `json:"assigned_at" bun:"assigned_at,nullzero"`
}

type ProfileChannel struct {
	bun.BaseModel `bun:"profile_channel,alias:pc"`
	ProfileUuid   string `json:"profile_uuid" bun:"profile_uuid,type:char(36)"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull"`
	Channel       string `json:"channel" bun:"channel,type:text,notnull"`
	Vendor        string `json:"vendor" bun:"vendor,type:text"`
	Data          string `json:"data" bun:"data,type:text"`
	PageId        string `json:"page_id" bun:"page_id,type:text"`
}

type ProfileNote struct {
	bun.BaseModel `bun:"profile_note,alias:pn"`
	ProfileUuid   string    `json:"profile_uuid" bun:"profile_uuid,type:char(36)"`
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

type ProfilePhone struct {
	bun.BaseModel `bun:"profile_phone,alias:pp"`
	ProfileUuid   string `json:"profile_uuid" bun:"profile_uuid,type:char(36),unique:unqx_profile_phone"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull"`
	Data          string `json:"data" bun:"data,type:text,unique:unqx_profile_phone"`
	Type          string `json:"type" bun:"type,type:text,notnull"`
}

type ProfileEmail struct {
	bun.BaseModel `bun:"profile_email,alias:pe"`
	ProfileUuid   string `json:"profile_uuid" bun:"profile_uuid,type:char(36),unique:unqx_profile_email"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull"`
	Data          string `json:"data" bun:"data,type:text,unique:unqx_profile_email"`
	EmailType     string `json:"email_type" bun:"email_type,type:text,notnull"`
	EmailAddress  string `json:"email_address" bun:"email_address,type:text"`
	EmailPrimary  int32  `json:"email_primary" bun:"email_primary,type:integer,notnull"`
}

type ProfileInfo struct {
	ProfileUuid          string                `json:"profile_uuid"`
	AvatarUrl            string                `json:"avatar_url"`
	MoreInformation      any                   `json:"more_information"`
	RefId                string                `json:"ref_id"`
	RefCode              string                `json:"ref_code"`
	JobTitle             string                `json:"job_title"`
	Description          string                `json:"description"`
	Birthday             string                `json:"birthday"`
	Gender               string                `json:"gender"`
	Fullname             string                `json:"fullname"`
	IdentityNumber       string                `json:"identity_number"`
	IdentityIssueOn      string                `json:"identity_issue_on"`
	IdentityIssueAt      string                `json:"identity_issue_at"`
	Passport             string                `json:"passport"`
	Address              string                `json:"address"`
	PhoneNumber          string                `json:"phone_number"`
	Email                string                `json:"email"`
	Country              string                `json:"country"`
	Province             string                `json:"province"`
	District             string                `json:"district"`
	Ward                 string                `json:"ward"`
	Emails               []ContactMapData      `json:"emails"`
	Phones               []ContactMapData      `json:"phones"`
	UserOwners           []UserOwnerShortInfo  `json:"user_owners"`
	Status               bool                  `json:"status"`
	ContactSocial        *omni.Social          `json:"contact_social"`
	CallId               string                `json:"call_id"`
	StatusCall           string                `json:"status_call"`
	AltStatusCall        string                `json:"alt_status_call"`
	AdditionalInfo       string                `json:"additional_info"`
	AssignnerUuid        string                `json:"assignner_uuid"`
	AssigneeUuid         string                `json:"assignee_uuid"`
	ProfileType          string                `json:"profile_type"`
	ProfileName          string                `json:"profile_name"`
	ProfileCode          string                `json:"profile_code"`
	RelatedProfileUuid   string                `json:"related_profile_uuid"`
	ListRelatedProfile   []*ProfileView        `json:"list_related_profile"`
	ContactUuid          string                `json:"contact_uuid"`
	SocialMappingContact *SocialMappingContact `json:"social_mapping_contact"`
}

type SocialMappingContact struct {
	Facebook string `json:"facebook"`
	Zalo     string `json:"zalo"`
}

type ProfileManageView struct {
	bun.BaseModel        `bun:"profile,alias:p"`
	DomainUuid           string               `json:"domain_uuid" bun:"domain_uuid"`
	ProfileUuid          string               `json:"profile_uuid" bun:"profile_uuid,pk"`
	PhoneNumber          string               `json:"phone_number" bun:"phone_number"`
	Email                string               `json:"email" bun:"email"`
	AvatarUrl            string               `json:"avatar_url" bun:"avatar_url"`
	Fullname             string               `json:"fullname" bun:"fullname"`
	Gender               string               `json:"gender" bun:"gender"`
	Birthday             string               `json:"birthday" bun:"birthday"`
	JobTitle             string               `json:"job_title" bun:"job_title"`
	Address              string               `json:"address" bun:"address"`
	IdentityNumber       string               `json:"identity_number" bun:"identity_number"`
	IdentityIssueOn      string               `json:"identity_issue_on" bun:"identity_issue_on"`
	IdentityIssueAt      string               `json:"identity_issue_at" bun:"identity_issue_at"`
	Passport             string               `json:"passport" bun:"passport"`
	Country              string               `json:"country" bun:"country"`
	Province             string               `json:"province" bun:"province"`
	District             string               `json:"district" bun:"district"`
	Ward                 string               `json:"ward" bun:"ward"`
	Phones               []*ProfilePhone      `json:"-" bun:"rel:has-many,join:profile_uuid=profile_uuid"`
	Emails               []*ProfileEmail      `json:"-" bun:"rel:has-many,join:profile_uuid=profile_uuid"`
	ContactSocial        *omni.Social         `json:"contact_social" bun:"contact_social"`
	ContactUuid          string               `json:"contact_uuid" bun:"contact_uuid"`
	UserOwners           []*ProfileOwner      `json:"-" bun:"rel:has-many,join:profile_uuid=profile_uuid"`
	RelatedProfileUuid   string               `json:"related_profile_uuid" bun:"related_profile_uuid,type:uuid,default:NULL"`
	ListRelatedProfile   []*ProfileView       `json:"list_related_profile" bun:"rel:has-many,join:related_profile_uuid=profile_uuid"`
	Status               bool                 `json:"status"`
	ProfileType          string               `json:"profile_type" bun:"profile_type"`
	ProfileName          string               `json:"profile_name" bun:"profile_name"`
	ProfileCode          string               `json:"profile_code"`
	SocialMappingContact SocialMappingContact `json:"social_mapping_contact" bun:"social_mapping_contact"`
	TotalOpen            int                  `json:"total_open" bun:"total_open"`
	TotalProcessing      int                  `json:"total_processing" bun:"total_processing"`
	TotalWaiting         int                  `json:"total_waiting" bun:"total_waiting"`
	TotalSolved          int                  `json:"total_solved" bun:"total_solved"`
	TotalPending         int                  `json:"total_pending" bun:"total_pending"`
	TotalReopen          int                  `json:"total_reopen" bun:"total_reopen"`
	CreatedBy            string               `json:"created_by" bun:"created_by"`
	UserCreatedBy        string               `json:"user_created_by" bun:"user_created_by"`
	UpdatedBy            string               `json:"updated_by" bun:"updated_by"`
	CreatedAt            time.Time            `json:"created_at" bun:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at" bun:"updated_at"`

	// Primary profile
	// pp.profile_uuid as primary_profile_uuid, pp.fullname as primary_profile_fullname, pp.profile_type as primary_profile_type, pp.profile_name as primary_profile_name
	PrimaryProfileUuid     string `json:"primary_profile_uuid" bun:"primary_profile_uuid"`
	PrimaryProfileFullname string `json:"primary_profile_fullname" bun:"primary_profile_fullname"`
	PrimaryProfileType     string `json:"primary_profile_type" bun:"primary_profile_type"`
	PrimaryProfileName     string `json:"primary_profile_name" bun:"primary_profile_name"`
}
