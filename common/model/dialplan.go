package model

import "github.com/uptrace/bun"

type Dialplan struct {
	bun.BaseModel       `bun:"v_dialplans"`
	DomainUuid          string `json:"domain_uuid" bun:"domain_uuid"`
	DialplanUuid        string `json:"dialplan_uuid" bun:"dialplan_uuid"`
	AppUuid             string `json:"app_uuid" bun:"app_uuid"`
	Hostname            string `json:"hostname" bun:"hostname"`
	DialplanContext     string `json:"dialplan_context" bun:"dialplan_context"`
	DialplanName        string `json:"dialplan_name" bun:"dialplan_name"`
	DialplanNumber      string `json:"dialplan_number" bun:"dialplan_number"`
	DialplanContinue    string `json:"dialplan_continue" bun:"dialplan_continue"`
	DialplanXml         string `json:"dialplan_xml" bun:"dialplan_xml"`
	DialplanOrder       int    `json:"dialplan_order" bun:"dialplan_order"`
	DialplanEnabled     string `json:"dialplan_enabled" bun:"dialplan_enabled"`
	DialplanDescription string `json:"dialplan_description" bun:"dialplan_description"`
}

type DialplanDetail struct {
	bun.BaseModel        `bun:"v_dialplan_details"`
	DomainUuid           string `json:"domain_uuid" bun:"domain_uuid"`
	DialplanUuid         string `json:"dialplan_uuid" bun:"dialplan_uuid"`
	DialplanDetailUuid   string `json:"dialplan_detail_uuid" bun:"dialplan_detail_uuid"`
	DialplanDetailTag    string `json:"dialplan_detail_tag" bun:"dialplan_detail_tag"`
	DialplanDetailType   string `json:"dialplan_detail_type" bun:"dialplan_detail_type"`
	DialplanDetailData   string `json:"dialplan_detail_data" bun:"dialplan_detail_data"`
	DialplanDetailBreak  string `json:"dialplan_detail_break" bun:"dialplan_detail_break"`
	DialplanDetailInline string `json:"dialplan_detail_inline" bun:"dialplan_detail_inline"`
	DialplanDetailGroup  int    `json:"dialplan_detail_group" bun:"dialplan_detail_group"`
	DialplanDetailOrder  int    `json:"dialplan_detail_order" bun:"dialplan_detail_order"`
}
