package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IUserLog interface {
	GetUserLogByInfo(ctx context.Context, domainUuid, callId string) (*model.UserLog, error)
	PatchUserLog(ctx context.Context, domainUuid string, userLog model.UserLog) error
}

var UserLogRepo IUserLog
