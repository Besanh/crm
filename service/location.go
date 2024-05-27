package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/response"
	"contactcenter-api/repository"
	"context"
)

type (
	ILocation interface {
		GetLocationProvince(ctx context.Context, limit, offset int) (int, any)
		GetLocationDistrict(ctx context.Context, provinceCode string, limit, offset int) (int, any)
		GetLocationWard(ctx context.Context, districtCode string, limit, offset int) (int, any)

		GetLocationWardByName(ctx context.Context, wardName string, limit, offset int) (int, any)
	}
	Location struct{}
)

func NewLocation() ILocation {
	return &Location{}
}
func (s *Location) GetLocationProvince(ctx context.Context, limit, offset int) (int, any) {
	total, provinces, err := repository.LocationProvinceRepo.GetProvinceByNameOrCode(ctx, "", limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Pagination(provinces, total, limit, offset)
}

func (s *Location) GetLocationDistrict(ctx context.Context, provinceCode string, limit, offset int) (int, any) {
	provinceExist, err := repository.LocationProvinceRepo.GetProvinceByCode(ctx, provinceCode)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(provinceExist.ProvinceCode) < 1 {
		return response.ServiceUnavailableMsg("province is not exist")
	}
	total, districts, err := repository.LocationDistrictRepo.GetDistrictByEntity(ctx, provinceExist.ProvinceCode, "", limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Pagination(districts, total, limit, offset)
}

func (s *Location) GetLocationWard(ctx context.Context, districtCode string, limit, offset int) (int, any) {
	districtExist, err := repository.LocationDistrictRepo.GetDistrictByCode(ctx, districtCode)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(districtExist.DistrictCode) < 1 {
		return response.ServiceUnavailableMsg("province is not exist")
	}
	total, wards, err := repository.LocationWardRepo.GetWardByEntity(ctx, districtExist.DistrictCode, "", limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Pagination(wards, total, limit, offset)
}

func (s *Location) GetLocationWardByName(ctx context.Context, wardName string, limit, offset int) (int, any) {
	total, wards, err := repository.LocationWardRepo.GetWardByEntity(ctx, "", wardName, limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Pagination(wards, total, limit, offset)
}
