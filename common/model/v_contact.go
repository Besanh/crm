package model

import "github.com/uptrace/bun"

type VContact struct {
	bun.BaseModel       `bun:"v_contacts,alias:c"`
	ContactUuid         string `json:"contact_uuid" bun:"contact_uuid"`
	DomainUuid          string `json:"domain_uuid" bun:"domain_uuid,nullzero"`
	ContactParentUuid   string `json:"contact_parent_uuid" bun:"contact_parent_uuid,nullzero"`
	ContactType         string `json:"contact_type" bun:"contact_type,nullzero"`
	ContactOrganization string `json:"contact_organization" bun:"contact_organization,nullzero"`
	ContactNamePrefix   string `json:"contact_name_prefix" bun:"contact_name_prefix,nullzero"`
	ContactNameGiven    string `json:"contact_name_given" bun:"contact_name_given,nullzero"`
	ContactNameMiddle   string `json:"contact_name_middle" bun:"contact_name_middle,nullzero"`
	ContactNameFamily   string `json:"contact_name_family" bun:"contact_name_family,nullzero"`
	ContactNameSuffix   string `json:"contact_name_suffix" bun:"contact_name_suffix,nullzero"`
	ContactNickname     string `json:"contact_nickname" bun:"contact_nickname,nullzero"`
	ContactTitle        string `json:"contact_title" bun:"contact_title,nullzero"`
	ContactRole         string `json:"contact_role" bun:"contact_role,nullzero"`
	ContactCategory     string `json:"contact_category" bun:"contact_category,nullzero"`
	ContactUrl          string `json:"contact_url" bun:"contact_url,nullzero"`
	ContactTimeZone     string `json:"contact_time_zone" bun:"contact_time_zone,nullzero"`
	ContactNote         string `json:"contact_note" bun:"contact_note,nullzero"`
	LastModDate         string `json:"last_mod_date" bun:"last_mod_date,nullzero"`
	LastModUser         string `json:"last_mod_user" bun:"last_mod_user,nullzero"`
}
