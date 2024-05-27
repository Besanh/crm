package service

import (
	"contactcenter-api/common/constants"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	IRedis "contactcenter-api/internal/redis"
	"net/http"
	"sort"
)

type (
	IExport interface {
		GetExports(domainUuid, userUuid string) (int, any)
		DownloadExport(domainUuid, id string) (int, any)
	}
	Export struct {
	}
)

func NewExport() IExport {
	return &Export{}
}

func (service *Export) GetExports(domainUuid, userUuid string) (int, any) {
	res, err := IRedis.Redis.HGetAll(constants.EXPORT_KEY + userUuid)
	if err != nil {
		return response.ServiceUnavailableMsg(err)
	}

	exports := make(Exports, 0)
	index := 1
	for _, v := range res {
		exportMap := model.ExportMap{}
		err := util.ParsesStringToStruct(v, &exportMap)
		if err != nil {
			return response.ServiceUnavailable()
		}
		exports = append(exports, &exportMap)
		index++
	}
	sort.Sort(exports)
	return response.Data(http.StatusOK, exports)
}

func (service *Export) DownloadExport(userUuid, id string) (int, any) {
	res, err := IRedis.Redis.HGet(constants.EXPORT_KEY+userUuid, id)
	if err != nil {
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(res)
}
