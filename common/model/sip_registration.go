package model

import "github.com/uptrace/bun"

type SipRegistration struct {
	bun.BaseModel  `bun:"sip_registrations,alias:sr"`
	CallId         string `json:"call_id" bun:"call_id"`
	SipUser        string `json:"sip_user" bun:"sip_user"`
	SipHost        string `json:"sip_host" bun:"sip_host"`
	Status         string `json:"status" bun:"status"`
	SipUsername    string `json:"sip_username" bun:"sip_username"`
	SipRealm       string `json:"sip_realm" bun:"sip_realm"`
	MwiUser        string `json:"mwi_user" bun:"mwi_user"`
	MwiHost        string `json:"mwi_host" bun:"mwi_host"`
	OrigServerHost string `json:"orig_server_host" bun:"orig_server_host"`
	OrigHostname   string `json:"orig_hostname" bun:"orig_hostname"`
	SubHost        string `json:"sub_host" bun:"sub_host"`
}
