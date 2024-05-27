package model

import (
	"time"

	"github.com/uptrace/bun"
)

type EventCalendar struct {
	bun.BaseModel   `bun:"event_calendar,alias:ec"`
	DomainUuid      string                     `json:"domain_uuid" bun:"domain_uuid,type:char(36),notnull"`
	EcUuid          string                     `json:"ec_uuid" bun:"ec_uuid,type:char(36),pk,notnull"`
	EccUuid         string                     `json:"ecc_uuid" bun:"ecc_uuid,type:uuid,notnull"`
	Category        *EventCalendarCategory     `json:"category" bun:"rel:has-one,join:ecc_uuid=ecc_uuid"`
	Title           string                     `json:"title" bun:"title,type:text,notnull"`
	Description     string                     `json:"description" bun:"description,type:text"`
	Todo            []*EventCalendarTodo       `json:"todo" bun:"rel:has-many,join:ec_uuid=ec_uuid"`
	Attachment      []*EventCalendarAttachment `json:"attachment" bun:"rel:has-many,join:ec_uuid=ec_uuid"`
	RemindTypeEvent string                     `json:"remind_type_event" bun:"remind_type_event,type:text,notnull"` // type remind cua event: after,specific
	// if equal after
	RemindTime int    `json:"remind_time" bun:"remind_time,type:integer,notnull,default:0"`
	RemindType string `json:"remind_type" bun:"remind_type,type:text,notnull"` // type remind: minute, hour, day

	// if equal specific time
	IsWholeDay bool      `json:"is_whole_day" bun:"is_whole_day,type:boolean,notnull"`
	StartTime  time.Time `json:"start_time" bun:"start_time,type:timestamp,nullzero"`
	EndTime    time.Time `json:"end_time" bun:"end_time,type:timestamp,nullzero"`

	// Notify via failover
	IsNotifyWeb   bool `json:"is_notify_web" bun:"is_notify_web,type:boolean,notnull"`
	IsNotifyEmail bool `json:"is_notify_email" bun:"is_notify_email,type:boolean,notnull"`
	IsNotifySms   bool `json:"is_notify_sms" bun:"is_notify_sms,type:boolean,notnull"`
	IsNotifyZns   bool `json:"is_notify_zns" bun:"is_notify_zns,type:boolean,notnull"`
	IsNotifyCall  bool `json:"is_notify_call" bun:"is_notify_call,type:boolean,notnull"`

	Repeat      int `json:"repeat" bun:"repeat,type:integer,default:0"`
	CountRepeat int `json:"count_repeat" bun:"count_repeat,type:integer,default:0"`

	// doi tuong dinh kem

	Status    bool      `json:"status" bun:"status,type:boolean,notnull"`
	CreatedBy string    `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy string    `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
}

type EventCalendarForm struct {
	DomainUuid      string                     `json:"domain_uuid" form:"domain_uuid"`
	EcUuid          string                     `json:"ec_uuid" form:"ec_uuid"`
	EccUuid         string                     `json:"ecc_uuid" form:"ecc_uuid"`
	Title           string                     `json:"title" form:"title"`
	Description     string                     `json:"description" form:"description"`
	Todo            []*EventCalendarTodo       `json:"todo" form:"todo"`
	Attachment      []*EventCalendarAttachment `json:"attachment" form:"attachment"`
	RemindTypeEvent string                     `json:"remind_type_event" form:"remind_type_event"` // type remind cua event: after,specific
	// if equal after
	RemindTime int    `json:"remind_time" form:"remind_time"`
	RemindType string `json:"remind_type" form:"remind_type"` // type remind: minute, hour, day

	// if equal specific time
	IsWholeDay bool      `json:"is_whole_day" form:"is_whole_day"`
	StartTime  time.Time `json:"start_time" form:"start_time"`
	EndTime    time.Time `json:"end_time" form:"end_time"`

	// Notify via failover
	IsNotifyWeb   bool `json:"is_notify_web" form:"is_notify_web"`
	IsNotifyEmail bool `json:"is_notify" form:"is_notify"`
	IsNotifySms   bool `json:"is_notify_sms" form:"is_notify_sms"`
	IsNotifyZns   bool `json:"is_notify_zns" form:"is_notify_zns"`
	IsNotifyCall  bool `json:"is_notify_call" form:"is_notify_call"`

	Repeat      int    `json:"repeat" form:"repeat"`
	CountRepeat string `json:"count_repeat" form:"count_repeat"`

	// doi tuong dinh kem

	Status    bool      `json:"status" form:"status"`
	CreatedBy string    `json:"created_by" form:"created_by"`
	UpdatedBy string    `json:"updated_by" form:"updated_by"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
}
