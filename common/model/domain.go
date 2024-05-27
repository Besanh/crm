package model

import (
	"contactcenter-api/common/variable"
	"errors"
	"slices"

	"github.com/uptrace/bun"
)

type Domain struct {
	bun.BaseModel     `bun:"v_domains"`
	DomainUuid        string `json:"domain_uuid" bun:"domain_uuid,type:uuid,pk"`
	DomainParentUuid  string `json:"domain_parent_uuid" bun:"domain_parent_uuid,type:uuid,nullzero"`
	DomainName        string `json:"domain_name" bun:"domain_name,type:text"`
	DomainEnabled     string `json:"domain_enabled" bun:"domain_enabled,type:text"`
	DomainDescription string `json:"domain_description" bun:"domain_description,type:text"`
}

type DomainConfig struct {
	bun.BaseModel  `bun:"crm_domain_config,alias:cdc"`
	DomainUuid     string             `json:"domain_uuid" bun:"domain_uuid,type:char(36),pk"`
	Version        string             `json:"version" bun:"version,type:text"`
	Logo           string             `json:"logo" bun:"logo,type:text"`
	APIUrl         string             `json:"api_url" bun:"api_url,type:text"`
	Color          string             `json:"color" bun:"color,type:text"`
	CdrStatus      []CdrStatus        `json:"cdr_status" bun:"cdr_status,type:text"`
	MissCallStatus []string           `json:"miss_call_status" bun:"miss_call_status,array"`
	Setting        DomainFieldSetting `json:"setting" bun:"setting,type:jsonb"`
}

type DomainConfigPut struct {
	bun.BaseModel  `bun:"crm_domain_config,alias:cdc"`
	DomainUuid     string             `json:"domain_uuid" bun:"domain_uuid,type:char(36),pk"`
	Version        string             `json:"version" bun:"version,type:text"`
	Logo           string             `json:"logo" bun:"logo,type:text"`
	APIUrl         string             `json:"api_url" bun:"api_url,type:text"`
	Color          string             `json:"color" bun:"color,type:text"`
	MissCallStatus []string           `json:"miss_call_status" bun:"miss_call_status,array"`
	Setting        DomainFieldSetting `json:"setting" bun:"setting,type:jsonb"`
}

type DomainConfigView struct {
	bun.BaseModel  `bun:"v_domains,alias:d"`
	DomainUuid     string             `json:"domain_uuid" bun:"domain_uuid"`
	DomainName     string             `json:"domain_name" bun:"domain_name"`
	Version        string             `json:"version" bun:"version"`
	Logo           string             `json:"logo" bun:"logo"`
	APIUrl         string             `json:"api_url" bun:"api_url"`
	Color          string             `json:"color" bun:"color"`
	CdrStatus      []CdrStatus        `json:"cdr_status" bun:"cdr_status"`
	MissCallStatus []string           `json:"miss_call_status" bun:"miss_call_status,array"`
	Setting        DomainFieldSetting `json:"setting" bun:"setting,type:jsonb"`
}

type CdrStatus struct {
	UnitUuid string   `json:"unit_uuid"`
	PressKey string   `json:"press_key"`
	Status   []string `json:"status"`
}

type DomainFieldSetting struct {
	ExternalPluginUuid string `json:"external_plugin_uuid"`
}

func (m *DomainConfigPut) ValidatePut() (err error) {
	if len(m.MissCallStatus) > 0 {
		for _, item := range m.MissCallStatus {
			if !slices.Contains(variable.CDRSTATUS, item) {
				return errors.New("status " + item + " invalid")
			}
		}
	}

	return
}
