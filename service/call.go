package service

import (
	"contactcenter-api/common/cache"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"context"
	"encoding/json"
	"time"
)

type (
	ICallTransfer interface {
		GetCallTransferByInfo(ctx context.Context, domainUuid, phoneNumber string) (int, any)
		CreateCallTransfer(ctx context.Context, domainUuid string, callTransfer model.CallTransfer) (int, any)
	}
	CallTransfer struct{}
)

func NewCall() ICallTransfer {
	return &CallTransfer{}
}

func (service *CallTransfer) GetCallTransferByInfo(ctx context.Context, domainUuid, phoneNumber string) (int, any) {
	dataCache, err := cache.RCache.Get(constants.CALL_TRANSFER + "_" + phoneNumber)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	var data model.CallTransfer
	if dataCache != "" {
		if err := json.Unmarshal([]byte(dataCache), &data); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error)
		}
	}

	return response.OK(data)
}

func (service *CallTransfer) CreateCallTransfer(ctx context.Context, domainUuid string, callTransfer model.CallTransfer) (int, any) {
	err := cache.RCache.SetTTL(constants.CALL_TRANSFER+"_"+callTransfer.PhoneNumber, callTransfer, 60*time.Minute)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error)
	}
	return response.OK(map[string]any{
		"call_id": callTransfer.CallId,
	})
}
