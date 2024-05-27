package model

import "github.com/uptrace/bun"

type RingGroup struct {
	bun.BaseModel               `bun:"v_ring_groups"`
	DomainUuid                  string `json:"domain_uuid" bun:"domain_uuid"`
	RingGroupUuid               string `json:"ring_group_uuid" bun:"ring_group_uuid,pk"`
	RingGroupName               string `json:"ring_group_name" bun:"ring_group_name"`
	RingGroupExtension          string `json:"ring_group_extension" bun:"ring_group_extension"`
	RingGroupGreeting           string `json:"ring_group_greeting" bun:"ring_group_greeting,nullzero"`
	RingGroupContext            string `json:"ring_group_context" bun:"ring_group_context"`
	RingGroupCallTimeout        int64  `json:"ring_group_call_timeout" bun:"ring_group_call_timeout,nullzero"`
	RingGroupForwardDestination string `json:"ring_group_forward_destination" bun:"ring_group_forward_destination,nullzero"`
	RingGroupForwardEnabled     string `json:"ring_group_forward_enabled" bun:"ring_group_forward_enabled"`
	RingGroupCallerIdName       string `json:"ring_group_caller_id_name" bun:"ring_group_caller_id_name,nullzero"`
	RingGroupCallerIdNumber     string `json:"ring_group_caller_id_number" bun:"ring_group_caller_id_number,nullzero"`
	RingGroupCidNamePrefix      string `json:"ring_group_cid_name_prefix" bun:"ring_group_cid_name_prefix,nullzero"`
	RingGroupCidNumberPrefix    string `json:"ring_group_cid_number_prefix" bun:"ring_group_cid_number_prefix,nullzero"`
	RingGroupStrategy           string `json:"ring_group_strategy" bun:"ring_group_strategy"`
	RingGroupTimeoutApp         string `json:"ring_group_timeout_app" bun:"ring_group_timeout_app,nullzero"`
	RingGroupTimeoutData        string `json:"ring_group_timeout_data" bun:"ring_group_timeout_data,nullzero"`
	RingGroupDistinctiveRing    string `json:"ring_group_distinctive_ring" bun:"ring_group_distinctive_ring,nullzero"`
	RingGroupRingback           string `json:"ring_group_ringback" bun:"ring_group_ringback,nullzero"`
	RingGroupMissedCallApp      string `json:"ring_group_missed_call_app" bun:"ring_group_missed_call_app,nullzero"`
	RingGroupMissedCallData     string `json:"ring_group_missed_call_data" bun:"ring_group_missed_call_data,nullzero"`
	RingGroupEnabled            string `json:"ring_group_enabled" bun:"ring_group_enabled"`
	RingGroupDescription        string `json:"ring_group_description" bun:"ring_group_description,nullzero"`
	DialplanUuid                string `json:"dialplan_uuid" bun:"dialplan_uuid"`
	RingGroupForwardTollAllow   string `json:"ring_group_forward_toll_allow" bun:"ring_group_forward_toll_allow,nullzero"`
}

type RingGroupDestination struct {
	bun.BaseModel            `bun:"v_ring_group_destinations"`
	RingGroupDestinationUuid string `json:"ring_group_destination_uuid" bun:"ring_group_destination_uuid,pk"`
	DomainUuid               string `json:"domain_uuid" bun:"domain_uuid"`
	RingGroupUuid            string `json:"ring_group_uuid" bun:"ring_group_uuid"`
	DestinationNumber        string `json:"destination_number" bun:"destination_number"`
	DestinationDelay         int64  `json:"destination_delay" bun:"destination_delay"`
	DestinationTimeout       int64  `json:"destination_timeout" bun:"destination_timeout"`
	DestinationPrompt        int64  `json:"destination_prompt" bun:"destination_prompt,nullzero"`
}
