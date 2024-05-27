package model

import "github.com/uptrace/bun"

type AgentInfo struct {
	bun.BaseModel `bun:"v_users,alias:u"`
	UserUuid      string                `json:"user_uuid" bun:"user_uuid,pk"`
	DomainUuid    string                `json:"domain_uuid" bun:"domain_uuid"`
	Username      string                `json:"username" bun:"username"`
	UserEnabled   string                `json:"user_enabled" bun:"user_enabled"`
	UserStatus    string                `json:"user_status" bun:"user_status"`
	Level         string                `json:"level" bun:"level"`
	LastName      string                `json:"last_name" bun:"last_name"`
	MiddleName    string                `json:"middle_name" bun:"middle_name"`
	FirstName     string                `json:"first_name" bun:"first_name"`
	Campaigns     []*CampaignOptionView `json:"campaigns" bun:"m2m:campaign_users,join:User=Campaign"`
	UserLive      *UserLive             `json:"user_live" bun:"rel:has-one,join:user_uuid=user_uuid"`
}
