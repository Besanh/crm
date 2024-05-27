package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
)

type UserLogRepo struct{}

func NewUserLog() repository.IUserLog {
	return &UserLogRepo{}
}

func (repo *UserLogRepo) PatchUserLog(ctx context.Context, domainUuid string, userLog model.UserLog) error {
	query := repository.FusionSqlClient.GetDB().NewUpdate().Model(&userLog).
		Where("domain_uuid = ?", domainUuid).
		WherePK()
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("update user log failed")
	}
	return nil
}

func (repo *UserLogRepo) GetUserLogByInfo(ctx context.Context, domainUuid, callId string) (*model.UserLog, error) {
	result := new(model.UserLog)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Where("domain_uuid = ?", domainUuid).
		Where("sub_status = ?", "HANGUP").
		Where("call_uuid = ?", callId)
	err := query.Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else {
		return result, nil
	}
}
