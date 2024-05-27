package common

import (
	"contactcenter-api/common/cache"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/util"
	"contactcenter-api/repository"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"golang.org/x/exp/slices"
)

func BuidTree(ctx context.Context, domainUuid string) (map[string]any, error) {
	tree := map[string]any{}
	filter := model.UnitFilter{
		IsParent: true,
		Level:    "0",
	}
	totalParent, unitParents, err := repository.UnitRepo.GetUnits(ctx, domainUuid, -1, 0, filter)
	if err != nil {
		log.Error(err)
		return map[string]any{}, err
	} else if totalParent > 0 {
		for _, val := range *unitParents {
			ProcessTree(ctx, domainUuid, tree, val, val.Level)
		}
	}

	return tree, nil
}

func BuildList(ctx context.Context, domainUuid string) (map[int]any, error) {
	tree := map[int]any{}
	filter := model.UnitFilter{
		IsParent: true,
		Level:    "0",
	}
	totalParent, unitParents, err := repository.UnitRepo.GetUnits(ctx, domainUuid, -1, 0, filter)
	if err != nil {
		log.Error(err)
		return map[int]any{}, err
	} else if totalParent > 0 {
		for _, val := range *unitParents {
			ProcessList(ctx, domainUuid, 0, tree, val, val.Level, val.UnitCode, 0)
		}
	}

	return tree, nil
}

func BuildListById(ctx context.Context, domainUuid, unitId string) (map[int]any, error) {
	tree := map[int]any{}
	unitParents, err := repository.UnitRepo.GetUnitById(ctx, domainUuid, unitId)
	if err != nil {
		log.Error(err)
		return map[int]any{}, err
	} else if unitParents != nil {
		ProcessList(ctx, domainUuid, 0, tree, *unitParents, unitParents.Level, unitParents.UnitCode, 0)
	}

	return tree, nil
}

func ProcessTree(ctx context.Context, domainUuid string, tree map[string]any, unit model.UnitInfo, level string) (map[string]any, error) {
	filter := model.UnitFilter{
		ParentUnitUuid: unit.UnitUuid,
	}
	if level == "0" {
		unitAppend := make([]any, 0)
		if tree[""] != nil {
			unitAppend = append(unitAppend, tree[""])
		} else {
			unitAppend = append(unitAppend, unit)
		}
		tree[""] = unitAppend
	} else {
		unitAppend := make([]any, 0)
		if tree[unit.ParentUnitUuid] != nil {
			for k, v := range tree {
				if k == unit.ParentUnitUuid {
					for _, val := range v.([]any) {
						val, ok := val.(model.UnitInfo)
						if ok {
							total, users, err := repository.UserCrmRepo.GetUserCrms(ctx, domainUuid, -1, 0, model.UserFilter{
								UnitUuid: val.UnitUuid,
							})
							if err != nil {
								log.Error(err)
								continue
							}
							if total > 0 {
								for _, u := range users {
									user := &model.User{}
									if err := util.ParseAnyToAny(u, user); err != nil {
										log.Error(err)
										continue
									}
									unit.Users = append(unit.Users, user)
								}
							}
						}
						unitAppend = append(unitAppend, val)
					}
				}
			}
		}
		if len(unit.UnitUuid) > 1 {
			unitAppend = append(unitAppend, unit)
		}

		tree[unit.ParentUnitUuid] = unitAppend
	}

	total, data, err := repository.UnitRepo.GetUnits(ctx, domainUuid, -1, 0, filter)
	if err != nil {
		log.Error(err)
		return map[string]any{}, err
	} else if total > 0 {
		for _, val := range *data {
			levelTmp, _ := strconv.Atoi(level)
			levelTmp++
			levelInt := strconv.Itoa(levelTmp)
			ProcessTree(ctx, domainUuid, tree, val, levelInt)
		}
	}

	return tree, nil
}

func ProcessList(ctx context.Context, domainUuid string, i int, tree map[int]any, unit any, level, parentUnitCode string, quantity int) map[int]any {
	tmp := unit.(model.UnitInfo)
	filter := model.UnitFilter{
		ParentUnitUuid: tmp.UnitUuid,
	}

	// If collapse smaller, open these comments
	// rootParentCode := findParentInTree(tmp.ParentUnitUuid, tree)
	// if len(rootParentCode) < 1 {
	// 	rootParentCode = tmp.UnitCode
	// }
	// tmp.ParentUnitCode = rootParentCode
	tmp.ParentUnitCode = parentUnitCode
	tmp.Quantity = quantity
	tree[len(tree)+1] = tmp

	total, data, err := repository.UnitRepo.GetUnits(ctx, domainUuid, -1, 0, filter)
	if err != nil {
		log.Error(err)
	} else if total > 0 {
		for _, val := range *data {
			levelTmp, _ := strconv.Atoi(level)
			levelTmp++
			levelInt := strconv.Itoa(levelTmp)
			filterQuantity := model.UnitFilter{
				ParentUnitUuid: val.UnitUuid,
			}
			totalQuantity, _, err := repository.UnitRepo.GetUnits(ctx, domainUuid, -1, 0, filterQuantity)
			if err != nil {
				log.Error(err)
				continue
			}
			ProcessList(ctx, domainUuid, i, tree, val, levelInt, parentUnitCode, totalQuantity)
		}
	}
	return tree
}

func CacheUserInUnit(ctx context.Context, domainUuid, unit string, user *model.UserView) ([]model.UserView, error) {
	listUsers, err := cache.RCache.Get(constants.PRIVILEGE_USER_UNIT + unit)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if listUsers != "" {
		var users []model.UserView
		if err := json.Unmarshal([]byte(listUsers), &users); err != nil {
			log.Error(err)
			return nil, err
		}
		return users, nil
	} else {
		tree, err := BuildListById(ctx, domainUuid, unit)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if len(tree) < 1 {
			return nil, nil
		}
		var listUnits []string
		for i := 1; i <= len(tree); i++ {
			listUnits = append(listUnits, tree[i].(model.UnitInfo).UnitUuid)
		}
		var usersInUnit []model.UserView
		var usersExcept []model.UserView
		if len(listUnits) > 0 {
			for _, v := range listUnits {
				_, userInUnit, err := repository.UserCrmRepo.GetUserCrms(ctx, domainUuid, -1, 0, model.UserFilter{
					UnitUuid: v,
				})
				if err != nil {
					log.Error(err)
					return nil, err
				}
				usersInUnit = append(usersInUnit, userInUnit...)
			}
		}

		var levelExcept []string
		if user.Level == constants.ROLE_MANAGER {
			levelExcept = []string{constants.ROLE_ADMIN}
		} else if user.Level == constants.ROLE_LEADER {
			levelExcept = []string{constants.ROLE_ADMIN, constants.ROLE_MANAGER}
		} else if user.Level == constants.ROLE_USER {
			levelExcept = []string{constants.ROLE_ADMIN, constants.ROLE_MANAGER, constants.ROLE_LEADER}
		}
		if len(levelExcept) > 0 {
			for _, v := range listUnits {
				_, userExcept, err := repository.UserCrmRepo.GetUserCrms(ctx, domainUuid, -1, 0, model.UserFilter{
					Levels:   levelExcept,
					UnitUuid: v,
				})
				if err != nil {
					log.Error(err)
					return nil, err
				}
				usersExcept = append(usersExcept, userExcept...)
			}
		}
		//remove user in usersExcept from usersInUnit
		if len(usersExcept) > 0 {
			for _, v := range usersExcept {
				for i, val := range usersInUnit {
					if val.UserUuid == v.UserUuid {
						usersInUnit = append(usersInUnit[:i], usersInUnit[i+1:]...)
					}
				}
			}
		}

		var userUuids []string
		for _, val := range usersInUnit {
			userUuids = append(userUuids, val.UserUuid)
		}
		if len(userUuids) > 0 {
			user.ManageCampaignUuids, err = repository.CampaignRepo.GetCampaignUuidsOfUsers(ctx, user.DomainUuid, userUuids)
			if err != nil {
				log.Error(err)
				return nil, err
			}
		}

		//set usersInUnit to cache
		if len(usersInUnit) > 0 {
			err = cache.RCache.SetTTL(constants.PRIVILEGE_USER_UNIT+unit, usersInUnit, 10*time.Minute)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			return usersInUnit, nil
		} else {
			return []model.UserView{}, nil
		}
	}
}

func GetParentTreeUnitUuid(ctx context.Context, domainUuid, unitUuid string) (map[int]any, error) {
	tree := map[int]any{}
	unit, err := repository.UnitRepo.GetUnitById(ctx, domainUuid, unitUuid)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	isLimit := false
	i := 0
	if unit.Level == "0" || len(unit.Level) < 1 {
		isLimit = true
		tree[0] = unit
		return tree, nil
	}
	for isLimit == false {
		if i > 200 || isLimit {
			isLimit = true
			return tree, nil
		}
		isLimit, tree = getParentUnit(ctx, domainUuid, i, unit.ParentUnitUuid, tree)
	}
	return tree, nil
}

func getParentUnit(ctx context.Context, domainUuid string, i int, parentUnitUuid string, tree map[int]any) (bool, map[int]any) {
	if len(parentUnitUuid) > 0 {
		parent, err := repository.UnitRepo.GetUnitById(ctx, domainUuid, parentUnitUuid)
		if err != nil {
			log.Error(err)
			return true, tree
		}
		if !checkItemExistInMap(tree, parent.UnitUuid) {
			if len(parent.UnitUuid) > 0 {
				tree[i] = parent
				i++
				getParentUnit(ctx, domainUuid, i, parent.ParentUnitUuid, tree)
			} else {
				return true, tree
			}
		}
	}

	return true, tree
}

func checkItemExistInMap(tree map[int]any, unitUuid string) bool {
	isExist := false
	for _, val := range tree {
		unit := &model.UnitInfo{}
		if err := util.ParseAnyToAny(val, unit); err != nil {
			log.Error(err)
			continue
		}
		if unitUuid == unit.UnitUuid {
			isExist = true
		}
	}
	return isExist
}

// formular: >=, <=, >, <
func GetUnitFormular(ctx context.Context, domainUuid, unitUuid, fromUnitLevel string, toUnitLevel string, encompass []string, formular string) (map[int]any, error) {
	tree := map[int]any{}
	unitParents, err := repository.UnitRepo.GetUnitById(ctx, domainUuid, unitUuid)
	if err != nil {
		log.Error(err)
		return tree, err
	} else if unitParents != nil {
		ProcessListReverse(ctx, domainUuid, 0, tree, *unitParents, unitParents.Level, unitParents.UnitCode, 0, fromUnitLevel, toUnitLevel, encompass)
	}

	return tree, nil
}

func ProcessListReverse(ctx context.Context, domainUuid string, i int, tree map[int]any, unit any, level, parentUnitCode string, quantity int, fromUnitLevel, toUnitLevel string, encompass []string) map[int]any {
	tmp := unit.(model.UnitInfo)
	filter := model.UnitFilter{
		UnitUuid: tmp.ParentUnitUuid,
	}

	tmp.ParentUnitCode = parentUnitCode
	tmp.Quantity = quantity
	if slices.Contains[[]string](encompass, tmp.Level) {
		tree[len(tree)+1] = tmp
	}

	filterEqual := model.UnitFilter{
		ParentUnitUuid: tmp.UnitUuid,
	}
	totalEqual, dataEqual, err := repository.UnitRepo.GetUnits(ctx, domainUuid, -1, 0, filterEqual)
	if err != nil {
		log.Error(err)
	} else if totalEqual > 0 {
		for _, val := range *dataEqual {
			if slices.Contains[[]string](encompass, val.Level) {
				if !checkUnitExistInTree(tree, val) {
					tree[len(tree)+1] = val
				}
			}
		}
	}

	toUnitLevelInt, _ := strconv.Atoi(toUnitLevel)
	fromUnitLevelInt, _ := strconv.Atoi(fromUnitLevel)

	total, data, err := repository.UnitRepo.GetUnits(ctx, domainUuid, -1, 0, filter)
	if err != nil {
		log.Error(err)
	} else if total > 0 {
		for _, val := range *data {
			levelTmp, _ := strconv.Atoi(level)
			levelTmp--
			if levelTmp < 0 {
				break
			}
			if levelTmp < fromUnitLevelInt || levelTmp > toUnitLevelInt {
				continue
			}
			levelInt := strconv.Itoa(levelTmp)
			filterQuantity := model.UnitFilter{
				UnitUuid: val.ParentUnitUuid,
			}
			totalQuantity, _, err := repository.UnitRepo.GetUnits(ctx, domainUuid, -1, 0, filterQuantity)
			if err != nil {
				log.Error(err)
				continue
			}
			ProcessListReverse(ctx, domainUuid, i, tree, val, levelInt, parentUnitCode, totalQuantity, fromUnitLevel, toUnitLevel, encompass)
		}
	}
	return tree
}

func checkUnitExistInTree(tree map[int]any, item model.UnitInfo) bool {
	for _, val := range tree {
		if val.(model.UnitInfo).UnitUuid == item.UnitUuid {
			return true
		}
	}
	return false
}
