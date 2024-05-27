package location

import "github.com/uptrace/bun"

type LocationDistrict struct {
	bun.BaseModel    `bun:"location_district,alias:ld"`
	DistrictCode     string            `json:"district_code" bun:"district_code,type:text,pk,notnull"`
	ProvinceCode     string            `json:"province_code" bun:"province_code,type:text,notnull"`
	DistrictName     string            `json:"district_name" bun:"district_name,type:text,notnull"`
	Status           bool              `json:"status" bun:"status,type:bool,default:false"`
	LocationProvince *LocationProvince `json:"location_province" bun:"rel:belongs-to"`
}
