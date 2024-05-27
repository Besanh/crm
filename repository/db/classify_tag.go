package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
)

type ClassifyTag struct{}

func NewClassifyTagRepo() repository.IClassifyTagRepo {
	return &ClassifyTag{}
}

func (repo *ClassifyTag) InsertClassifyTag(ctx context.Context, classifyTag model.ClassifyTag) error {
	resp, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&classifyTag).Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("insert classify tag failed")
	}

	return nil
}

func (repo *ClassifyTag) GetClassifyTags(ctx context.Context, domainUuid string, limit, offset int) (int, []model.ClassifyTag, error) {
	var result []model.ClassifyTag
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(&result).
		Distinct().
		Where("domain_uuid = ?", domainUuid)
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	query.Order("created_at DESC")
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return 0, result, err
	} else if err == sql.ErrNoRows {
		return 0, result, errors.New("classify tag not found")
	}
	return total, result, nil
}

func (repo *ClassifyTag) GetClassifyTagById(ctx context.Context, domainUuid, id string) (*model.ClassifyTag, error) {
	result := new(model.ClassifyTag)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(&result).
		Where("domain_uuid = ?", domainUuid).
		Where("classify_tag_uuid = ?", id)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *ClassifyTag) PutClassifyTagById(ctx context.Context, domainUuid string, classifyTag model.ClassifyTag) error {
	resp, err := repository.FusionSqlClient.GetDB().NewUpdate().Model(&classifyTag).
		Where("domain_uuid = ?", domainUuid).
		Where("classify_tag_uuid = ?", classifyTag.ClassifyTagUuid).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("update classify tag failed")
	}

	return nil
}
