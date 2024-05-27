package model

import (
	"time"

	"github.com/uptrace/bun"
)

type RoleGroup struct {
	bun.BaseModel             `bun:"role_group,alias:rg"`
	DomainUuid                string                    `json:"domain_uuid" bun:"domain_uuid,type: char(36),notnull"`
	RoleGroupUuid             string                    `json:"role_group_uuid" bun:"role_group_uuid,pk,type: char(36),notnull"`
	RoleGroupName             string                    `json:"role_group_name" bun:"role_group_name,type:text,notnull"`
	Status                    bool                      `json:"status" bun:"status,type:boolean"`
	CharacteristicColor       string                    `json:"characteristic_color" bun:"characteristic_color,type:text"`
	Description               string                    `json:"description" bun:"description,type:text"`
	PermissionMain            Permission                `json:"permission_main" bun:"permission_main,type:text"`
	PermissionAdvance         PermissionAdvance         `json:"permission_advance" bun:"permission_advance,type:text"`
	PermissionUser            PermissionUser            `json:"permission_user" bun:"permission_user,type:text"`
	PermissionMainOptimize    PermissionMainOptimize    `json:"permission_main_optimize" bun:"-"`
	PermissionAdvanceOptimize PermissionAdvanceOptimize `json:"permission_advance_optimize" bun:"-"`
	CreatedBy                 string                    `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy                 string                    `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt                 time.Time                 `json:"created_at" bun:"created_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt                 time.Time                 `json:"updated_at" bun:"updated_at,type:timestamp,notnull,nullzero,default:current_timestamp"`
}
