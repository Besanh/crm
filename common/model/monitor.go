package model

import (
	"time"

	"github.com/uptrace/bun"
)

type SubscriberMonitorView struct {
	DomainUuid     string    `json:"domain_uuid"`
	UserUuid       string    `json:"user_uuid"`
	Username       string    `json:"username"`
	Extension      string    `json:"extension"`
	TotalCalls     int64     `json:"total_calls"`
	TotalAnswered  int64     `json:"total_answered"`
	TotalTalktime  int64     `json:"total_talktime"`
	LastName       string    `json:"last_name"`
	MiddleName     string    `json:"middle_name"`
	FirstName      string    `json:"first_name"`
	Status         string    `json:"status"`
	LoginTime      time.Time `json:"login_time"`
	LastUpdateTime time.Time `json:"last_update_time"`
}

type SipRegistrationMonitorView struct {
	SipUser    string    `json:"sip_user"`
	ServerHost string    `json:"server_host"`
	NetworkIp  string    `json:"network_ip"`
	UserAgent  string    `json:"user_agent"`
	CallTime   time.Time `json:"call_time"`
	CallState  string    `json:"call_state"`
	CallId     string    `json:"call_id"`
	Caller     string    `json:"caller"`
	Callee     string    `json:"callee"`
	Direction  string    `json:"direction"`
}

type SipDialogMonitorView struct {
	CallTime  time.Time `json:"call_time"`
	CallState string    `json:"call_state"`
	CallId    string    `json:"call_id"`
	Caller    string    `json:"caller"`
	Callee    string    `json:"callee"`
	Direction string    `json:"direction"`
}

type HopperCampaignMonitorView struct {
	CampaignUuid   string `json:"campaign_uuid"`
	CampaignName   string `json:"campaign_name"`
	LocalStartTime string `json:"local_start_time"`
	LocalEndTime   string `json:"local_end_time"`
	Type           string `json:"type"`
	TotalReady     int    `json:"total_ready"`
	TotalIncall    int    `json:"total_incall"`
	TotalDnc       int    `json:"total_dnc"`
	TotalHopper    int    `json:"total_hopper"`
	ConcurrentCall int    `json:"concurrent_call"`
	Ratio          int    `json:"ratio"`
	TotalLead      int    `json:"total_lead"`
	TotalCalled    int    `json:"total_called"`
}

type HopperLeadMonitorView struct {
	LeadUuid    string `json:"lead_uuid"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
	CallTime    string `json:"call_time"`
}

type UserLiveMonitor struct {
	bun.BaseModel `bun:"v_users,alias:u"`
	UserUuid      string `json:"user_uuid" bun:"user_uuid,pk"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid"`
	Username      string `json:"username" bun:"username"`
	LastName      string `json:"last_name" bun:"last_name"`
	MiddleName    string `json:"middle_name" bun:"middle_name"`
	FirstName     string `json:"first_name" bun:"first_name"`
	Extension     string `json:"extension" bun:"extension"`
	LiveStatus    string `json:"live_status" bun:"live_status"`
	PauseCode     string `json:"pause_code" bun:"pause_code"`
	CampaignName  string `json:"campaign_name" bun:"campaign_name"`
	LogStatus     string `json:"log_status" bun:"log_status"`
	CallUuid      string `json:"call_uuid" bun:"call_uuid"`
	LatestTime    string `json:"latest_time" bun:"latest_time"`
	EarliestTime  string `json:"earliest_time" bun:"earliest_time"`
	TotalPauseSec int64  `json:"total_pause_sec" bun:"total_pause_sec"`
	TotalWaitSec  int64  `json:"total_wait_sec" bun:"total_wait_sec"`
	TotalTalkSec  int64  `json:"total_talk_sec" bun:"total_talk_sec"`
	TotalDispoSec int64  `json:"total_dispo_sec" bun:"total_dispo_sec"`
	UnitUuid      string `json:"unit_uuid" bun:"unit_uuid"`
}

type MonitorCampaignData struct {
	bun.BaseModel  `bun:"campaign,alias:ca"`
	CampaignUuid   string    `json:"campaign_uuid" bun:"campaign_uuid"`
	CampaignName   string    `json:"campaign_name" bun:"campaign_name"`
	Type           string    `json:"type" bun:"type"`
	Status         string    `json:"status" bun:"status"`
	TotalLead      int       `json:"total_lead" bun:"total_lead"`
	TotalCompleted int       `json:"total_completed" bun:"total_completed"`
	TotalIncall    int       `json:"total_incall" bun:"total_incall"`
	TotalReady     int       `json:"total_ready" bun:"total_ready"`
	TotalHopper    int       `json:"total_hopper" bun:"total_hopper"`
	TotalRemain    int       `json:"total_remain" bun:"total_remain"`
	TotalRecall    int       `json:"total_recall" bun:"total_recall"`
	TotalAgent     int       `json:"total_agent" bun:"total_agent"`
	CreatedAt      time.Time `json:"created_at" bun:"created_at"`
	LocalStartTime string    `json:"local_start_time" bun:"local_start_time"`
	LocalEndTime   string    `json:"local_end_time" bun:"local_end_time"`
	Active         bool      `json:"active" bun:"active"`
}

type LiveCallMonitor struct {
	AgentName     string `json:"agent_name" bun:"agent_name"`
	BCIDNum       string `json:"b_cid_num" bun:"b_cid_num"`
	CallID        string `json:"call_id" bun:"call_id"`
	CallState     string `json:"call_state" bun:"call_state"`
	CallTime      string `json:"call_time" bun:"call_time"`
	CIDNum        string `json:"cid_num" bun:"cid_num"`
	Destination   string `json:"destination" bun:"destination"`
	Direction     string `json:"direction" bun:"direction"`
	DomainName    string `json:"domain_name" bun:"domain_name"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid"`
	Extension     string `json:"extension" bun:"extension"`
	ExtensionUuid string `json:"extension_uuid" bun:"extension_uuid"`
	LastName      string `json:"last_name" bun:"last_name"`
	MiddleName    string `json:"middle_name" bun:"middle_name"`
	FirstName     string `json:"first_name" bun:"first_name"`
	NetworkIP     string `json:"network_ip" bun:"network_ip"`
	ServerHost    string `json:"server_host" bun:"server_host"`
	UserAgent     string `json:"user_agent" bun:"user_agent"`
	UserUuid      string `json:"user_uuid" bun:"user_uuid,pk"`
	Username      string `json:"username" bun:"username"`
	UnitUuid      string `json:"unit_uuid" bun:"unit_uuid"`
	UnitName      string `json:"unit_name" bun:"unit_name"`
}

type CallCenterQueueMonitor struct {
	CallCenterQueueUuid string    `json:"call_center_queue_uuid"`
	CallCenterQueueName string    `json:"call_center_queue_name"`
	CallCenterAgentUuid string    `json:"call_center_agent_uuid"`
	CallCenterAgentName string    `json:"call_center_agent_name"`
	Mobile              string    `json:"mobile"`
	Status              string    `json:"status"`
	JoinedTime          time.Time `json:"joined_time"`
	AbandonedTime       time.Time `json:"abandoned_time"`
	BridgedTime         time.Time `json:"bridged_time"`
	RejoinedTime        time.Time `json:"rejoined_time"`
	CallUuid            string    `json:"call_uuid"`
}
