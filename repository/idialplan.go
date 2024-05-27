package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IDialplan interface {
	GetDialplanById(ctx context.Context, domainUuid, dialplanUuid string) (*model.Dialplan, error)
	GetDialplanByNumber(ctx context.Context, domainUuid, dialplanNumber string) (*model.Dialplan, error)
	UpdateDialplanTransaction(ctx context.Context, dialplan model.Dialplan, dialplanDetails []model.DialplanDetail) error
}

var DialplanRepo IDialplan
