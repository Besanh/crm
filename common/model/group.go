package model

import "github.com/uptrace/bun"

type Group struct {
	bun.BaseModel    `bun:"v_groups"`
	GroupUuid        string `json:"group_uuid" bun:"group_uuid,pk"`
	DomainUuid       string `json:"domain_uuid" bun:"domain_uuid"`
	GroupName        string `json:"group_name" bun:"group_name"`
	GroupDescription string `json:"group_description" bun:"group_description,nullzero"`
	GroupProtected   string `json:"-" bun:"group_protected"`
	DepartmentUuid   string `json:"department_uuid" bun:"department_uuid,nullzero"`
}

type GroupView struct {
	bun.BaseModel    `bun:"v_groups,alias:g"`
	GroupUuid        string                `json:"group_uuid" bun:"group_uuid,pk"`
	DomainUuid       string                `json:"domain_uuid" bun:"domain_uuid"`
	GroupName        string                `json:"group_name" bun:"group_name"`
	GroupDescription string                `json:"group_description" bun:"group_description,nullzero"`
	GroupPriority    int                   `json:"group_priority" bun:"group_priority"`
	Users            []*UserOptionView     `json:"users" bun:"m2m:v_group_users,join:Group=User"`
	Campaigns        []*CampaignOptionView `json:"campaigns" bun:"m2m:campaign_groups,join:Group=Campaign"`
	DepartmentUuid   string                `json:"department_uuid" pg:"department_uuid"`
	DepartmentName   string                `json:"department_name" pg:"department_name"`
	GroupKPI         *GroupKPI             `json:"kpi" bun:"rel:belongs-to,join:group_uuid=group_uuid"`
}

type GroupViewWithTotalUser struct {
	bun.BaseModel    `bun:"v_groups,alias:g"`
	GroupUuid        string `json:"group_uuid" bun:"group_uuid,pk"`
	DomainUuid       string `json:"domain_uuid" bun:"domain_uuid"`
	DomainName       string `json:"domain_name" bun:"domain_name"`
	GroupName        string `json:"group_name" bun:"group_name"`
	GroupDescription string `json:"group_description" bun:"group_description,nullzero"`
	TotalUsers       int    `json:"total_users" bun:"total_users"`
	DepartmentUuid   string `json:"department_uuid" pg:"department_uuid"`
	DepartmentName   string `json:"department_name" pg:"department_name"`
}

type GroupUser struct {
	bun.BaseModel `bun:"v_group_users"`
	GroupUserUuid string    `json:"group_user_uuid" bun:"group_user_uuid"`
	DomainUuid    string    `json:"domain_uuid" bun:"domain_uuid"`
	GroupName     string    `json:"group_name" bun:"group_name"`
	GroupUuid     string    `json:"group_uuid" bun:"group_uuid,pk"`
	UserUuid      string    `json:"user_uuid" bun:"user_uuid,pk"`
	User          *UserView `bun:"rel:belongs-to,join:user_uuid=user_uuid"`
	Group         *Group    `bun:"rel:belongs-to,join:group_uuid=group_uuid"`
}

type GroupPut struct {
	bun.BaseModel    `bun:"v_groups"`
	GroupUuid        string    `json:"group_uuid"`
	DomainUuid       string    `json:"domain_uuid"`
	GroupName        string    `json:"group_name"`
	GroupDescription string    `json:"group_description"`
	DepartmentUuid   string    `json:"department_uuid"`
	Campaigns        []string  `json:"campaign_assigns"`
	Users            []string  `json:"user_assigns"`
	GroupKPI         *GroupKPI `json:"kpi"`
}

type GroupKPI struct {
	bun.BaseModel `bun:"group_kpi,alias:gk"`
	GroupUuid     string `json:"group_uuid" bun:"group_uuid,type:uuid,pk,notnull"`
	MinCallPerDay int    `json:"min_call_per_day" bun:"min_call_per_day,type:integer,use_zero"`
	MissCallRatio int    `json:"miss_call_ratio" bun:"miss_call_ratio,type:integer,use_zero"`
	WrapUpTime    int    `json:"wrap_up_time" bun:"wrap_up_time,type:integer,use_zero"`
	CallReceive5s int    `json:"call_receive_5s" bun:"call_receive_5s,type:integer,use_zero"`
}
