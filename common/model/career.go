package model

import "github.com/uptrace/bun"

type Career struct {
	bun.BaseModel `bun:"career"`
	CareerCode    string `json:"career_code" bun:"career_code,type:text,notnull"`
	CareerName    string `json:"career_name" bun:"career_name,type:text,notnull"`
	Source        string `json:"source" bun:"source,type:text"`
	Status        bool   `json:"status" bun:"status,type:bool,default:false"`
}

type CareerView struct {
	bun.BaseModel `bun:"career"`
	CareerCode    string `json:"career_code" bun:"career_code"`
	CareerName    string `json:"career_name" bun:"career_name"`
	Status        bool   `json:"status" bun:"status"`
}
