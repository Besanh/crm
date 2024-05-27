package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
)

type ClassifyGroup struct{}

func NewClassifyGroupRepo() repository.IClassifyGroupRepo {
	return &ClassifyGroup{}
}

func (repo *ClassifyGroup) InsertClassifyGroup(ctx context.Context, classifyGroup model.ClassifyGroup) error {
	resp, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&classifyGroup).Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("insert classify group failed")
	}

	return nil
}

func (repo *ClassifyGroup) GetClassifyGroups(ctx context.Context, domainUuid string, limit, offset int) (int, []model.ClassifyGroup, error) {
	result := new([]model.ClassifyGroup)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(result).
		Distinct().
		Where("domain_uuid = ?", domainUuid).
		Where("status = true")
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	query.Order("created_at DESC")
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return 0, *result, err
	} else if err == sql.ErrNoRows {
		return 0, *result, errors.New("classify group not found")
	}
	return total, *result, nil
}

func (repo *ClassifyGroup) GetClassifyGroupById(ctx context.Context, domainUuid, id string) (*model.ClassifyGroup, error) {
	result := new(model.ClassifyGroup)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(&result).
		Where("domain_uuid = ?", domainUuid).
		Where("classify_group_uuid = ?", id)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *ClassifyGroup) PutClassifyGroupById(ctx context.Context, domainUuid string, classifyGroup model.ClassifyGroup) error {
	resp, err := repository.FusionSqlClient.GetDB().NewUpdate().Model(&classifyGroup).
		Where("domain_uuid = ?", domainUuid).
		Where("classify_group_uuid = ?", classifyGroup.ClassifyGroupUuid).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("update classify group failed")
	}

	return nil
}
