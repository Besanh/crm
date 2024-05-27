package model

import "github.com/uptrace/bun"

type CallCenterQueue struct {
	bun.BaseModel                          `bun:"v_call_center_queues"`
	CallCenterQueueUuid                    string `json:"call_center_queue_uuid" bun:"call_center_queue_uuid"`
	DomainUuid                             string `json:"domain_uuid" bun:"domain_uuid"`
	DialplanUuid                           string `json:"dialplan_uuid" bun:"dialplan_uuid"`
	QueueName                              string `json:"queue_name" bun:"queue_name"`
	QueueExtension                         string `json:"queue_extension" bun:"queue_extension"`
	QueueStrategy                          string `json:"queue_strategy" bun:"queue_strategy"`
	QueueMohSound                          string `json:"queue_moh_sound" bun:"queue_moh_sound"`
	QueueRecordTemplate                    string `json:"queue_record_template" bun:"queue_record_template"`
	QueueTimeBaseScore                     string `json:"queue_time_base_score" bun:"queue_time_base_score"`
	QueueMaxWaitTime                       int    `json:"queue_max_wait_time" bun:"queue_max_wait_time"`
	QueueMaxWaitTimeWithNoAgent            int    `json:"queue_max_wait_time_with_no_agent" bun:"queue_max_wait_time_with_no_agent"`
	QueueMaxWaitTimeWithNoAgentTimeReached int    `json:"queue_max_wait_time_with_no_agent_time_reached" bun:"queue_max_wait_time_with_no_agent_time_reached"`
	QueueTierRulesApply                    string `json:"queue_tier_rules_apply" bun:"queue_tier_rules_apply"`
	QueueTierRuleWaitSecond                int    `json:"queue_tier_rule_wait_second" bun:"queue_tier_rule_wait_second"`
	QueueTierRuleNoAgentNoWait             string `json:"queue_tier_rule_no_agent_no_wait" bun:"queue_tier_rule_no_agent_no_wait"`
	QueueTimeoutAction                     string `json:"queue_timeout_action" bun:"queue_timeout_action,nullzero"`
	QueueDiscardAbandonedAfter             int    `json:"queue_discard_abandoned_after" bun:"queue_discard_abandoned_after"`
	QueueAbandonedResumeAllowed            string `json:"queue_abandoned_resume_allowed" bun:"queue_abandoned_resume_allowed"`
	QueueTierRuleWaitMultiplyLevel         string `json:"queue_tier_rule_wait_multiply_level" bun:"queue_tier_rule_wait_multiply_level"`
	QueueCidPrefix                         string `json:"queue_cid_prefix" bun:"queue_cid_prefix,nullzero"`
	QueueAnnounceSound                     string `json:"queue_announce_sound" bun:"queue_announce_sound,nullzero"`
	QueueAnnounceFrequency                 int    `json:"queue_announce_frequency" bun:"queue_announce_frequency,nullzero"`
	QueueCcExitKeys                        string `json:"queue_cc_exit_keys" bun:"queue_cc_exit_keys,nullzero"`
	QueueDescription                       string `json:"queue_description" bun:"queue_description"`
}

type CallCenterAgent struct {
	bun.BaseModel          `bun:"v_call_center_agents,alias:cca"`
	CallCenterAgentUuid    string `json:"call_center_agent_uuid" bun:"call_center_agent_uuid"`
	DomainUuid             string `json:"domain_uuid" bun:"domain_uuid,nullzero"`
	UserUuid               string `json:"user_uuid" bun:"user_uuid,nullzero"`
	AgentName              string `json:"agent_name" bun:"agent_name,nullzero"`
	AgentType              string `json:"agent_type" bun:"agent_type,nullzero"`
	AgentCallTimeout       int    `json:"agent_call_timeout" bun:"agent_call_timeout,nullzero"`
	AgentId                string `json:"agent_id" bun:"agent_id,nullzero"`
	AgentPassword          string `json:"agent_password" bun:"agent_password,nullzero"`
	AgentContact           string `json:"agent_contact" bun:"agent_contact,nullzero"`
	AgentStatus            string `json:"agent_status" bun:"agent_status,nullzero"`
	AgentLogout            string `json:"agent_logout" bun:"agent_logout,nullzero"`
	AgentMaxNoAnswer       int    `json:"agent_max_no_answer" bun:"agent_max_no_answer,nullzero"`
	AgentWrapUpTime        int    `json:"agent_wrap_up_time" bun:"agent_wrap_up_time,nullzero"`
	AgentRejectDelayTime   int    `json:"agent_reject_delay_time" bun:"agent_reject_delay_time,nullzero"`
	AgentBusyDelayTime     int    `json:"agent_busy_delay_time" bun:"agent_busy_delay_time,nullzero"`
	AgentNoAnswerDelayTime string `json:"agent_no_answer_delay_time" bun:"agent_no_answer_delay_time,nullzero"`
}

type CallCenterTier struct {
	bun.BaseModel       `bun:"v_call_center_tiers"`
	CallCenterTierUuid  string `json:"call_center_tier_uuid" bun:"call_center_tier_uuid"`
	DomainUuid          string `json:"domain_uuid" bun:"domain_uuid"`
	CallCenterQueueUuid string `json:"call_center_queue_uuid" bun:"call_center_queue_uuid"`
	CallCenterAgentUuid string `json:"call_center_agent_uuid" bun:"call_center_agent_uuid"`
	AgentName           string `json:"agent_name" bun:"agent_name"`
	QueueName           string `json:"queue_name" bun:"queue_name"`
	TierLevel           int    `json:"tier_level" bun:"tier_level,type:numeric"`
	TierPosition        int    `json:"tier_position" bun:"tier_position,type:numeric"`
}

type CallCenterTierWithExtension struct {
	bun.BaseModel       `bun:"v_call_center_tiers,alias:cct"`
	CallCenterTierUuid  string `json:"call_center_tier_uuid" bun:"call_center_tier_uuid"`
	DomainUuid          string `json:"domain_uuid" bun:"domain_uuid"`
	CallCenterQueueUuid string `json:"call_center_queue_uuid" bun:"call_center_queue_uuid"`
	CallCenterAgentUuid string `json:"call_center_agent_uuid" bun:"call_center_agent_uuid"`
	Extension           string `json:"extension" bun:"extension"`
}
