package model

import (
	"time"

	"github.com/uptrace/bun"
)

type CampaignSchedule struct {
	bun.BaseModel      `bun:"campaign_schedule,alias:sp"`
	DomainUuid         string        `json:"domain_uuid" bun:"domain_uuid,type:uuid,notnull"`
	Id                 string        `json:"id" bun:"id,type:uuid,notnull,default:uuid_generate_v4()"`
	CampaignUuid       string        `json:"campaign_uuid" bun:"campaign_uuid,type:uuid,notnull"`
	MondayStartTime    time.Duration `json:"monday_start_time" bun:"monday_start_time,type:bigint"`
	MondayEndTime      time.Duration `json:"monday_end_time" bun:"monday_end_time,type:bigint"`
	TuesdayStartTime   time.Duration `json:"tuesday_start_time" bun:"tuesday_start_time,type:bigint"`
	TuesdayEndTime     time.Duration `json:"tuesday_end_time" bun:"tuesday_end_time,type:bigint"`
	WednesdayStartTime time.Duration `json:"wednesday_start_time" bun:"wednesday_start_time,type:bigint"`
	WednesdayEndTime   time.Duration `json:"wednesday_end_time" bun:"wednesday_end_time,type:bigint"`
	ThursdayStartTime  time.Duration `json:"thursday_start_time" bun:"thursday_start_time,type:bigint"`
	ThursdayEndTime    time.Duration `json:"thursday_end_time" bun:"thursday_end_time,type:bigint"`
	FridayStartTime    time.Duration `json:"friday_start_time" bun:"friday_start_time,type:bigint"`
	FridayEndTime      time.Duration `json:"friday_end_time" bun:"friday_end_time,type:bigint"`
	SaturdayStartTime  time.Duration `json:"saturday_start_time" bun:"saturday_start_time,type:bigint"`
	SaturdayEndTime    time.Duration `json:"saturday_end_time" bun:"saturday_end_time,type:bigint"`
	SundayStartTime    time.Duration `json:"sunday_start_time" bun:"sunday_start_time,type:bigint"`
	SundayEndTime      time.Duration `json:"sunday_end_time" bun:"sunday_end_time,type:bigint"`
	CreatedAt          time.Time     `json:"-" bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt          time.Time     `json:"-" bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}