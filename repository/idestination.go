package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IDestination interface {
	GetDestinationById(ctx context.Context, domainUuid, destinationUuid string) (*model.Destination, error)
	GetDestinationByNumber(ctx context.Context, domainUuid, destinationNumber string) (*model.Destination, error)
	InsertDestination(ctx context.Context, destination *model.Destination) error
}

var DestinationRepo IDestination
