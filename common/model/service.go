package model

import "github.com/uptrace/bun"

type UserData struct {
	bun.BaseModel `bun:"v_users,alias:u"`
	UserUuid      string     `json:"user_uuid" bun:"user_uuid,pk"`
	Username      string     `json:"username" bun:"username"`
	DomainUuid    string     `json:"domain_uuid" bun:"domain_uuid"`
	Level         string     `json:"level" bun:"level"`
	RoleUuid      string     `json:"role_uuid" bun:"role_uuid"`
	UnitUuid      string     `json:"unit_uuid" bun:"unit_uuid"`
	RoleGroups    *RoleGroup `json:"role_groups" bun:"rel:has-one,join:role_uuid=role_group_uuid"`
	Units         *Unit      `json:"units" bun:"rel:has-one,join:unit_uuid=unit_uuid"`
	// Extensions       []*ExtensionData `json:"extensions" bun:"extensions,m2m:v_extension_users"`
	ManageExtensions []ExtensionData `json:"manage_extensions" bun:"-"`
	ManageUsers      []UserInfoData  `json:"manage_users" bun:"-"`
}

type GroupData struct {
	bun.BaseModel  `bun:"v_groups,alias:g"`
	GroupUuid      string `json:"group_uuid" bun:"group_uuid,pk"`
	GroupName      string `json:"group_name" bun:"group_name"`
	GroupWeight    int    `json:"group_weight" bun:"group_weight"`
	DepartmentUuid string `json:"department_uuid" bun:"department_uuid"`
}

type ExtensionData struct {
	bun.BaseModel `bun:"v_extensions,alias:e"`
	ExtensionUuid string `json:"extension_uuid" bun:"extension_uuid,pk"`
	UserUuid      string `json:"user_uuid" bun:"user_uuid,pk"`
	DomainUuid    string `json:"domain_uuid" bun:"domain_uuid"`
	Extension     string `json:"extension" bun:"extension"`
}
type DepartmentData struct {
	bun.BaseModel  `bun:"department,alias:d"`
	DomainUuid     string `json:"domain_uuid" bun:"domain_uuid"`
	DepartmentUuid string `json:"department_uuid" bun:"department_uuid,pk"`
	DepartmentName string `json:"department_name" bun:"department_name"`
}

type UserInfoData struct {
	bun.BaseModel `bun:"v_users,alias:u"`
	UserUuid      string     `json:"user_uuid" bun:"user_uuid,pk"`
	ExtensionUuid string     `json:"extension_uuid" bun:"extension_uuid,pk"`
	DomainUuid    string     `json:"domain_uuid" bun:"domain_uuid"`
	Extension     string     `json:"extension" bun:"extension"`
	Username      string     `json:"username" bun:"username"`
	Firstname     string     `json:"first_name" bun:"first_name"`
	Middlename    string     `json:"middle_name" bun:"middle_name"`
	Lastname      string     `json:"last_name" bun:"last_name"`
	Level         string     `json:"level" bun:"level"`
	Email         string     `json:"email" bun:"email"`
	RoleUuid      string     `json:"role_uuid" bun:"role_uuid"`
	UnitUuid      string     `json:"unit_uuid" bun:"unit_uuid"`
	RoleGroups    *RoleGroup `json:"role_groups" bun:"rel:has-one,join:role_uuid=role_group_uuid"`
	Units         *Unit      `json:"units" bun:"rel:has-one,join:unit_uuid=unit_uuid"`
}
