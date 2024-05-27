package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IDomain interface {
	GetDomains(ctx context.Context, limit, offset int) (*[]model.Domain, int, error)
	GetDomainById(ctx context.Context, domainUuid string) (domain *model.Domain, err error)
	GetDomainByName(ctx context.Context, domainName string) (domain *model.Domain, err error)
	GetDomainConfigById(ctx context.Context, domainUuid string) (*model.DomainConfigView, error)
	GetDomainConfigs(ctx context.Context, limit, offset int) (*[]model.DomainConfigView, int, error)
	InsertOrUpdateDomainConfig(ctx context.Context, domainConfig *model.DomainConfig) error
	PutDomainConfig(ctx context.Context, domainUuid string, domainConfig model.DomainConfig) error

	PostDomain(ctx context.Context, domain model.Domain) error
}

var DomainRepo IDomain
