package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/repository"
	"context"

	"github.com/google/uuid"
)

type (
	IDomain interface {
		PostDomain(ctx context.Context, domain model.Domain) (int, any)
		GetDomainConfigById(ctx context.Context, domainUuid string) (int, any)
		GetDomainConfigs(ctx context.Context, limit, offset int) (int, any)
		PutDomainConfig(ctx context.Context, domainUuid, userUuid string, domainConfig model.DomainConfig) (int, any)
		PostDomainConfig(ctx context.Context, domainUuid string, config model.DomainConfigPut) (int, any)
		PutDomainConfigById(ctx context.Context, domainUuid string, domainConfig model.DomainConfigPut) (int, any)
	}
	Domain struct {
	}
)

func NewDomain() IDomain {
	return &Domain{}
}

func (s *Domain) GetDomainConfigById(ctx context.Context, domainUuid string) (int, any) {
	domain, err := repository.DomainRepo.GetDomainConfigById(ctx, domainUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if domain == nil {
		return response.NotFoundMsg("domain not found")
	}
	return response.OK(domain)
}

func (s *Domain) GetDomainConfigs(ctx context.Context, limit, offset int) (int, any) {
	data, total, err := repository.DomainRepo.GetDomainConfigs(ctx, limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.Pagination(data, total, limit, offset)
}

func (s *Domain) PutDomainConfig(ctx context.Context, domainUuid, userUuid string, domainConfig model.DomainConfig) (int, any) {
	domainConfigExist, err := repository.DomainRepo.GetDomainConfigById(ctx, domainUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	domainConfigExist.CdrStatus = domainConfig.CdrStatus
	domainConfigUpdate := model.DomainConfig{
		DomainUuid:     domainUuid,
		Version:        domainConfig.Version,
		Logo:           domainConfig.Logo,
		APIUrl:         domainConfig.APIUrl,
		Color:          domainConfig.Color,
		CdrStatus:      domainConfig.CdrStatus,
		MissCallStatus: domainConfigExist.MissCallStatus,
		Setting:        domainConfig.Setting,
	}
	if err := repository.DomainRepo.PutDomainConfig(ctx, domainUuid, domainConfigUpdate); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"id": domainUuid,
	})
}

func (s *Domain) PostDomainConfig(ctx context.Context, domainUuid string, config model.DomainConfigPut) (int, any) {
	if domain, err := repository.DomainRepo.GetDomainById(ctx, domainUuid); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if domain == nil {
		return response.NotFoundMsg("domain is not exist")
	}
	cfg := model.DomainConfig{
		DomainUuid:     domainUuid,
		Version:        config.Version,
		Logo:           config.Logo,
		APIUrl:         config.APIUrl,
		Color:          config.Color,
		MissCallStatus: config.MissCallStatus,
		Setting:        config.Setting,
	}
	if err := repository.DomainRepo.InsertOrUpdateDomainConfig(ctx, &cfg); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(config)
}

func (s *Domain) PutDomainConfigById(ctx context.Context, domainUuid string, domainConfig model.DomainConfigPut) (int, any) {
	domain, err := repository.DomainRepo.GetDomainConfigById(ctx, domainUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if domain == nil {
		return response.NotFoundMsg("domain is not exist")
	}

	domainConfigUpdate := model.DomainConfig{
		DomainUuid: domainUuid,
		Version:    domainConfig.Version,
		Logo:       domainConfig.Logo,
		APIUrl:     domainConfig.APIUrl,
		Color:      domainConfig.Color,
		// CdrStatus:      domainConfig.CdrStatus,
		MissCallStatus: domainConfig.MissCallStatus,
		Setting:        domainConfig.Setting,
	}
	if err := repository.DomainRepo.PutDomainConfig(ctx, domainUuid, domainConfigUpdate); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"id": domainUuid,
	})
}

func (s *Domain) PostDomain(ctx context.Context, domain model.Domain) (int, any) {
	domainExist, err := repository.DomainRepo.GetDomainByName(ctx, domain.DomainName)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(domainExist.DomainUuid) > 0 {
		return response.ServiceUnavailableMsg("domain already exists")
	}

	domain.DomainUuid = uuid.New().String()
	if err := repository.DomainRepo.PostDomain(ctx, domain); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Created(map[string]any{
		"id": domain.DomainUuid,
	})
}
