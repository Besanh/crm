package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
)

type Destination struct {
}

func NewDestination() repository.IDestination {
	return &Destination{}
}

func (repo *Destination) GetDestinationById(ctx context.Context, domainUuid, destinationUuid string) (*model.Destination, error) {
	destination := new(model.Destination)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(destination).
		Where("domain_uuid = ?", domainUuid).
		Where("destination_uuid = ?", destinationUuid).
		Limit(1)
	err := query.Scan(ctx)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return destination, err
}

func (repo *Destination) GetDestinationByNumber(ctx context.Context, domainUuid, destinationNumber string) (*model.Destination, error) {
	destination := new(model.Destination)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(destination).
		Where("domain_uuid = ?", domainUuid).
		Where("destination_number = ?", destinationNumber).
		Limit(1)
	err := query.Scan(ctx)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return destination, err
}

func (repo *Destination) InsertDestination(ctx context.Context, destination *model.Destination) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(destination)
	_, err := query.Exec(ctx)
	return err
}
