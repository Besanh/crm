package db

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type RoleGroupRepo struct {
}

func NewRoleGroup() repository.IRoleGroup {
	repo := &RoleGroupRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *RoleGroupRepo) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.RoleGroup)(nil)); err != nil {
		panic(err)
	}
}

func (repo *RoleGroupRepo) InitColumn() {

}

func (repo *RoleGroupRepo) InitIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.RoleGroup)(nil)).IfNotExists().Index("idx_role_group_domain_uuid").Column("domain_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.RoleGroup)(nil)).IfNotExists().Index("idx_role_group_uuid").Column("role_group_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.RoleGroup)(nil)).IfNotExists().Index("idx_role_group_name").Column("role_group_name").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.RoleGroup)(nil)).IfNotExists().Index("idx_role_group_status").Column("status").Exec(ctx); err != nil {
		log.Error(err)
	}
}

func (repo *RoleGroupRepo) GetRoleGroup(ctx context.Context, domainUuid string, limit, offset int, filter model.RoleGroupFilter) (int, *[]model.RoleGroup, error) {
	result := new([]model.RoleGroup)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(result).
		Where("domain_uuid = ?", domainUuid)
	if len(filter.RoleGroupName) > 0 {
		query.Where("role_group_name = ?", filter.RoleGroupName)
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status.Bool)
	}
	if filter.StartTime.Valid {
		query.Where("start_time = ?", filter.StartTime.Time)
	}
	if filter.EndTime.Valid {
		query.Where("end_time = ?", filter.EndTime.Time)
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	query.Order("created_at DESC")

	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return 0, nil, err
	} else if err == sql.ErrNoRows {
		return 0, nil, nil
	}

	return total, result, nil
}

func (repo *RoleGroupRepo) InsertRoleGroup(ctx context.Context, domainUuid string, userUuid string, roleGroup model.RoleGroup) error {
	resp, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&roleGroup).Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("insert role group failed")
	}

	return nil
}

func (repo *RoleGroupRepo) GetRoleGroupById(ctx context.Context, domainUuid, id string) (*model.RoleGroup, error) {
	result := new(model.RoleGroup)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(result).
		Where("domain_uuid = ?", domainUuid).
		Where("role_group_uuid = ?", id)
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, errors.New("role group not found")
	}

	return result, nil
}

func (repo *RoleGroupRepo) UpdateRoleGroupById(ctx context.Context, domainUuid string, roleGroup model.RoleGroup) error {
	resp, err := repository.FusionSqlClient.GetDB().NewUpdate().Model(&roleGroup).WherePK().Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("update role group failed")
	}

	return nil
}

func (repo *RoleGroupRepo) DeleteRoleGroupById(ctx context.Context, domainUuid, id string) error {
	res, err := repository.FusionSqlClient.GetDB().NewDelete().Model(&model.RoleGroup{}).
		Where("domain_uuid = ?", domainUuid).
		Where("role_group_uuid = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("delete role group failed")
	}

	return nil
}
