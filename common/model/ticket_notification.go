package model

import (
	"time"

	"github.com/uptrace/bun"
)

type TicketNotification struct {
	bun.BaseModel    `bun:"ticket_notification,alias:ticket_notification"`
	DomainUuid       string    `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	NotificationUuid string    `json:"notification_uuid" bun:"notification_uuid,pk,type: char(36),default:uuid_generate_v4()"`
	IsRead           bool      `json:"is_read" bun:"is_read,type:boolean,default:'false',notnull"`
	Subject          string    `json:"subject" bun:"subject,type:text"`
	Content          string    `json:"content" bun:"content,type:text"`
	NotificationType string    `json:"notification_type" bun:"notification_type,type:text"`
	ReceiverId       string    `json:"receiver_id" bun:"receiver_id,type: char(36)"`
	GroupUuid        string    `json:"group_uuid" bun:"group_uuid,type: char(36)"`
	Priority         string    `json:"priority" bun:"priority,type:text"`
	Active           string    `json:"active" bun:"active,type:text"`
	Channel          string    `json:"channel" bun:"channel,type:text"`
	StartTime        time.Time `json:"start_time" bun:"start_time,nullzero,notnull,default:current_timestamp"`
	// EndTime          time.Time `json:"end_time" bun:"end_time,type:timestamptz,default:now(),notnull"`
	IsDeleted   bool      `json:"-" bun:"is_deleted,type:boolean,default:'false',notnull"`
	DeletedAt   time.Time `json:"-" bun:",soft_delete,nullzero"`
	CreatedAt   time.Time `json:"-" bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
	CreatedBy   string    `json:"created_by" bun:"created_by,type: char(36),nullzero"`
	UpdatedBy   string    `json:"updated_by" bun:"updated_by,type: char(36),nullzero"`
	User        *User     `json:"user" bun:"rel:belongs-to,join:receiver_id=user_uuid"`
	Group       *Group    `json:"group" bun:"rel:has-one,join:group_uuid=group_uuid"`
	TicketUuid  string    `json:"ticket_uuid" bun:"-"`
	Id          string    `json:"id" bun:"-"`
	IsPerTime   bool      `json:"is_per_time" bun:"-"`
	IsNotify    bool      `json:"is_notify" bun:"-"`
	IsSystem    bool      `json:"is_system" bun:"-"`
	IsPush      bool      `json:"is_push" bun:"-"`
	Status      string    `json:"status"`
	ReleaseDate time.Time `json:"release_date" bun:"release_date,nullzero,notnull,default:current_timestamp"`
	TicketCode  string    `json:"ticket_code" bun:"ticket_code,type:text"`
}
