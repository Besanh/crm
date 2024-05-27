package model

type FSCallCenterMember struct {
	Queue          string `json:"queue"`
	InstanceId     string `json:"instance_id"`
	Uuid           string `json:"uuid"`
	SessionUuid    string `json:"session_uuid"`
	CidNumber      string `json:"cid_number"`
	CidName        string `json:"cid_name"`
	SystemEpoch    string `json:"system_epoch"`
	JoinedEpoch    string `json:"joined_epoch"`
	RejoinedEpoch  string `json:"rejoined_epoch"`
	BridgeEpoch    string `json:"bridge_epoch"`
	AbandonedEpoch string `json:"abandoned_epoch"`
	BaseScore      string `json:"base_score"`
	SkillScore     string `json:"skill_score"`
	ServingAgent   string `json:"serving_agent"`
	ServingSystem  string `json:"serving_system"`
	State          string `json:"state"`
	Score          string `json:"score"`
	QueueName      string `json:"-"`
}

type FSCallCenterAgent struct {
	BusyDelayTime      string `json:"busy_delay_time"`
	CallsAnswered      string `json:"calls_answered"`
	Contact            string `json:"contact"`
	ExternalCallsCount string `json:"external_calls_count"`
	InstanceId         string `json:"instance_id"`
	LastBridgeEnd      string `json:"last_bridge_end"`
	LastBridgeStart    string `json:"last_bridge_start"`
	LastOfferedCall    string `json:"last_offered_call"`
	LastStatusChange   string `json:"last_status_change"`
	MaxNoAnswer        string `json:"max_no_answer"`
	Name               string `json:"name"`
	NoAnswerCount      string `json:"no_answer_count"`
	NoAnswerDelayTime  string `json:"no_answer_delay_time"`
	ReadyTime          string `json:"ready_time"`
	RejectDelayTime    string `json:"reject_delay_time"`
	State              string `json:"state"`
	Status             string `json:"status"`
	TalkTime           string `json:"talk_time"`
	Type               string `json:"type"`
	Uuid               string `json:"uuid"`
	WrapUpTime         string `json:"wrap_up_time"`
	QueueName          string
}
