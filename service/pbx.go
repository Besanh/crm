package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type (
	IPbx interface {
		PostPbx(ctx context.Context, domainUuid, userUuid string, pbx model.Pbx) (int, any)
		GetPbxs(ctx context.Context, domainUuid, userUuid string, filter model.PbxFilter) (int, any)
		GetPbxByUnitId(ctx context.Context, domainUuid, userUuid, id string) (int, any)
		GetPbxById(ctx context.Context, domainUuid, userUuid, id string) (int, any)
		PutPbxById(ctx context.Context, domainUuid, userUuid, id string, pbxPut model.Pbx) (int, any)
		DeletePbxById(ctx context.Context, domainUuid, userUuid, id string) (int, any)
	}
	Pbx struct{}
)

func NewPbx() IPbx {
	return &Pbx{}
}

func (s *Pbx) PostPbx(ctx context.Context, domainUuid, userUuid string, pbx model.Pbx) (int, any) {
	if len(pbx.PbxUuid) > 0 {
		pbxExist, err := repository.PbxRepo.GetPbxByUnitId(ctx, domainUuid, pbx.PbxUuid)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if len(pbxExist.PbxUuid) > 0 {
			return response.ServiceUnavailableMsg("pbx already exist in unit")
		}
	}

	pbx.DomainUuid = domainUuid
	pbx.PbxUuid = uuid.NewString()
	pbx.CreatedBy = userUuid
	pbx.CreatedAt = time.Now()

	if err := repository.PbxRepo.InsertPbx(ctx, pbx); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Created(map[string]any{
		"id": pbx.PbxUuid,
	})
}

func (s *Pbx) GetPbxs(ctx context.Context, domainUuid, userUuid string, filter model.PbxFilter) (int, any) {
	user, err := repository.UserCrmRepo.GetUserCrmById(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(user.UserUuid) < 1 {
		return response.ServiceUnavailableMsg("user not found")
	}
	pbxs, err := repository.PbxRepo.GetPbxs(ctx, domainUuid, user.UnitUuid, filter)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(pbxs)
}

func (s *Pbx) GetPbxByUnitId(ctx context.Context, domainUuid, userUuid, id string) (int, any) {
	pbx, err := repository.PbxRepo.GetPbxByUnitId(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(pbx)
}

func (s *Pbx) GetPbxById(ctx context.Context, domainUuid, userUuid, id string) (int, any) {
	pbx, err := repository.PbxRepo.GetPbxById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(pbx.PbxUuid) < 1 {
		return response.ServiceUnavailableMsg("pbx not found")
	}

	return response.OK(pbx)
}

func (s *Pbx) PutPbxById(ctx context.Context, domainUuid, userUuid, id string, pbxPut model.Pbx) (int, any) {
	pbx, err := repository.PbxRepo.GetPbxById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(pbx.PbxUuid) < 1 {
		return response.ServiceUnavailableMsg("pbx not found")
	}

	pbxUnitExist, err := repository.PbxRepo.GetPbxByUnitId(ctx, domainUuid, pbxPut.UnitUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(pbxUnitExist.PbxUuid) > 0 {
		return response.ServiceUnavailableMsg("pbx is exist in unit")
	}

	pbx.PbxName = pbxPut.PbxName
	pbx.Status = pbxPut.Status
	pbx.Domain = pbxPut.Domain
	pbx.OutboundProxy = pbxPut.OutboundProxy
	pbx.Wss = pbxPut.Wss
	pbx.Transport = pbxPut.Transport
	pbx.Port = pbxPut.Port
	pbx.UnitUuid = pbxPut.UnitUuid
	pbx.UpdatedBy = userUuid
	pbx.UpdatedAt = time.Now()

	if err := repository.PbxRepo.PutPbxById(ctx, domainUuid, pbx); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"pbx_uuid": id,
	})
}

func (s *Pbx) DeletePbxById(ctx context.Context, userUuid, domainUuid, id string) (int, any) {
	pbx, err := repository.PbxRepo.GetPbxById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(pbx.PbxUuid) < 1 {
		return response.ServiceUnavailableMsg("pbx not found")
	}

	if err := repository.PbxRepo.DeletePbxById(ctx, domainUuid, id); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"pbx_uuid": id,
	})
}
