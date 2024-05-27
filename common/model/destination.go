package model

import "github.com/uptrace/bun"

type Destination struct {
	bun.BaseModel             `bun:"v_destinations"`
	DomainUuid                string `json:"domain_uuid" bun:"domain_uuid"`
	DestinationUuid           string `json:"destination_uuid" bun:"destination_uuid"`
	DialplanUuid              string `json:"dialplan_uuid" bun:"dialplan_uuid"`
	FaxUuid                   string `json:"fax_uuid" bun:"fax_uuid,nullzero"`
	DestinationType           string `json:"destination_type" bun:"destination_type"`
	DestinationNumber         string `json:"destination_number" bun:"destination_number"`
	DestinationNumberRegex    string `json:"destination_number_regex" bun:"destination_number_regex"`
	DestinationCallerIdName   string `json:"destination_caller_id_name" bun:"destination_caller_id_name,nullzero"`
	DestinationCallerIdNumber string `json:"destination_caller_id_number" bun:"destination_caller_id_number,nullzero"`
	DestinationCidNamePrefix  string `json:"destination_cid_name_prefix" bun:"destination_cid_name_prefix,nullzero"`
	DestinationContext        string `json:"destination_context" bun:"destination_context"`
	DestinationRecord         string `json:"destination_record" bun:"destination_record,nullzero"`
	DestinationAccountcode    string `json:"destination_accountcode" bun:"destination_accountcode,nullzero"`
	DestinationApp            string `json:"destination_app" bun:"destination_app"`
	DestinationData           string `json:"destination_data" bun:"destination_data"`
	DestinationEnabled        string `json:"destination_enabled" bun:"destination_enabled"`
	DestinationDescription    string `json:"destination_description" bun:"destination_description,nullzero"`
}
