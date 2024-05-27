package location

import "github.com/uptrace/bun"

type LocationWard struct {
	bun.BaseModel    `bun:"location_ward,alias:lw"`
	WardCode         string            `json:"ward_code" bun:"ward_code,type:text,pk,notnull"`
	DistrictCode     string            `json:"district_code" bun:"district_code,type:text,notnull"`
	WardName         string            `json:"ward_name" bun:"ward_name,type:text,notnull"`
	Status           bool              `json:"status" bun:"status,type:bool,default:false"`
	LocationDistrict *LocationDistrict `json:"location_district" bun:"rel:belongs-to"`
}
