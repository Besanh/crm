package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	"contactcenter-api/repository"
	"contactcenter-api/service/common"
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type (
	IUnit interface {
		GetUnits(ctx context.Context, domainUuid, userUuid string, limit, offset int, filter model.UnitFilter) (int, any)
		GetUnitsRecursive(ctx context.Context, domainUuid, viewType string) (int, any)
		GetUnitTreeRender(ctx context.Context, domainUuid, userUuid string) (int, any)
		PostUnit(ctx context.Context, domainUuid, userUuid string, unit model.Unit) (int, any)
		GetUnitById(ctx context.Context, domainUuid, userUuid, id string) (int, any)
		PutUnitById(ctx context.Context, domainUuid, userUuid, id string, unit model.Unit) (int, any)
		DeleteUnitById(ctx context.Context, domainUuid, userUuid, id string) (int, any)
		PatchUnitLogo(ctx context.Context, domainUuid, userUuid, UnitUuid, data string) (int, any)
		GetUnitLogo(ctx context.Context, fileName string) (string, error)
		DeleteUnitLogo(ctx context.Context, domainUuid, userUuid, id string) (int, any)
		ExportUnits(ctx context.Context, domainUuid, userUuid string, filter model.UnitFilter) (string, error)
		GetUnitParentTree(ctx context.Context, domainUuid, userUuid string, filter model.ParentUnitFilter) (int, any)
		GetUnitChildTree(ctx context.Context, domainUuid, userUuid string, filter model.UnitFilter) (int, any)
		GetTreeExcludeCurrentUnit(ctx context.Context, domainUuid, userUuid string, filter model.UnitFilter) (int, any)
	}
	Unit struct{}
)

func NewUnit() IUnit {
	return &Unit{}
}

func (u *Unit) GetUnits(ctx context.Context, domainUuid, userUuid string, limit, offset int, filter model.UnitFilter) (int, any) {
	total, units, err := repository.UnitRepo.GetUnits(ctx, domainUuid, limit, offset, filter)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Pagination(units, total, limit, offset)
}

func (u *Unit) GetUnitsRecursive(ctx context.Context, domainUuid, viewType string) (int, any) {
	if viewType == "tree" {
		tree, err := common.BuidTree(ctx, domainUuid)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		return response.Data(http.StatusOK, tree)
	} else if viewType == "list" {
		list, err := common.BuildList(ctx, domainUuid)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		return response.Data(http.StatusOK, list)
	}

	return response.BadRequestMsg("invalid view type")
}

func (u *Unit) GetUnitTreeRender(ctx context.Context, domainUuid, userUuid string) (int, any) {
	userCrm, err := repository.UserCrmRepo.GetUserCrmById(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(userCrm.UserUuid) < 1 {
		return response.BadRequestMsg("user is not exists")
	}
	tree := map[int]any{}
	if len(userCrm.UnitUuid) > 0 {
		tree, err = common.BuildListById(ctx, domainUuid, userCrm.UnitUuid)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
	}

	return response.Data(http.StatusOK, tree)
}

func (u *Unit) PostUnit(ctx context.Context, domainUuid, userUuid string, unit model.Unit) (int, any) {
	if len(unit.ParentUnitUuid) > 0 {
		parentUnit, err := repository.UnitRepo.GetUnitById(ctx, domainUuid, unit.ParentUnitUuid)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if parentUnit == nil {
			return response.NotFoundMsg("parent unit is not found")
		}
		levelTmp, _ := strconv.Atoi(parentUnit.Level)
		levelTmp++
		unit.Level = strconv.Itoa(levelTmp)
	} else {
		unit.Level = "0"
	}

	unit.UnitUuid = uuid.NewString()
	unit.DomainUuid = domainUuid
	unit.CreatedAt = time.Now()
	unit.CreatedBy = userUuid

	if err := repository.UnitRepo.InsertUnit(ctx, unit); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Created(map[string]any{
		"unit_uuid": unit.UnitUuid,
	})
}

func (u *Unit) GetUnitById(ctx context.Context, domainUuid, userUuid, id string) (int, any) {
	unit, err := repository.UnitRepo.GetUnitById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if unit == nil {
		return response.NotFoundMsg("unit is not found")
	}

	return response.Data(http.StatusOK, unit)
}

func (u *Unit) PutUnitById(ctx context.Context, domainUuid, userUuid, id string, unit model.Unit) (int, any) {
	unitExist, err := repository.UnitRepo.GetUnitById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if unitExist == nil {
		return response.NotFoundMsg("unit is not found")
	}

	if err := util.ParseAnyToAny(unit, unitExist); err != nil {
		log.Error(err)
		return response.BadRequestMsg(err.Error())
	}
	unit.UpdatedBy = userUuid
	unitExist.UpdatedAt = time.Now()

	if len(unit.ParentUnitUuid) > 0 {
		parentUnit, err := repository.UnitRepo.GetUnitById(ctx, domainUuid, unit.ParentUnitUuid)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if parentUnit == nil {
			return response.NotFoundMsg("parent unit is not found")
		}
		levelTmp, _ := strconv.Atoi(parentUnit.Level)
		levelTmp++
		unitExist.Level = strconv.Itoa(levelTmp)
	}

	if err := repository.UnitRepo.PutUnit(ctx, domainUuid, *unitExist); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"unit_uuid": id,
	})
}

func (u *Unit) DeleteUnitById(ctx context.Context, domainUuid, userUuid, id string) (int, any) {
	unit, err := repository.UnitRepo.GetUnitById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if unit == nil {
		return response.NotFoundMsg("unit is not found")
	}

	if err := repository.UnitRepo.DeleteUnitById(ctx, domainUuid, id); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"unit_uuid": id,
	})
}

func (u *Unit) PatchUnitLogo(ctx context.Context, domainUuid, userUuid, UnitUuid, data string) (int, any) {
	_, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	unit, err := repository.UnitRepo.GetUnitById(ctx, domainUuid, UnitUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if unit == nil {
		return response.BadRequestMsg("contact is not existed")
	}
	dir := util.UPLOAD_DIR + "unit/"
	fileName := "logo_" + unit.UnitUuid
	if fileName, err = util.DecodeAndSaveImageBase64(data, dir, fileName); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err)
	}
	unit.UnitConfig.Logo = fileName
	unitInfo := &model.UnitInfo{}
	if err := util.ParseAnyToAny(unit, unitInfo); err != nil {
		log.Error(err)
		return response.BadRequestMsg(err.Error())
	}

	if err := repository.UnitRepo.PutUnit(ctx, domainUuid, *unitInfo); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(map[string]any{
		"unit_uuid": unit.UnitUuid,
	})
}

func (s *Unit) GetUnitLogo(ctx context.Context, fileName string) (string, error) {
	logo := ""
	if _, err := os.Stat(util.UPLOAD_DIR + "unit/" + fileName); err != nil {
		return "", errors.New("file " + fileName + " is not exist")

	}
	logo = util.UPLOAD_DIR + "unit/" + fileName
	return logo, nil
}

func (s *Unit) DeleteUnitLogo(ctx context.Context, domainUuid, userUuid, id string) (int, any) {
	unit, err := repository.UnitRepo.GetUnitById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if unit == nil {
		return response.ServiceUnavailableMsg("unit is not exist")
	}
	logo := unit.UnitConfig.Logo
	if _, err := os.Stat(util.UPLOAD_DIR + "unit/" + logo); err != nil {
		return response.ServiceUnavailableMsg(errors.New("file " + logo + " is not exist"))

	}
	logo = util.UPLOAD_DIR + "unit/" + logo
	if err := os.Remove(logo); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	unit.UnitConfig.Logo = ""
	unit.UpdatedAt = time.Now()
	unit.UpdatedBy = userUuid
	if err := repository.UnitRepo.PutUnit(ctx, domainUuid, *unit); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"unit_uuid": id,
	})
}

func (s *Unit) GetUnitParentTree(ctx context.Context, domainUuid, userUuid string, filter model.ParentUnitFilter) (int, any) {
	user, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if len(user.UnitUuid) < 1 {
		return response.ServiceUnavailableMsg("unit not found")
	}
	result, err := common.GetParentTreeUnitUuid(ctx, domainUuid, user.UnitUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if len(filter.Level) > 0 {
		data := map[int]any{}
		i := 0
		for _, val := range result {
			item := model.Unit{}
			if err := util.ParseAnyToAny(val, &item); err != nil {
				log.Error(err)
				break
			}
			if item.Level == filter.Level {
				data[i] = item
				i++
			}
		}
		return response.OK(data)
	} else {
		return response.OK(result)
	}
}

func (s *Unit) GetUnitChildTree(ctx context.Context, domainUuid, userUuid string, filter model.UnitFilter) (int, any) {
	var tree map[int]any
	user, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if len(user.UnitUuid) < 1 {
		return response.ServiceUnavailableMsg("unit not found")
	}
	tree, err = common.GetUnitFormular(ctx, domainUuid, user.UnitUuid, filter.FromUnitLevel, filter.ToUnitLevel, filter.Encompass, filter.Formular)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(tree)
}

func (s *Unit) GetTreeExcludeCurrentUnit(ctx context.Context, domainUuid, userUuid string, filter model.UnitFilter) (int, any) {
	userInfo, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	tree := []any{}
	unitLv1 := []model.UnitInfo{}
	treeTmp, err := common.GetUnitFormular(ctx, domainUuid, userInfo.UnitUuid, "1", "4", []string{"1"}, "")
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if len(treeTmp) > 0 {
		item := model.Unit{}
		if err := util.ParseAnyToAny(treeTmp[len(treeTmp)], &item); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		filterUnit := model.UnitFilter{
			Level: "1",
		}
		_, unitLevel1, err := repository.UnitRepo.GetUnits(ctx, domainUuid, -1, 0, filterUnit)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}

		if len(*unitLevel1) > 0 {
			for _, val := range *unitLevel1 {
				if val.UnitUuid != item.UnitUuid {
					unitLv1 = append(unitLv1, val)
				}
			}
		}
	}

	if len(unitLv1) > 0 {
		for _, item := range unitLv1 {
			unit := model.Unit{}
			if err := util.ParseAnyToAny(item, &unit); err != nil {
				log.Error(err)
				continue
			}
			unitChild, err := common.BuildListById(ctx, domainUuid, unit.UnitUuid)
			if err != nil {
				log.Error(err)
				continue
			}
			if len(unitChild) > 0 {
				for _, item := range unitChild {
					tree = append(tree, item)
				}
			}
		}
	}

	return response.OK(tree)
}

func CheckUnitExistInTree(unitUuid string, tree map[int]any) (isExist bool) {
	for _, item := range tree {
		units := []model.Unit{}
		if err := util.ParseAnyToAny(item, &units); err != nil {
			log.Error(err)
			continue
		}
		if len(units) > 0 {
			for _, val := range units {
				if unitUuid == val.UnitUuid {
					isExist = true
					return
				}
			}
		}
	}
	return
}
