package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Transaction struct {
	bun.BaseModel   `bun:"transaction,alias:t"`
	DomainUuid      string    `json:"domain_uuid" bun:"domain_uuid,type:uuid,notnull"`
	TransactionUuid string    `json:"transaction_uuid" bun:"transaction_uuid,pk,type:uuid,notnull"`
	Entity          string    `json:"entity" bun:"entity,type:text"`
	EntityUuid      string    `json:"entity_uuid" bun:"entity_uuid,type:text"`
	Action          string    `json:"action" bun:"action,type:text"` //add update delete import export
	Status          string    `json:"status" bun:"status,type:text"`
	OldData         string    `json:"old_data" bun:"old_data,type:text"`
	NewData         string    `json:"new_data" bun:"new_data,type:text"`
	Result          string    `json:"result" bun:"result,type:text"`
	UserUuid        string    `json:"user_uuid" bun:"user_uuid,type:text,nullzero"`
	CreatedAt       time.Time `json:"created_at" bun:"created_at,type:timestamp,default:now()"`
}

type TransactionLog struct {
	bun.BaseModel `bun:"transaction_log,alias:tl"`
	Id            string    `json:"id" bun:"id,type:uuid,pk"`
	Level         string    `json:"level" bun:"level,type:text"`
	Message       string    `json:"message" bun:"message,type:text"`
	Data          string    `json:"data" bun:"data,type:text"`
	EntityId      string    `json:"entity_id" bun:"entity_id,type:text"`
	Meta          string    `json:"meta" bun:"meta,type:text"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,type:timestamp,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,type:timestamp,nullzero,notnull,default:current_timestamp"`
}
