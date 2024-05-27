package location

type Location struct {
	ProvinceCode string `json:"province_code" bun:"province_code,type:text,notnull"`
	ProvinceName string `json:"province_name" bun:"province_name,type:text,notnull"`
	DistrictCode string `json:"district_code" bun:"district_code,type:text,notnull"`
	DistrictName string `json:"district_name" bun:"district_name,type:text,notnull"`
	WardCode     string `json:"ward_code" bun:"ward_code,type:text,notnull"`
	WardName     string `json:"ward_name" bun:"ward_name,type:text,notnull"`
	Status       bool   `json:"status" bun:"status,type:bool,default:false"`
}
