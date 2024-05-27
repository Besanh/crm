package db

import (
	"contactcenter-api/common/model"
	sqlclient "contactcenter-api/internal/sqlclient"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type Domain struct {
}

func NewDomain() repository.IDomain {
	repo := &Domain{}
	repo.InitTable()
	repo.InitColumn()
	return repo
}

func (repo *Domain) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.DomainConfig)(nil)); err != nil {
		panic(err)
	}
}

func (repo *Domain) InitColumn() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.DomainConfig)(nil), "color TEXT NULL DEFAULT NULL"); err != nil {
		panic(err)
	}
	if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.DomainConfig)(nil), "cdr_status TEXT NULL DEFAULT NULL"); err != nil {
		panic(err)
	}
	if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.DomainConfig)(nil), "setting jsonb NULL DEFAULT NULL"); err != nil {
		panic(err)
	}
}

func (repo *Domain) GetDomainById(ctx context.Context, domainUuid string) (*model.Domain, error) {
	domain := new(model.Domain)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(domain).
		Where("domain_uuid = ?", domainUuid)
	err := query.Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else {
		return domain, nil
	}
}

func (repo *Domain) GetDomains(ctx context.Context, limit, offset int) (*[]model.Domain, int, error) {
	domains := new([]model.Domain)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(domains)
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	total, err := query.ScanAndCount(ctx)
	if err != nil && err != sql.ErrNoRows {
		return domains, 0, err
	} else {
		return domains, total, nil
	}
}

func (repo *Domain) InsertOrUpdateDomainConfig(ctx context.Context, domainConfig *model.DomainConfig) error {
	if repository.FusionSqlClient.GetDriver() == sqlclient.MYSQL {
		query := repository.FusionSqlClient.GetDB().NewInsert().Model(domainConfig).On("DUPLICATE KEY UPDATE")
		resp, err := query.Exec(ctx)
		if err != nil {
			return err
		} else if affected, _ := resp.RowsAffected(); affected != 1 {
			return errors.New("insert crm domain config failed")
		}
	} else {
		query := repository.FusionSqlClient.GetDB().NewInsert().Model(domainConfig).On("CONFLICT (domain_uuid) DO UPDATE")
		resp, err := query.Exec(ctx)
		if err != nil {
			return err
		} else if affected, _ := resp.RowsAffected(); affected != 1 {
			return errors.New("insert crm domain config failed")
		}
	}
	return nil
}

func (repo *Domain) GetDomainConfigById(ctx context.Context, domainUuid string) (*model.DomainConfigView, error) {
	domainConfig := new(model.DomainConfigView)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(domainConfig).
		Column("d.domain_name", "d.domain_uuid").
		Column("dc.api_url", "dc.version", "dc.logo", "dc.cdr_status", "dc.setting").
		Join("LEFT JOIN crm_domain_config dc on dc.domain_uuid = d.domain_uuid").
		Where("d.domain_uuid = ?", domainUuid).
		Limit(1)
	err := query.Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else {
		return domainConfig, nil
	}
}

func (repo *Domain) GetDomainConfigs(ctx context.Context, limit, offset int) (*[]model.DomainConfigView, int, error) {
	domainConfigs := new([]model.DomainConfigView)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		TableExpr("v_domains as d").
		Column("d.domain_name", "d.domain_uuid").
		Column("dc.api_url", "dc.version", "dc.logo", "dc.setting").
		Join("LEFT JOIN crm_domain_config dc on dc.domain_uuid = d.domain_uuid")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	total, err := query.ScanAndCount(ctx, domainConfigs)
	if err != nil && err != sql.ErrNoRows {
		return domainConfigs, 0, err
	} else {
		return domainConfigs, total, nil
	}
}

func (repo *Domain) PutDomainConfig(ctx context.Context, domainUuid string, domainConfig model.DomainConfig) error {
	query := repository.FusionSqlClient.GetDB().NewUpdate().
		Model(&domainConfig).
		Where("domain_uuid = ?", domainUuid)
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected != 1 {
		return errors.New("update domain config failed")
	}
	return nil
}

func (repo *Domain) PostDomain(ctx context.Context, domain model.Domain) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(&domain)
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected != 1 {
		return errors.New("insert domain failed")
	}

	return nil
}

func (repo *Domain) GetDomainByName(ctx context.Context, name string) (*model.Domain, error) {
	domain := new(model.Domain)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(domain).
		Where("domain_name = ?", name)
	err := query.Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else {
		return domain, nil
	}
}
