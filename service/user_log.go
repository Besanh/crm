package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/response"
	"contactcenter-api/repository"
	"context"
	"fmt"
	"strconv"
)

type (
	IUserLog interface {
		PatchUserLog(ctx context.Context, domainUuid, userUuid, callId string, dispoSec string) (int, any)
	}
	UserLog struct{}
)

func NewUserLog() IUserLog {
	return &UserLog{}
}

func (s *UserLog) PatchUserLog(ctx context.Context, domainUuid, userUuid, callId string, dispoSec string) (int, any) {
	userLog, err := repository.UserLogRepo.GetUserLogByInfo(ctx, domainUuid, callId)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(userLog.UserLogUuid) > 0 {
		userLog.Comment = fmt.Sprintf("crm user_log call_id %s update disposec from %d to %s", callId, userLog.DispoSec, dispoSec)
		dispo, _ := strconv.ParseUint(dispoSec, 10, 64)
		userLog.DispoSec = int64(dispo)
		if err := repository.UserLogRepo.PatchUserLog(ctx, domainUuid, *userLog); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}

		return response.OK(map[string]any{
			"user_log_uuid": userLog.UserLogUuid,
		})
	}

	return response.OK(map[string]any{
		"message": fmt.Sprintf("user log uuid: %s not found", userLog.UserLogUuid),
	})
}
