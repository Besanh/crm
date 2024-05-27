package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
)

type RingGroup struct {
}

func NewRingGroup() repository.IRingGroup {
	return &RingGroup{}
}

func (repo *RingGroup) GetRingGroupById(ctx context.Context, domainUuid, ringGroupUuid string) (*model.RingGroup, error) {
	ringGroup := new(model.RingGroup)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ringGroup).
		Where("domain_uuid = ?", domainUuid).
		Where("follow_me_uuid = ?", ringGroupUuid).
		Limit(1)
	err := query.Scan(ctx)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return ringGroup, err
}

func (repo *RingGroup) GetRingGroupByExtension(ctx context.Context, domainUuid, ringGroupExtension string) (*model.RingGroup, error) {
	ringGroup := new(model.RingGroup)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ringGroup).
		Where("domain_uuid = ?", domainUuid).
		Where("ring_group_extension = ?", ringGroupExtension).
		Limit(1)
	err := query.Scan(ctx)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return ringGroup, err
}

func (repo *RingGroup) GetRingGroupDestinationByRingGroupUuid(ctx context.Context, domainUuid, ringGroupUuid string) (*model.RingGroupDestination, error) {
	ringGroupDestination := new(model.RingGroupDestination)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ringGroupDestination).
		Where("domain_uuid = ?", domainUuid).
		Where("follow_me_uuid = ?", ringGroupUuid).
		Limit(1)
	err := query.Scan(ctx)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return ringGroupDestination, err
}

func (repo *RingGroup) GetRingGroupDestinationOfExtension(ctx context.Context, domainUuid, ringGroupUuid, extension string) (*model.RingGroupDestination, error) {
	ringGroupDestination := new(model.RingGroupDestination)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(ringGroupDestination).
		Where("domain_uuid = ?", domainUuid).
		Where("ring_group_uuid = ?", ringGroupUuid).
		Where("destination_number = ?", extension).
		Limit(1)
	err := query.Scan(ctx)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return ringGroupDestination, err
}
