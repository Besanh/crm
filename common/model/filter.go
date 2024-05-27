package model

import (
	"database/sql"
	"time"
)

type ProfileFilter struct {
	Active         bool         `json:"active"`
	Fullname       string       `json:"fullname"`
	PhoneNumber    string       `json:"phone_number"`
	Email          string       `json:"email"`
	CustomerId     string       `json:"customer_id"`
	NationalId     string       `json:"national_id"`
	Address        string       `json:"address"`
	Gender         string       `json:"gender"`
	Country        string       `json:"country"`
	Province       string       `json:"province"`
	District       string       `json:"district"`
	Ward           string       `json:"ward"`
	Status         sql.NullBool `json:"status"`
	UserOwerUuid   string       `json:"user_owner_uuid"`
	RefId          string       `json:"ref_id"`
	RefCode        string       `json:"ref_code"`
	ProfileCode    string       `json:"profile_code"`
	Birthday       string       `json:"birthday"`
	JobTitle       string       `json:"job_title"`
	IdentityNumber string       `json:"identity_number"`
	Passport       string       `json:"password"`
	StartTime      time.Time    `json:"start_time"`
	EndTime        time.Time    `json:"end_time"`
	Common         string       `json:"common"` // fullname, profile_code, phone_number, email
	ProfileType    string       `json:"profile_type"`
	FacebookUserId string       `json:"facebook_user_id"`
	ZaloUserId     string       `json:"zalo_user_id"`
	CreatedBy      string       `json:"created_by"`
	IsAll          bool         `json:"is_all"`
	ProfileUuids   []string     `json:"profile_uuids"`
}

type ContactFilter struct {
	ContactType string
	ContactName string
	StartTime   time.Time
	EndTime     time.Time
}

type EventCalendarCategoryFilter struct {
	Title     string
	Color     string
	Status    sql.NullBool
	CreatedBy string
}

type EventCalendarFilter struct {
	EccUuid   string
	StartTime time.Time
	EndTime   time.Time
	CreatedBy string
}

type SolutionFilter struct {
	SolutionCode string       `json:"solution_code"`
	SolutionName string       `json:"solution_name"`
	Status       sql.NullBool `json:"status"`
	FileType     string       `json:"file_type"`
}

type UnitFilter struct {
	ParentUnitUuid string       `json:"parent_unit_uuid"`
	UnitUuid       string       `json:"unit_uuid"`
	UnitName       string       `json:"unit_name"`
	UnitCode       string       `json:"unit_code"`
	UnitLeader     string       `json:"unit_leader"`
	UnitBasis      sql.NullBool `json:"unit_basis"`
	Status         sql.NullBool `json:"status"`
	Level          string       `json:"level"`
	IsParent       bool         `json:"is_parent"`
	UnitUuids      []string     `json:"unit_uuids"`
	FromUnitLevel  string       `json:"from_unit_level"`
	ToUnitLevel    string       `json:"to_unit_level"`
	Formular       string       `json:"formular"`
	Encompass      []string     `json:"encompass"`
}

type RoleGroupFilter struct {
	RoleGroupName string       `json:"role_group_name"`
	Status        sql.NullBool `json:"status"`
	StartTime     sql.NullTime `json:"start_time"`
	EndTime       sql.NullTime `json:"end_time"`
	Limit         int          `json:"limit"`
	Offset        int          `json:"offset"`
}

type UserFilter struct {
	Fullname               string `json:"fullname"`
	Name                   string `json:"username"`
	Level                  string `json:"level"`
	Levels                 []string
	UserUuid               []string
	StartTime              string `json:"start_time"`
	EndTime                string `json:"end_time"`
	Email                  string `json:"email"`
	IsAll                  bool
	Order                  string `json:"order"`
	OrderDirection         string `json:"order_direction"`
	Enabled                string `json:"user_enabled"`
	IsMapExtension         string
	RoleUuid               string `json:"role_uuid"`
	UnitUuid               string `json:"unit_uuid"`
	ManageUserUuids        []string
	ManageUsers            []*UserView `json:"manage_users"`
	ManageExtensionUuids   []string
	Extension              string
	Common                 string
	ManageExcludeUnitUuids []string
}

type MonitorFilter struct {
	GroupUuids    []string
	CampaignUuids []string
	UserUuids     []string
}

type ExtensionFilter struct {
	Extension            string
	Username             string
	Fullname             string
	ManageExtensionUuids []string
	Enabled              sql.NullBool
	UnitUuid             string
	Common               string
}

type GroupFilter struct {
	GroupName        string
	ManageGroupUuids []string
}

type IsBool struct {
	Value  bool
	IsNull bool
}

type IsString struct {
	Value  string
	IsNull bool
}
type ParamAutoImportLead struct {
	HopperCampaignUuid string `json:"hopper_campaign_uuid"`
	DomainUuid         string `json:"domain_uuid"`
	CampaignUuid       string `json:"campaign_uuid"`
	ListUuid           string `json:"list_uuid"`
	LeadName           string `json:"lead_name"`
	PhoneNumber        string `json:"phone_number"`
	Status             string `json:"status"`
	StartTime          string `json:"start_time"`
	EndTime            string `json:"end_time"`
	Interval           int    `json:"interval"`
}

type LeadFilter struct {
	LeadName     string    `json:"lead_name"`
	PhoneNumber  string    `json:"phone_number"`
	Status       string    `json:"status"`
	AltStatus    string    `json:"alt_status"`
	CampaignUuid string    `json:"campaign_uuid"`
	ListUuid     string    `json:"list_uuid"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	OrderBy      string    `json:"order_by"`
	UserUuid     string    `json:"user_uuid"`
	AssigneeUuid string    `json:"assignee_uuid"`

	// Contract
	CustomerName   string   `json:"customer_name"`
	ProductName    string   `json:"product_name"`
	GroupCodeUuids []string `json:"group_code_uuid"`
}

type HopperFilter struct {
	CampaignUuids []string
	ListUuids     []string
}

type StatusFilter struct {
	CampaignUuid       string
	CategoryStatusUuid string
	StatusName         string
	StatusCode         string
}

type MailTemplateFilter struct {
}

type ReportFilter struct {
	ManageUserUuids      []string
	ManageExtensionUuids []string
	AllowedGroups        []string
	GroupUuids           []string
	CampaignUuids        []string
	DomainName           string
	Level                string
	StartTime            time.Time
	EndTime              time.Time
}

type AudioFilter struct {
	AudioName string
	UserUuid  string
}

type CampaignFilter struct {
	CampaignName        string
	CampaignUuids       []string
	Types               []string
	ManageCampaignUuids []string
	Active              IsBool
	Status              string
	UnitUuid            string
}

type CarrierFilter struct {
	DialType         string
	CarrierName      string
	ManageGroupUuids []string
}

type CallCenterQueueFilter struct {
	CallCenterQueueUuids []string
	CampaignUuids        []string
}

type GeneralFilter struct {
	AllowedGroups          []string
	GroupUuids             []string
	CampaignUuids          []string
	Extensions             []string
	StartTime              time.Time
	EndTime                time.Time
	DomainName             string
	Level                  string
	ManageUserUuids        []string
	ManageCampaignUuids    []string
	Directions             []string
	CampaignTypes          []string
	Limit                  int
	Offset                 int
	CarrierUuids           []string
	ManageExcludeUnitUuids []string
	UserStatus             string
	UnitUuid               string
	ManageNotInUserUuids   []string
	ManageUsers            []*UserView `json:"manage_users"`
	ManageExtensionUuids   []string
	TypeUnit               string `json:"type_unit"`
	ManageUnitUuids        []string
	ManageNotInUnitUuids   []string
}

type ContactTagFilter struct {
	TagName         string       `json:"tag_name"`
	TagType         string       `json:"tag_type"`
	LimitedFunction string       `json:"limited_function"`
	Status          sql.NullBool `json:"status"`
	StartTime       string       `json:"start_time"`
	EndTime         string       `json:"end_time"`
	FileType        string       `json:"file_type"`
}

type ContactGroupFilter struct {
	GroupType string       `json:"group_type"` // staff, member
	GroupName string       `json:"group_name"`
	Status    sql.NullBool `json:"status"`
	StartTime string       `json:"start_time"`
	EndTime   string       `json:"end_time"`
	FileType  string       `json:"file_type"`
}

type ContactCareerFilter struct {
	CareerType string       `json:"career_type"` // career
	CareerName string       `json:"career_name"`
	Status     sql.NullBool `json:"status"`
	StartTime  string       `json:"start_time"`
	EndTime    string       `json:"end_time"`
	FileType   string       `json:"file_type"`
}

type OmniFilter struct {
	OmniName  string       `json:"omni_name"`
	OmniType  string       `json:"omni_type"`
	Supplier  string       `json:"supplier"`
	Status    sql.NullBool `json:"status"`
	StartTime time.Time    `json:"start_time"`
	EndTime   time.Time    `json:"end_time"`
}

type WorkDayFilter struct {
	WorkDayName string       `json:"workday_name"`
	Status      sql.NullBool `json:"status"`
	StartTime   string       `json:"start_time"`
	EndTime     string       `json:"end_time"`
	FileType    string       `json:"file_type"`
	UnitUuids   []string     `json:"unit_uuids"`
	Limit       int          `json:"limit"`
	Offset      int          `json:"offset"`
}

type PbxFilter struct {
	PbxName   string       `json:"pbx_name"`
	Status    sql.NullBool `json:"status"`
	Verified  sql.NullBool `json:"verified"`
	StartTime string       `json:"start_time"`
	EndTime   string       `json:"end_time"`
	FileType  string       `json:"file_type"`
}

type ClassifyTagFilter struct {
	TagName         string       `json:"tag_name"`
	TagType         string       `json:"tag_type"`
	LimitedFunction string       `json:"limited_function"`
	Status          sql.NullBool `json:"status"`
	StartTime       string       `json:"start_time"`
	EndTime         string       `json:"end_time"`
	FileType        string       `json:"file_type"`
}

type ClassifyGroupFilter struct {
	GroupType string       `json:"group_type"` // staff, member
	GroupName string       `json:"group_name"`
	Status    sql.NullBool `json:"status"`
	StartTime string       `json:"start_time"`
	EndTime   string       `json:"end_time"`
	FileType  string       `json:"file_type"`
}

type ClassifyCareerFilter struct {
	CareerType string       `json:"career_type"` // career
	CareerName string       `json:"career_name"`
	Status     sql.NullBool `json:"status"`
	StartTime  string       `json:"start_time"`
	EndTime    string       `json:"end_time"`
	FileType   string       `json:"file_type"`
}

type ParentUnitFilter struct {
	Level string `json:"level"`
}

type CareerFilter struct {
	CareerCode      string       `json:"career_code"`
	CareerName      string       `json:"career_name"`
	Source          []string     `json:"source"`
	IsSearchExactly sql.NullBool `json:"is_search_exactly"`
}
