package model

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"v_users,alias:u"`
	UserUuid      string `json:"user_uuid" bun:"user_uuid,pk"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid"`
	Username      string `json:"username" bun:"username"`
	Password      string `json:"password" bun:"password"`
	Salt          string `json:"-" bun:"salt"`
	ApiKey        string `json:"api_key" bun:"api_key"`
	UserEnabled   string `json:"user_enabled" bun:"user_enabled"`
	ContactUuid   string `json:"contact_uuid" bun:"contact_uuid"`
	UserStatus    string `json:"user_status" bun:"user_status"`
	AddDate       string `json:"add_date" bun:"add_date"`
	AddUser       string `json:"add_user" bun:"add_user"`
	Level         string `json:"level" bun:"level"`
	EnableWebrtc  bool   `json:"enable_webrtc" bun:"enable_webrtc"`
	UnitUuid      string `json:"unit_uuid" bun:"unit_uuid"`
	RoleUuid      string `json:"role_uuid" bun:"role_uuid"`
}
type UserOwnerShortInfo struct {
	bun.BaseModel `bun:"v_users,alias:u"`
	Data          string `json:"data" bun:"data"`
	Type          string `json:"type" bun:"type"`
	Info          struct {
		UserUuid   string    `json:"user_uuid" bun:"user_uuid,pk"`
		Username   string    `json:"username" bun:"username"`
		AssignedAt time.Time `json:"assigned_at" bun:"assigned_at"`
	} `json:"info" bun:"info"`
}

type UserBasicInfo struct {
	bun.BaseModel `bun:"v_users,alias:u"`
	UserUuid      string `json:"user_uuid" bun:"user_uuid,pk"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid"`
	Username      string `json:"username" bun:"username"`
}

type UserInfo struct {
	bun.BaseModel    `bun:"v_users,alias:u"`
	UserUuid         string          `json:"user_uuid" bun:"user_uuid,pk"`
	Username         string          `json:"username" bun:"username"`
	DomainUuid       string          `json:"domain_uuid" bun:"domain_uuid"`
	Level            string          `json:"level" bun:"level"`
	ManageExtensions []ExtensionData `json:"manage_extensions" bun:"-"`
	ManageUsers      []UserInfoData  `json:"manage_users" bun:"-"`
	RoleUuid         string          `json:"role_uuid" bun:"role_uuid"`
	UnitUuid         string          `json:"unit_uuid" bun:"unit_uuid"`
	RoleGroups       *RoleGroup      `json:"role_groups" bun:"rel:has-one,join:role_uuid=role_group_uuid"`
	Units            *Unit           `json:"units" bun:"rel:has-one,join:unit_uuid=unit_uuid"`
}

type UserView struct {
	bun.BaseModel       `bun:"v_users,alias:u"`
	UserUuid            string          `json:"user_uuid" bun:"user_uuid,pk"`
	DomainUuid          string          `json:"domain_uuid" bun:"domain_uuid"`
	Username            string          `json:"username" bun:"username"`
	Password            string          `json:"password" bun:"password"`
	ApiKey              string          `json:"api_key" bun:"api_key"`
	UserEnabled         string          `json:"user_enabled" bun:"user_enabled"`
	UserStatus          string          `json:"user_status" bun:"user_status"`
	Level               string          `json:"level" bun:"level"`
	LastName            string          `json:"last_name" bun:"last_name"`
	MiddleName          string          `json:"middle_name" bun:"middle_name"`
	FirstName           string          `json:"first_name" bun:"first_name"`
	UnitUuid            string          `json:"unit_uuid" bun:"unit_uuid"`
	UnitName            string          `json:"unit_name" bun:"unit_name"`
	RoleUuid            string          `json:"role_uuid" bun:"role_uuid"`
	RoleName            string          `json:"role_group_name" bun:"role_group_name"`
	Email               string          `json:"email" bun:"email"`
	EnableWebrtc        bool            `json:"enable_webrtc" bun:"enable_webrtc"`
	Extension           string          `json:"extension" bun:"extension"`
	ExtensionUuid       string          `json:"extension_uuid" bun:"extension_uuid"`
	RoleGroups          *RoleGroup      `json:"role_groups" bun:"rel:has-one,join:role_uuid=role_group_uuid"`
	Units               *Unit           `json:"units" bun:"rel:has-one,join:unit_uuid=unit_uuid"`
	ManageExtensions    []ExtensionData `json:"manage_extensions" bun:"-"`
	ManageUsers         []*UserView     `json:"manage_users" bun:"-"`
	ManageUnits         []*UnitInfo     `json:"manage_units" bun:"-"`
	ManageCampaignUuids []string        `json:"manage_campaign_uuids" bun:"-"`
}

type UserLoginView struct {
	bun.BaseModel `bun:"v_users,alias:u"`
	UserUuid      string `json:"user_uuid" bun:"user_uuid,pk"`
	ExtensionUuid string `json:"extension_uuid" bun:"extension_uuid"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid"`
	DomainName    string `json:"domain_name" bun:"domain_name"`
	Username      string `json:"username" bun:"username"`
	Extension     string `json:"extension" bun:"extension"`
}

type UserPost struct {
	UnitUuid          string `json:"unit_uuid" bun:"unit_uuid"`
	GroupUuid         string `json:"group_uuid"`
	RoleUuid          string `json:"role_uuid"`
	LastName          string `json:"last_name"`
	MiddleName        string `json:"middle_name"`
	FirstName         string `json:"first_name"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	Level             string `json:"level"`
	IsCreateExtension bool   `json:"is_create_extension"`
	Extension         string `json:"extension"`
	UserEnabled       string `json:"user_enabled"`
	Email             string `json:"email"`
	EnableWebrtc      bool   `json:"enable_webrtc"`
	IsRemapping       bool   `json:"is_remapping"`
}
type UserLive struct {
	bun.BaseModel  `bun:"user_live,alias:ul"`
	Id             string    `json:"user_live_uuid" bun:"user_live_uuid,pk,type:char(36)"`
	UserUuid       string    `json:"user_uuid" bun:"user_uuid,type:char(36),unique"`
	Username       string    `json:"username" bun:"username,type:text"`
	Extension      string    `json:"extension" bun:"extension,type:text"`
	ExtensionUuid  string    `json:"extension_uuid" bun:"extension_uuid,type:char(36)"`
	DomainUuid     string    `json:"domain_uuid" bun:"domain_uuid,type:char(36)"`
	DomainName     string    `json:"domain_name" bun:"domain_name,type:varchar(20)"`
	CampaignUuid   string    `json:"campaign_uuid" bun:"campaign_uuid,type:char(36)"`
	Status         string    `json:"status" bun:"status,type:varchar(20)"`
	LoginTime      time.Time `json:"login_time" bun:"login_time,type:timestamp"`
	LastUpdateTime time.Time `json:"last_update_time" bun:"last_update_time,type:timestamp"`
}

type UserLiveView struct {
	bun.BaseModel  `bun:"user_live,alias:ul"`
	Id             string    `json:"user_live_uuid" bun:"user_live_uuid,pk"`
	UserUuid       string    `json:"user_uuid" bun:"user_uuid"`
	Username       string    `json:"username" bun:"username"`
	Extension      string    `json:"extension" bun:"extension"`
	ExtensionUuid  string    `json:"extension_uuid" bun:"extension_uuid"`
	DomainUuid     string    `json:"domain_uuid" bun:"domain_uuid"`
	DomainName     string    `json:"domain_name" bun:"domain_name"`
	CampaignUuid   string    `json:"campaign_uuid" bun:"campaign_uuid"`
	Status         string    `json:"status" bun:"status"`
	LoginTime      time.Time `json:"login_time" bun:"login_time"`
	LastUpdateTime time.Time `json:"last_update_time" bun:"last_update_time"`
	LastName       string    `json:"last_name" bun:"last_name"`
	MiddleName     string    `json:"middle_name" bun:"middle_name"`
	FirstName      string    `json:"first_name" bun:"first_name"`
}

type UserExtensionView struct {
	bun.BaseModel `bun:"v_users,alias:u"`
	UserUuid      string `json:"user_uuid" bun:"user_uuid,pk"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid"`
	Username      string `json:"username" bun:"username"`
	LastName      string `json:"last_name" bun:"last_name"`
	MiddleName    string `json:"middle_name" bun:"middle_name"`
	FirstName     string `json:"first_name" bun:"first_name"`
	Extension     string `json:"extension" bun:"extension"`
}

type UserOptionView struct {
	bun.BaseModel `bun:"v_users,alias:u"`
	UserUuid      string `json:"user_uuid" bun:"user_uuid,pk"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid"`
	Level         string `json:"level" bun:"level"`
	Username      string `json:"username" bun:"username"`
	Firstname     string `json:"first_name" bun:"first_name"`
	Middlename    string `json:"middle_name" bun:"middle_name"`
	Lastname      string `json:"last_name" bun:"last_name"`
}

type UserPasswordPut struct {
	CurrentPassword string `json:"current_password"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UserAuth struct {
	bun.BaseModel `bun:"v_users,alias:u"`
	UserUuid      string `json:"user_uuid" bun:"user_uuid,pk"`
	DomainUuid    string `json:"domain_id" bun:"domain_uuid"`
	DomainName    string `json:"domain_name" bun:"domain_name"`
	Username      string `json:"username" bun:"username"`
	Password      string `json:"password" bun:"password"`
	Salt          string `json:"salt" bun:"salt"`
	ApiKey        string `json:"api_key" bun:"api_key"`
	UserEnabled   string `json:"user_enabled" bun:"user_enabled"`
	Level         string `json:"level" bun:"level"`
}

type UserCustomData struct {
	bun.BaseModel      `bun:"user_custom_data"`
	UserCustomDataUuid string `json:"user_custom_data_uuid" bun:"user_custom_data_uuid,pk,type:char(36)"`
	UserUuid           string `json:"user_uuid" bun:"user_uuid,type:char(36),notnull,unique:unq_user_custom_data_key"`
	Key                string `json:"key" bun:"key,type:text,notnull,unique:unq_user_custom_data_key"`
	Value              string `json:"value" bun:"value,type:text,nullzero"`
}

type UserCampaignView struct {
	bun.BaseModel `bun:"v_users,alias:u"`
	UserUuid      string `json:"user_uuid" bun:"user_uuid,pk"`
	Extension     string `json:"extension" bun:"extension"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid"`
}
