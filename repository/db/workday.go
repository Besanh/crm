package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/uptrace/bun"
)

type WorkDay struct {
}

func NewWorkDay() repository.IWorkDay {
	repo := &WorkDay{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return &WorkDay{}
}

func (repo *WorkDay) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := repository.FusionSqlClient.GetDB().Query(
		`DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'workday_type') THEN
				CREATE TYPE workday_type AS ENUM ('work', 'holiday');
			END IF;
		END
		$$
		`,
		nil,
	); err != nil {
		panic(err)
	}
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.WorkDay)(nil)); err != nil {
		panic(err)
	}
}

func (repo *WorkDay) InitColumn() {

}

func (repo *WorkDay) InitIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.WorkDay)(nil)).IfNotExists().Index("idx_workday_domain_uuid").Column("domain_uuid").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.WorkDay)(nil)).IfNotExists().Index("idx_workday_uuid").Column("workday_uuid").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.WorkDay)(nil)).IfNotExists().Index("idx_workday_day").Column("day").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.WorkDay)(nil)).IfNotExists().Index("idx_workday_start_time").Column("start_time").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.WorkDay)(nil)).IfNotExists().Index("idx_workday_end_time").Column("end_time").Exec(ctx); err != nil {
		panic(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.WorkDay)(nil)).IfNotExists().Index("idx_workday_is_work").Column("is_work").Exec(ctx); err != nil {
		panic(err)
	}
}

func (repo *WorkDay) InsertWorkDays(ctx context.Context, workDays *[]model.WorkDay) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, workDay := range *workDays {
			tmp := workDay
			existed, err := tx.NewSelect().Model(&workDay).
				Where("day = ?", workDay.Day).
				Where("domain_uuid = ?", workDay.DomainUuid).
				Where("unit_uuid = ?", workDay.UnitUuid).
				Exists(ctx)
			if err != nil {
				return err
			} else if !existed {
				workDay = tmp
				_, err = tx.NewInsert().Model(&workDay).
					Where("day = ?", workDay.Day).
					Where("domain_uuid = ?", workDay.DomainUuid).
					Exec(ctx)
				if err != nil {
					return err
				}
			} else {
				workDay = tmp
				workDay.UpdatedAt = time.Now()
				_, err = tx.NewUpdate().Model(&workDay).
					Where("day = ?", workDay.Day).
					Where("domain_uuid = ?", workDay.DomainUuid).
					Column("workday_name", "status", "start_time", "end_time", "description", "is_work").
					Exec(ctx)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (repo *WorkDay) GetWorkDays(ctx context.Context, domainUuid string, filter model.WorkDayFilter) (int, *[]model.WorkDay, error) {
	workDays := new([]model.WorkDay)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(workDays).
		ColumnExpr("DISTINCT ON (wd.workday_id) workday_id").
		Column("wd.workday_name", "wd.status", "wd.description").
		Relation("Unit").
		Where("wd.domain_uuid = ?", domainUuid)
	if len(filter.WorkDayName) > 0 {
		query.Where("workday_name = ?", filter.WorkDayName)
	}
	if len(filter.UnitUuids) > 0 {
		query.Where("wd.unit_uuid IN (?)", bun.In(filter.UnitUuids))
	}

	if filter.Limit > 0 {
		query.Limit(filter.Limit).Offset(filter.Offset)
	}
	query.Order("wd.workday_id DESC")
	if total, err := query.ScanAndCount(ctx); err != nil && err != sql.ErrNoRows {
		return 0, workDays, err
	} else {
		return total, workDays, nil
	}
}

func (repo *WorkDay) GetWorkDayById(ctx context.Context, domainUuid, id string) (*model.WorkDay, error) {
	workDay := new(model.WorkDay)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(workDay).
		Relation("Unit").
		Where("domain_uuid = ?", domainUuid).
		Where("workday_uuid = ?", id).
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return workDay, nil
	}
}

func (repo *WorkDay) DeleteWorkDay(ctx context.Context, domainUuid, id string) error {
	query := repository.FusionSqlClient.GetDB().
		NewDelete().
		Model((*model.WorkDay)(nil)).
		Where("domain_uuid = ?", domainUuid).
		Where("workday_uuid = ?", id)
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 0 {
		return errors.New("delete workday fail")
	}
	return nil
}

func (repo *WorkDay) DeleteWorkdayByWorkdayId(ctx context.Context, domainUuid, workdayId string) error {
	query := repository.FusionSqlClient.GetDB().
		NewDelete().
		Model((*model.WorkDay)(nil)).
		Where("domain_uuid = ?", domainUuid).
		Where("workday_id = ?", workdayId)
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 0 {
		return errors.New("delete workday fail")
	}
	return nil
}

func (repo *WorkDay) GetWorkDayByWorkdayId(ctx context.Context, domainUuid, workdayId string) (*[]model.WorkDay, error) {
	workDay := new([]model.WorkDay)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(workDay).
		Relation("Unit").
		Where("wd.domain_uuid = ?", domainUuid).
		Where("workday_id = ?", workdayId)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return workDay, nil
}
