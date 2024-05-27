package model

import (
	"time"

	"github.com/uptrace/bun"
)

type SocialNetwork struct {
	bun.BaseModel `bun:"social_network,alias:sn"`
	DomainUuid    string    `json:"domain_uuid" bun:"domain_uuid,type: char(36)"`
	SocialUuid    string    `json:"social_uuid" bun:"social_uuid,type: char(36),pk"`
	SocialName    string    `json:"social_name" bun:"social_name,type:text"`
	SocialType    string    `json:"social_type" bun:"social_type,type:text"`
	UserUuid      string    `json:"user_uuid" bun:"user_uuid,type: char(36)"`
	Avatar        string    `json:"avatar" bun:"avatar,type:text"`
	ZaloId        string    `json:"zalo_id" bun:"zalo_id,type:text"`
	FacebookId    string    `json:"facebook_id" bun:"facebook_id,type:text"`
	StatusChat    bool      `json:"status_chat" bun:"status_chat,type:boolean,notnull,default:'true'"`
	StatusEmail   bool      `json:"status_email" bun:"status_email,type:boolean,notnull,default:'true'"`
	CreatedBy     string    `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy     string    `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,type:timestamp,nullzero"`
}

type SocialNetworkData struct {
	DomainUuid  string    `json:"domain_uuid" bun:"domain_uuid"`
	SocialUuid  string    `json:"social_uuid" bun:"social_uuid"`
	SocialName  string    `json:"social_name" bun:"social_name"`
	SocialType  string    `json:"social_type" bun:"social_type"`
	UserUuid    string    `json:"user_uuid" bun:"user_uuid"`
	Avatar      string    `json:"avatar" bun:"avatar"`
	ZaloId      string    `json:"zalo_id" bun:"zalo_id"`
	FacebookId  string    `json:"facebook_id" bun:"facebook_id"`
	StatusChat  bool      `json:"status_chat" bun:"status_chat"`
	StatusEmail bool      `json:"status_email" bun:"status_email"`
	Username    string    `json:"username" bun:"username"`
	Level       string    `json:"level" bun:"level"`
	CreatedBy   string    `json:"created_by" bun:"created_by"`
	UpdatedBy   string    `json:"updated_by" bun:"updated_by"`
	CreatedAt   time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bun:"updated_at"`
}

type ConversationResult struct {
	ConvData []Conv `json:"conv_data"`
}

type Conv struct {
	Id string `json:"_id" bun:"_id"`
}

type ConversationBody struct {
	OwnerId string `json:"owner_id"`
}

type SocialNetworkConfig struct {
	Domain string
	Token  string
	PageId string
	ShopId string
}
