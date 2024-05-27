package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Campaign struct {
	bun.BaseModel          `bun:"campaign,alias:ca"`
	DomainUuid             string    `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull,unique:domain_campaign_idx_unq"`
	Id                     string    `json:"campaign_uuid" bun:"campaign_uuid,pk,type:char(36)"`
	CampaignName           string    `json:"campaign_name" bun:"campaign_name,type:text,notnull,unique:domain_campaign_idx_unq"`
	Type                   string    `json:"type" bun:"type,type:text,notnull"`
	Description            string    `json:"description" bun:"description,type:text"`
	Active                 bool      `json:"active" bun:"active,type:bool,default:'1',nullzero"`
	ConcurrentCall         int       `json:"concurrent_call" bun:"concurrent_call,type:integer,default:1"`
	Ratio                  int       `json:"ratio" bun:"ratio,type:integer,default:100"`
	CarrierUuid            string    `json:"carrier_uuid" bun:"carrier_uuid,type:char(36),nullzero"`
	CallCenterQueueUuid    string    `json:"call_center_queue_uuid" bun:"call_center_queue_uuid,type:char(36),nullzero"`
	TemplateUuid           string    `json:"template_uuid" bun:"template_uuid,type:char(36),nullzero"`
	RecallTimes            int       `json:"recall_times" bun:"recall_times,type:integer,default:0,nullzero"`
	RecallDuration         int       `json:"limit_recall_duration" bun:"limit_recall_duration,type:integer,default:0,nullzero"`
	ScheduleRecall         string    `json:"schedule_recall" bun:"schedule_recall,type:text,default:NULL"`
	ScheduleRecallDuration int       `json:"schedule_recall_duration" bun:"schedule_recall_duration,type:integer,nullzero,default:0"`
	Hopper                 int       `json:"hopper" bun:"hopper,type:integer,default:20"`
	AnswerCallbackUrl      string    `json:"answer_callback_url" bun:"answer_callback_url,type:text,nullzero"`
	LocalStartTime         string    `json:"local_start_time" bun:"local_start_time,type:time,default:'00:00:00',nullzero"`
	LocalEndTime           string    `json:"local_end_time" bun:"local_end_time,type:time,default:'23:59:00',nullzero"`
	CustomerOrder          string    `json:"customer_order" bun:"customer_order,type:text,default:'id',nullzero"`
	AllowManualDial        bool      `json:"allow_manual_dial" bun:"allow_manual_dial,type:bool,default:'1',nullzero"`
	AllowSearchLead        bool      `json:"allow_search_lead" bun:"allow_search_lead,type:bool,default:'1',nullzero"`
	EnableCallbackAlert    bool      `json:"enable_callback_alert" bun:"enable_callback_alert,type:bool,default:'1',nullzero"`
	DefaultListUuid        string    `json:"default_list_uuid" bun:"default_list_uuid,type:char(36),nullzero"`
	ScriptUuid             string    `json:"script_uuid" bun:"script_uuid,type:char(36),nullzero"`
	Status                 string    `json:"status" bun:"status,type:text,notnull"`
	CreatedAt              time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt              time.Time `json:"updated_at" bun:"updated_at,type:timestamp,notnull,nullzero,default:current_timestamp"`

	// Autocall
	Network                   string       `json:"network" bun:"network,type:text,nullzero"`
	ModeCall                  string       `json:"mode_call" bun:"mode_call,type:text,nullzero"`
	RunId                     string       `json:"run_id" bun:"run_id,type:text,nullzero"`
	TypeAutocall              string       `json:"type_autocall" bun:"type_autocall,type:text,nullzero"`
	PriorityRecall            string       `json:"priority_recall" bun:"priority_recall,type:text,nullzero"`
	RecallStatus              []string     `json:"recall_status" bun:"recall_status,array"` //not type:array
	CallTimeout               int          `json:"call_timeout" bun:"call_timeout,type:integer,nullzero,default:60"`
	EnableEncrypt             bool         `json:"enable_encrypt" bun:"enable_encrypt,type:bool,default:'0',nullzero"`
	CallbackUrl               string       `json:"callback_url" bun:"callback_url,type:text,nullzero"`
	OrigCampaignUuid          string       `json:"orig_campaign_uuid" bun:"orig_campaign_uuid,type:text,nullzero"`
	LeadExpiredAt             bun.NullTime `json:"lead_expired_at" bun:"lead_expired_at,type:timestamp,nullzero,default:NULL"`
	AllowSelectStatus         bool         `json:"allow_select_status" bun:"allow_select_status,type:bool,default:'0',nullzero"`
	IsContractCampaign        bool         `json:"is_contract_campaign" bun:"is_contract_campaign,type:bool,default:'0',nullzero"`
	HotLeadInnerCampaignUuids []string     `json:"hot_lead_inner_campaign_uuids" bun:"hot_lead_inner_campaign_uuids,array"`
	CreatedBy                 string       `json:"created_by" bun:"created_by,type:char(36),nullzero"`
}

type CampaignOptionView struct {
	bun.BaseModel `bun:"campaign,alias:ca"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid"`
	Id            string `json:"campaign_uuid" bun:"campaign_uuid,pk"`
	CampaignName  string `json:"campaign_name" bun:"campaign_name"`
	Type          string `json:"type" bun:"type"`
}

type CampaignUser struct {
	bun.BaseModel `bun:"campaign_users"`
	DomainUuid    string              `json:"domain_uuid" bun:"domain_uuid,type:char(36)"`
	CampaignUuid  string              `json:"campaign_uuid" bun:"campaign_uuid,pk,type:char(36)"`
	UserUuid      string              `json:"user_uuid" bun:"user_uuid,pk,type:char(36)"`
	Campaign      *CampaignOptionView `bun:"rel:belongs-to,join:campaign_uuid=campaign_uuid"`
	User          *UserCampaignView   `bun:"rel:belongs-to,join:user_uuid=user_uuid"`
}

type CampaignView struct {
	bun.BaseModel           `bun:"campaign,alias:ca"`
	DomainUuid              string          `json:"domain_uuid" bun:"domain_uuid"`
	Id                      string          `json:"campaign_uuid" bun:"campaign_uuid,pk"`
	CampaignName            string          `json:"campaign_name" bun:"campaign_name"`
	Type                    string          `json:"type" bun:"type"`
	Description             string          `json:"description" bun:"description"`
	Active                  bool            `json:"active" bun:"active"`
	ConcurrentCall          int             `json:"concurrent_call" bun:"concurrent_call"`
	Ratio                   int             `json:"ratio" bun:"ratio"`
	CarrierUuid             string          `json:"carrier_uuid" bun:"carrier_uuid"`
	CallCenterQueueUuid     string          `json:"call_center_queue_uuid" bun:"call_center_queue_uuid"`
	TemplateUuid            string          `json:"template_uuid" bun:"template_uuid"`
	RecallTimes             int             `json:"recall_times" bun:"recall_times"`
	RecallDuration          int             `json:"limit_recall_duration" bun:"limit_recall_duration"`
	ScheduleRecall          string          `json:"schedule_recall" bun:"schedule_recall"`
	ScheduleRecallDuration  int             `json:"schedule_recall_duration" bun:"schedule_recall_duration"`
	Hopper                  int             `json:"hopper" bun:"hopper"`
	AnswerCallbackUrl       string          `json:"answer_callback_url" bun:"answer_callback_url"`
	LocalStartTime          string          `json:"local_start_time" bun:"local_start_time"`
	LocalEndTime            string          `json:"local_end_time" bun:"local_end_time"`
	CustomerOrder           string          `json:"customer_order" bun:"customer_order"`
	AllowManualDial         bool            `json:"allow_manual_dial" bun:"allow_manual_dial"`
	AllowSearchLead         bool            `json:"allow_search_lead" bun:"allow_search_lead"`
	EnableCallbackAlert     bool            `json:"enable_callback_alert" bun:"enable_callback_alert"`
	DefaultListUuid         string          `json:"default_list_uuid" bun:"default_list_uuid"`
	ScriptUuid              string          `json:"script_uuid" bun:"script_uuid"`
	Status                  string          `json:"status" bun:"status"`
	CreatedAt               time.Time       `json:"created_at" bun:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at" bun:"updated_at"`
	CallCenterQueueStrategy string          `json:"call_center_queue_strategy" bun:"call_center_queue_strategy"`
	TemplateName            string          `json:"template_name" bun:"template_name"`
	CarrierName             string          `json:"carrier_name" bun:"carrier_name"`
	ModeCall                string          `json:"mode_call" bun:"mode_call"`
	Network                 CampaignNetwork `json:"network" bun:"-"`
	NetworkValue            string          `json:"-" bun:"network"`
	RunId                   string          `json:"run_id" bun:"run_id"`
	TypeAutocall            string          `json:"type_autocall" bun:"type_autocall"`
	CallTimeout             int             `json:"call_timeout" bun:"call_timeout"`
	// Custom infor
	DomainName string                `json:"domain_name" bun:"domain_name"`
	Users      []*UserCampaignView   `json:"users" bun:"m2m:campaign_users,join:Campaign=User"`
	Groups     []*GroupModel         `json:"groups" bun:"m2m:campaign_groups,join:Campaign=Group"`
	Statuses   []*StatusCampaignView `json:"statuses" bun:"rel:has-many,join:campaign_uuid=campaign_uuid"`
	Schedules  []*CampaignSchedule   `json:"schedules" bun:"rel:has-many,join:campaign_uuid=campaign_uuid"`

	PriorityRecall     string       `json:"priority_recall" bun:"priority_recall"`
	RecallStatus       []string     `json:"recall_status" bun:"recall_status,array"` //not type:array
	EnableEncrypt      bool         `json:"enable_encrypt" bun:"enable_encrypt"`
	CallbackUrl        string       `json:"callback_url" bun:"callback_url"`
	OrigCampaignUuid   string       `json:"orig_campaign_uuid" bun:"orig_campaign_uuid"`
	LeadExpiredAt      bun.NullTime `json:"lead_expired_at" bun:"lead_expired_at"`
	AllowSelectStatus  bool         `json:"allow_select_status" bun:"allow_select_status"`
	IsContractCampaign bool         `json:"is_contract_campaign" bun:"is_contract_campaign"`

	HotLeadInnerCampaignUuids []string `json:"hot_lead_inner_campaign_uuids" bun:"hot_lead_inner_campaign_uuids,array"`
	CreatedBy                 string   `json:"created_by" bun:"created_by"`
	CreatedByUsername         string   `json:"created_by_username" bun:"created_by_username"`

	CustomData map[string]string `json:"custom_data" bun:"-"`
}

type CampaignPost struct {
	CampaignName   string `json:"campaign_name"`
	Type           string `json:"type"`
	Description    string `json:"description"`
	ConcurrentCall int    `json:"concurrent_call"`
	Ratio          int    `json:"ratio"`
	TemplateUuid   string `json:"template_uuid"`
	TypeAutocall   string `json:"type_autocall"`
}

type CampaignAssign struct {
	Users  []string `json:"users"`
	Groups []string `json:"groups"`
}

type GroupModel struct {
	bun.BaseModel `bun:"v_groups"`
	GroupUuid     string `json:"group_uuid" bun:"group_uuid,pk"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid"`
	GroupName     string `json:"group_name" bun:"group_name"`
}

type CampaignStatusView struct {
	bun.BaseModel `bun:"campaign,alias:ca"`
	DomainUuid    string        `json:"domain_uuid" bun:"domain_uuid"`
	Id            string        `json:"campaign_uuid" bun:"campaign_uuid,pk"`
	CampaignName  string        `json:"campaign_name" bun:"campaign_name"`
	Type          string        `json:"type" bun:"type"`
	Description   string        `json:"description" bun:"description"`
	Active        bool          `json:"active" bun:"active"`
	Statuses      []*StatusView `json:"statuses" bun:"rel:has-many,join:campaign_uuid=campaign_uuid"`
}

type CampaignPut struct {
	CampaignName              string       `json:"campaign_name"`
	Type                      string       `json:"type"`
	Description               string       `json:"description"`
	Active                    bool         `json:"active"`
	ConcurrentCall            int          `json:"concurrent_call"`
	Ratio                     int          `json:"ratio"`
	CarrierUuid               string       `json:"carrier_uuid"`
	CallCenterQueueUuid       string       `json:"call_center_queue_uuid"`
	TemplateUuid              string       `json:"template_uuid"`
	RecallTimes               int          `json:"recall_times"`
	RecallDuration            int          `json:"limit_recall_duration"`
	ScheduleRecall            string       `json:"schedule_recall"`
	ScheduleRecallDuration    int          `json:"schedule_recall_duration"`
	Hopper                    int          `json:"hopper"`
	AnswerCallbackUrl         string       `json:"answer_callback_url"`
	LocalStartTime            string       `json:"local_start_time"`
	LocalEndTime              string       `json:"local_end_time"`
	CustomerOrder             string       `json:"customer_order"`
	AllowManualDial           bool         `json:"allow_manual_dial"`
	AllowSearchLead           bool         `json:"allow_search_lead"`
	EnableCallbackAlert       bool         `json:"enable_callback_alert"`
	CallCenterQueueStrategy   string       `json:"call_center_queue_strategy"`
	DefaultListUuid           string       `json:"default_list_uuid"`
	ScriptUuid                string       `json:"script_uuid"`
	Network                   string       `json:"-"`
	ModeCall                  string       `json:"mode_call"`
	TypeAutocall              string       `json:"type_autocall"`
	Users                     []string     `json:"users"`
	Groups                    []string     `json:"groups"`
	Statuses                  []string     `json:"statuses"`
	CallTimeout               int          `json:"call_timeout"`
	PriorityRecall            string       `json:"priority_recall"`
	RecallStatus              []string     `json:"recall_status"`
	EnableEncrypt             bool         `json:"enable_encrypt"`
	CallbackUrl               string       `json:"callback_url"`
	OrigCampaignUuid          string       `json:"orig_campaign_uuid"`
	LeadExpiredAt             bun.NullTime `json:"lead_expired_at"`
	Status                    string       `json:"status"`
	AllowSelectStatus         bool         `json:"allow_select_status"`
	IsContractCampaign        bool         `json:"is_contract_campaign"`
	HotLeadInnerCampaignUuids []string     `json:"hot_lead_inner_campaign_uuids"`

	CustomData map[string]any `json:"custom_data"`
}

type CampaignAssignExtension struct {
	CampaignUuid string   `json:"campaign_uuid"`
	Extensions   []string `json:"extensions"`
}

type CampaignNetwork struct {
	Viettel int `json:"viettel"`
	Mobi    int `json:"mobi"`
	Vina    int `json:"vina"`
	Tel     int `json:"tel"`
	OffNet  int `json:"offnet"`
}

type CampaignPostCopy struct {
	CampaignName string `json:"campaign_name"`
	CampaignUuid string `json:"campaign_uuid"`
}

type CampaignCustomData struct {
	bun.BaseModel          `bun:"campaign_custom_data"`
	CampaignCustomDataUuid string `json:"campaign_custom_data_uuid" bun:"campaign_custom_data_uuid,pk,type:char(36)"`
	CampaignUuid           string `json:"campaign_uuid" bun:"campaign_uuid,type:char(36),notnull,unique:unq_campaign_custom_data_key"`
	Key                    string `json:"key" bun:"key,type:text,notnull,unique:unq_campaign_custom_data_key"`
	Value                  string `json:"value" bun:"value,type:text,nullzero"`
}

type CampaignGroup struct {
	bun.BaseModel `bun:"campaign_groups"`
	DomainUuid    string        `json:"domain_uuid" bun:"domain_uuid,type:char(36)"`
	CampaignUuid  string        `json:"campaign_uuid" bun:"campaign_uuid,pk,type:char(36)"`
	GroupUuid     string        `json:"group_uuid" bun:"group_uuid,pk,type:char(36)"`
	Campaign      *CampaignView `bun:"rel:belongs-to,join:campaign_uuid=campaign_uuid"`
	Group         *GroupModel   `bun:"rel:belongs-to,join:group_uuid=group_uuid"`
}
