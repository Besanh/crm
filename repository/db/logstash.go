package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"errors"
)

type LogstashRepo struct {
}

func NewLogstash() repository.ILogstash {
	return &LogstashRepo{}
}

func (repo *LogstashRepo) InsertLogstash(ctx context.Context, domainUuid string, logstash model.Logstash) error {
	res, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&logstash).
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("insert logstash fail")
	}

	return nil
}
