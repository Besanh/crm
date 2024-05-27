package model

import "github.com/uptrace/bun"

type Extension struct {
	bun.BaseModel                       `bun:"v_extensions,alias:e"`
	ExtensionUuid                       string `json:"extension_uuid" bun:"extension_uuid,pk"`
	DomainUuid                          string `json:"domain_id" bun:"domain_uuid"`
	Extension                           string `json:"extension" bun:"extension,nullzero"`
	NumberAlias                         string `json:"number_alias" bun:"number_alias,nullzero"`
	Password                            string `json:"password" bun:"password,nullzero"`
	Accountcode                         string `json:"accountcode" bun:"accountcode,nullzero"`
	EffectiveCallerIdName               string `json:"effective_caller_id_name" bun:"effective_caller_id_name,nullzero"`
	EffectiveCallerIdNumber             string `json:"effective_caller_id_number" bun:"effective_caller_id_number,nullzero"`
	OutboundCallerIdName                string `json:"outbound_caller_id_name" bun:"outbound_caller_id_name,nullzero"`
	OutboundCallerIdNumber              string `json:"outbound_caller_id_number" bun:"outbound_caller_id_number,nullzero"`
	EmergencyCallerIdName               string `json:"emergency_caller_id_name" bun:"emergency_caller_id_name,nullzero"`
	EmergencyCallerIdNumber             string `json:"emergency_caller_id_number" bun:"emergency_caller_id_number,nullzero"`
	DirectoryFirstName                  string `json:"directory_first_name" bun:"directory_first_name,nullzero"`
	DirectoryLastName                   string `json:"directory_last_name" bun:"directory_last_name,nullzero"`
	DirectoryVisible                    string `json:"directory_visible" bun:"directory_visible,nullzero"`
	DirectoryExtenVisible               string `json:"directory_exten_visible" bun:"directory_exten_visible,nullzero"`
	LimitMax                            string `json:"limit_max" bun:"limit_max,nullzero"`
	LimitDestination                    string `json:"limit_destination" bun:"limit_destination,nullzero"`
	MissedCallApp                       string `json:"missed_call_app" bun:"missed_call_app,nullzero"`
	MissedCallData                      string `json:"missed_call_data" bun:"missed_call_data,nullzero"`
	UserContext                         string `json:"user_context" bun:"user_context,nullzero"`
	TollAllow                           string `json:"toll_allow" bun:"toll_allow,nullzero"`
	CallTimeout                         int    `json:"call_timeout" bun:"call_timeout,nullzero"`
	CallGroup                           string `json:"call_group" bun:"call_group,nullzero"`
	CallScreenEnabled                   string `json:"call_screen_enabled" bun:"call_screen_enabled,nullzero"`
	UserRecord                          string `json:"user_record" bun:"user_record,nullzero"`
	HoldMusic                           string `json:"hold_music" bun:"hold_music,nullzero"`
	AuthAcl                             string `json:"auth_acl" bun:"auth_acl,nullzero"`
	Cidr                                string `json:"cidr" bun:"cidr,nullzero"`
	SipForceContact                     string `json:"sip_force_contact" bun:"sip_force_contact,nullzero"`
	NibbleAccount                       int    `json:"nibble_account" bun:"nibble_account,nullzero"`
	SipForceExpires                     int    `json:"sip_force_expires" bun:"sip_force_expires,nullzero"`
	MwiAccount                          string `json:"mwi_account" bun:"mwi_account,nullzero"`
	SipBypassMedia                      string `json:"sip_bypass_media" bun:"sip_bypass_media,nullzero"`
	UniqueId                            int    `json:"unique_id" bun:"unique_id,nullzero"`
	DialString                          string `json:"dial_string" bun:"dial_string,nullzero"`
	DialUser                            string `json:"dial_user" bun:"dial_user,nullzero"`
	DialDomain                          string `json:"dial_domain" bun:"dial_domain,nullzero"`
	DoNotDisturb                        string `json:"do_not_disturb" bun:"do_not_disturb,nullzero"`
	ForwardAllDestination               string `json:"forward_all_destination" bun:"forward_all_destination,nullzero"`
	ForwardAllEnabled                   string `json:"forward_all_enabled" bun:"forward_all_enabled,nullzero"`
	ForwardBusyDestination              string `json:"forward_busy_destination" bun:"forward_busy_destination,nullzero"`
	ForwardBusyEnabled                  string `json:"forward_busy_enabled" bun:"forward_busy_enabled,nullzero"`
	ForwardNoAnswerDestination          string `json:"forward_no_answer_destination" bun:"forward_no_answer_destination,nullzero"`
	ForwardNoAnswerEnabled              string `json:"forward_no_answer_enabled" bun:"forward_no_answer_enabled,nullzero"`
	ForwardUserNotRegisteredDestination string `json:"forward_user_not_registered_destination" bun:"forward_user_not_registered_destination,nullzero"`
	ForwardUserNotRegisteredEnabled     string `json:"forward_user_not_registered_enabled" bun:"forward_user_not_registered_enabled,nullzero"`
	FollowMeUuid                        string `json:"follow_me_uuid" bun:"follow_me_uuid,nullzero"`
	Enabled                             string `json:"enabled" bun:"enabled,nullzero"`
	Description                         string `json:"description" bun:"description,nullzero"`
	ForwardCallerIdUuid                 string `json:"forward_caller_id_uuid" bun:"forward_caller_id_uuid,nullzero"`
	AbsoluteCodecString                 string `json:"absolute_codec_string" bun:"absolute_codec_string,nullzero"`
	ForcePing                           string `json:"force_ping" bun:"force_ping,nullzero"`
}

type ExtensionView struct {
	bun.BaseModel `bun:"v_extensions,alias:e"`
	ExtensionUuid string `json:"extension_uuid" bun:"extension_uuid,pk"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid"`
	DomainName    string `json:"domain_name" bun:"domain_name"`
	Extension     string `json:"extension" bun:"extension"`
	Enabled       string `json:"enabled" bun:"enabled"`
}

type ExtensionInfo struct {
	bun.BaseModel    `bun:"v_extensions,alias:e"`
	ExtensionUuid    string `json:"extension_uuid" bun:"extension_uuid"`
	Extension        string `json:"extension" bun:"extension,nullzero"`
	UserUuid         string `json:"user_uuid" bun:"user_uuid,nullzero"`
	Username         string `json:"username" bun:"username,nullzero"`
	FirsName         string `json:"first_name" bun:"first_name"`
	LastName         string `json:"last_name" bun:"last_name"`
	MiddleName       string `json:"middle_name" bun:"middle_name"`
	UnitUuid         string `json:"unit_uuid" bun:"unit_uuid"`
	UnitName         string `json:"unit_name" bun:"unit_name"`
	Enabled          bool   `json:"enabled" bun:"enabled,nullzero"`
	DomainUuid       string `json:"domain_uuid" bun:"domain_uuid"`
	DomainName       string `json:"domain_name" bun:"domain_name"`
	IsLinkCallCenter bool   `json:"is_link_call_center" bun:"is_link_call_center"`
	// FollowMeUuid     string `json:"-" bun:"follow_me_uuid"`
	// FollowMe         string `json:"follow_me" bun:"-"`
}

type ExtensionInfoWithPassword struct {
	bun.BaseModel    `bun:"v_extensions,alias:e"`
	ExtensionUuid    string `json:"extension_uuid" bun:"extension_uuid,pk"`
	DomainUuid       string `json:"domain_uuid" bun:"domain_uuid"`
	DomainName       string `json:"domain_name" bun:"domain_name"`
	Extension        string `json:"extension" bun:"extension"`
	Password         string `json:"password" bun:"password"`
	Enabled          bool   `json:"enabled" bun:"enabled,nullzero"`
	IsLinkCallCenter bool   `json:"is_link_call_center" bun:"is_link_call_center"`
	UserUuid         string `json:"user_uuid" bun:"user_uuid,nullzero"`
	Username         string `json:"username" bun:"username,nullzero"`
	// FollowMeUuid     string `json:"-" bun:"follow_me_uuid"`
	FollowMe      string `json:"follow_me" bun:"-"`
	OutboundProxy string `json:"outbound_proxy" bun:"-"`
	Wss           string `json:"wss" bun:"-"`
	SipPort       string `json:"sip_port" bun:"-"`
	Transport     string `json:"transport" bun:"-"`
}

type ExtensionUser struct {
	bun.BaseModel     `bun:"v_extension_users"`
	ExtensionUserUuid string `json:"extension_user_uuid" bun:"extension_user_uuid"`
	DomainUuid        string `json:"domain_uuid" bun:"domain_uuid,nullzero"`
	ExtensionUuid     string `json:"extension_uuid" bun:"extension_uuid,nullzero"`
	UserUuid          string `json:"user_uuid" bun:"user_uuid,nullzero"`
}

type ExtensionPost struct {
	Extension       string `json:"extension"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Enabled         bool   `json:"enabled"`
	UserUuid        string `json:"user_uuid"`
	// NEW
	// IsFollowMe        bool            `json:"is_follow_me"`
	// IsDeleteFollowMe  bool            `json:"-"`
	// FollowMe          FollowMeConfig  `json:"follow_me"`
	IsRingGroup       bool            `json:"is_ring_group"`
	RingGroup         RingGroupConfig `json:"ring_group"`
	IsDeleteRingGroup bool            `json:"-"`
}

type FollowMeConfig struct {
	Destination string `json:"destination"`
	Delay       int    `json:"delay"`
	Timeout     int    `json:"timeout"`
	Confirm     bool   `json:"confirm"`
}

type RingGroupConfig struct {
	Destination string `json:"destination"`
	Main        string `json:"main"`
	Sub         string `json:"sub"`
	Script      string `json:"script"`
}
