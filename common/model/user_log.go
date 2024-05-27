package model

import (
	"time"

	"github.com/uptrace/bun"
)

type UserLog struct {
	bun.BaseModel `bun:"user_log,alias:ul"`
	UserLogUuid   string    `json:"user_log_uuid" bun:"user_log_uuid,type:char(36),pk"`
	DomainUuid    string    `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull"`
	UserUuid      string    `json:"user_uuid" bun:"user_uuid,type:char(36),notnull"`
	ExtensionUuid string    `json:"extension_uuid" bun:"extension_uuid,type:char(36),notnull"`
	EventTime     time.Time `json:"event_time" bun:"event_time,type:timestamp,default:current_timestamp"`
	CampaignUuid  string    `json:"campaign_uuid" bun:"campaign_uuid,type:char(36),nullzero"`
	LeadUuid      string    `json:"lead_uuid" bun:"lead_uuid,type:char(36),nullzero"`
	PauseEpoch    int64     `json:"pause_epoch" bun:"pause_epoch,type:integer,nullzero"`
	PauseSec      int64     `json:"pause_sec" bun:"pause_sec,type:integer"`
	WaitEpoch     int64     `json:"wait_epoch" bun:"wait_epoch,type:integer,nullzero"`
	WaitSec       int64     `json:"wait_sec" bun:"wait_sec,type:integer"`
	TalkEpoch     int64     `json:"talk_epoch" bun:"talk_epoch,type:integer,nullzero"`
	TalkSec       int64     `json:"talk_sec" bun:"talk_sec,type:integer"`
	DispoEpoch    int64     `json:"dispo_epoch" bun:"dispo_epoch,type:integer,nullzero"`
	DispoSec      int64     `json:"dispo_sec" bun:"dispo_sec,type:integer"`
	Status        string    `json:"status" bun:"status,type:text,type:varchar(20),nullzero"`
	SubStatus     string    `json:"sub_status" bun:"sub_status,type:varchar(20),nullzero"`
	Comment       string    `json:"comment" bun:"comment,type:text,nullzero"`
	PauseCode     string    `json:"pause_code" bun:"pause_code,type:varchar(20),nullzero"`
	PauseType     string    `json:"pause_type" bun:"pause_type,type:varchar(20),nullzero,default:'UNDEFINED'"`
	CallUuid      string    `json:"call_uuid" bun:"call_uuid,type:char(36),nullzero"`
}
