package model

import "github.com/uptrace/bun"

type FollowMe struct {
	bun.BaseModel        `bun:"v_follow_me"`
	DomainUuid           string `json:"domain_uuid" bun:"domain_uuid"`
	FollowMeUuid         string `json:"follow_me_uuid" bun:"follow_me_uuid,pk"`
	CidNamePrefix        string `json:"cid_name_prefix" bun:"cid_name_prefix,nullzero"`
	CidNumberPrefix      string `json:"cid_number_prefix" bun:"cid_number_prefix,nullzero"`
	DialString           string `json:"dial_string" bun:"dial_string"`
	FollowMeEnabled      bool   `json:"follow_me_enabled" bun:"follow_me_enabled"`
	FollowMeCallerIdUuid string `json:"follow_me_caller_id_uuid" bun:"follow_me_caller_id_uuid,nullzero"`
	FollowMeIgnoreBusy   bool   `json:"follow_me_ignore_busy" bun:"follow_me_ignore_busy"`
}

type FollowMeDestination struct {
	bun.BaseModel           `bun:"v_follow_me_destinations"`
	DomainUuid              string `json:"domain_uuid" bun:"domain_uuid"`
	FollowMeUuid            string `json:"follow_me_uuid" bun:"follow_me_uuid"`
	FollowMeDestinationUuid string `json:"follow_me_destination_uuid" bun:"follow_me_destination_uuid,pk"`
	FollowMeDestination     string `json:"follow_me_destination" bun:"follow_me_destination"`
	FollowMeDelay           int    `json:"follow_me_delay" bun:"follow_me_delay"`
	FollowMeTimeout         int    `json:"follow_me_timeout" bun:"follow_me_timeout"`
	FollowMePrompt          int    `json:"follow_me_prompt" bun:"follow_me_prompt,nullzero"`
	FollowMeOrder           int    `json:"follow_me_order" bun:"follow_me_order"`
}
