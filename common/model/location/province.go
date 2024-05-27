package location

import "github.com/uptrace/bun"

type LocationProvince struct {
	bun.BaseModel `bun:"location_province,alias:lp"`
	ProvinceCode  string `json:"province_code" bun:"province_code,type:text,pk,notnull"`
	ProvinceName  string `json:"province_name" bun:"province_name,type:text,notnull"`
	Status        bool   `json:"status" bun:"status,type:bool,default:false"`
}
