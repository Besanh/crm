package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type Dialplan struct {
}

func NewDialplan() repository.IDialplan {
	return &Dialplan{}
}

func (repo *Dialplan) GetDialplanById(ctx context.Context, domainUuid, dialplanUuid string) (*model.Dialplan, error) {
	dialplan := new(model.Dialplan)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(dialplan).
		Where("domain_uuid = ?", domainUuid).
		Where("dialplan_uuid = ?", dialplanUuid).
		Limit(1)
	err := query.Scan(ctx)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return dialplan, err
}

func (repo *Dialplan) GetDialplanByNumber(ctx context.Context, domainUuid, dialplanNumber string) (*model.Dialplan, error) {
	dialplan := new(model.Dialplan)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(dialplan).
		Where("domain_uuid = ?", domainUuid).
		Where("dialplan_number = ?", dialplanNumber).
		Limit(1)
	err := query.Scan(ctx)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return dialplan, err
}

func (repo *Dialplan) UpdateDialplanTransaction(ctx context.Context, dialplan model.Dialplan, dialplanDetails []model.DialplanDetail) error {
	err := repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if isExisted, err := tx.NewSelect().Model((*model.Dialplan)(nil)).
			Where("dialplan_number = ?", dialplan.DialplanNumber).
			Where("domain_uuid = ?", dialplan.DomainUuid).
			Exists(ctx); err != nil {
			return err
		} else if !isExisted {
			if _, err = tx.NewInsert().Model(&dialplan).
				Exec(ctx); err != nil {
				return err
			}
		} else {
			if _, err := tx.NewUpdate().
				Model(&dialplan).
				Column("dialplan_name", "dialplan_xml").
				Where("dialplan_number = ?", dialplan.DialplanNumber).
				Where("domain_uuid = ?", dialplan.DomainUuid).
				Exec(ctx); err != nil {
				return err
			}
		}
		if _, err := repository.FusionSqlClient.GetDB().NewDelete().Model((*model.DialplanDetail)(nil)).
			Where("dialplan_uuid = ?", dialplan.DialplanUuid).
			Exec(ctx); err != nil {
			return err
		}
		if _, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&dialplanDetails).
			Exec(ctx); err != nil {
			return err
		}
		return nil
	})
	return err
}
