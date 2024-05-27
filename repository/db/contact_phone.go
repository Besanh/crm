package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
)

func (repo *Contact) InsertContactPhone(ctx context.Context, contact *model.ContactPhone) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(contact)
	_, err := query.Exec(ctx)
	return err
}
