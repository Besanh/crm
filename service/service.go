package service

import (
	"contactcenter-api/common/constants"
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/util"
	"contactcenter-api/repository"
	"contactcenter-api/service/common"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

type (
	PBXInformation struct {
		PBXDomain, APIDomain, PBXPort, PBXWss, PBXOutboundProxy string
		PBXTransport                                            string
	}
	ByTime []map[string]any
)

var (
	API_HOST = "http://localhost:8000"
	PBXInfo  PBXInformation
)

func AddLog(level string, message string, data any) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	v := model.TransactionLog{
		Level:     level,
		Message:   message,
		Id:        uuid.NewString(),
		CreatedAt: time.Now(),
		Meta:      fmt.Sprintf("%s:%d", srcFile, numLine),
	}
	if data != nil {
		typeData := reflect.ValueOf(data)
		if typeData.Kind() == reflect.String {
			v.Data = typeData.String()
		} else if typeData.Kind() == reflect.Map || typeData.Kind() == reflect.Struct {
			if bytes, err := json.Marshal(data); err != nil {
				log.Error(err)
			} else {
				v.Data = string(bytes)
			}
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := repository.TransactionRepo.InsertTransactionLog(ctx, v); err != nil {
		log.Error(err)
	}
}

func GetUserInfo(ctx context.Context, userUuid string) (*model.UserView, error) {
	user, err := repository.UserCrmRepo.GetUserCrmById(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return nil, err
	} else if user != nil {
		if util.InArr([]string{constants.SUPERADMIN, constants.ADMIN, constants.MANAGER, constants.LEADER}, user.Level) {
			userUuids := make([]string, 0)
			unitUuids := make([]*model.UnitInfo, 0)
			userUnits := make([]*model.UserView, 0)
			/**
			* 1. Lay tat ca user cua unit
			*  Neu level admin -> lấy từ admin -> manager -> leader -> user
			*  Neu level manager -> leader
			*  Neu level leader -> user
			* 2. Lay all user cua unit thap hon
			* Ap dung tuong tu cho extension
			 */
			tree := map[string]any{}
			if len(user.UnitUuid) > 0 {
				unit, err := repository.UnitRepo.GetUnitById(ctx, user.DomainUuid, user.UnitUuid)
				if err != nil {
					log.Error(err)
					return nil, err
				} else if unit == nil {
					return nil, errors.New("unit_uuid is not exists")
				}
				common.ProcessTree(ctx, user.DomainUuid, tree, *unit, unit.Level)
			}
			if len(user.UnitUuid) > 0 {
				for _, val := range tree {
					if val != nil {
						items := []model.UnitInfo{}
						if err := util.ParseAnyToAny(val, &items); err != nil {
							log.Error(err)
							continue
						}
						for _, item := range items {
							if len(item.Users) > 0 {
								for _, v := range item.Users {
									user := &model.UserView{}
									if err := util.ParseAnyToAny(v, user); err != nil {
										log.Error(err)
										continue
									} else {
										userUnits = append(userUnits, user)
										unit, err := repository.UnitRepo.GetUnitById(ctx, user.DomainUuid, v.UnitUuid)
										if err != nil {
											log.Error(err)
											continue
										}
										unitUuids = append(unitUuids, unit)
										// userUuids = append(userUuids, v.UserUuid)
									}
								}
							}
						}
					}
				}
			}
			user.ManageUsers = userUnits
			user.ManageUnits = unitUuids
			if len(userUuids) > 0 {
				if extensions, err := repository.ExtensionRepo.GetExtensionsOfUserUuids(ctx, user.DomainUuid, userUuids); err != nil {
					log.Error(err)
					return nil, errors.New("user info is invalid")
				} else if extensions != nil {
					user.ManageExtensions = *extensions
				}
			}

			return user, nil
		} else {
			return user, nil
		}
	} else {
		return nil, errors.New("user info is invalid")
	}
}

func HandleCollectionInfoUser(ctx context.Context, userUuid string, filter model.UserFilter) (model.UserFilter, error) {
	user, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		return filter, err
	}

	if user.Level == constants.MANAGER {
		filter.ManageExtensionUuids = make([]string, 0)
		for _, e := range user.ManageExtensions {
			filter.ManageExtensionUuids = append(filter.ManageExtensionUuids, e.ExtensionUuid)
		}
	} else if user.Level == constants.LEADER {
		filter.ManageExtensionUuids = make([]string, 0)
		for _, e := range user.ManageExtensions {
			filter.ManageExtensionUuids = append(filter.ManageExtensionUuids, e.ExtensionUuid)
		}
	}

	if len(user.ManageUsers) > 0 {
		filter.ManageUserUuids = make([]string, 0)
		for _, u := range user.ManageUsers {
			filter.ManageUserUuids = append(filter.ManageUserUuids, u.UserUuid)
		}
	}

	return filter, nil
}

func AddTransaction(domainUuid, userUuid, entity string, entityUuid string, action, status, result string, oldValue, newValue any) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	oldValueStr := ""
	newValueStr := ""
	if oldValue != nil {
		if bytes, err := json.Marshal(oldValue); err != nil {
			log.Error(err)
		} else {
			oldValueStr = string(bytes)
		}
	}
	if newValue != nil {
		if bytes, err := json.Marshal(newValue); err != nil {
			log.Error(err)
		} else {
			newValueStr = string(bytes)
		}
	}
	transaction := model.Transaction{
		DomainUuid:      domainUuid,
		Entity:          entity,
		TransactionUuid: uuid.NewString(),
		EntityUuid:      entityUuid,
		Action:          action,
		Status:          status,
		OldData:         oldValueStr,
		NewData:         newValueStr,
		Result:          result,
		CreatedAt:       time.Now(),
		UserUuid:        userUuid,
	}
	if err := repository.TransactionRepo.InsertTransaction(ctx, &transaction); err != nil {
		log.Error(err)
	}
}

func ParseListCustomName(listUuid string) string {
	listUuidConverted := strings.Replace(listUuid, "-", "_", 4)
	return fmt.Sprintf("custom_%s", listUuidConverted)
}

func ParseTimeUTCFormated(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func RecoveryApp() {
	if r := recover(); r != nil {
		log.Error(r)
	}
}

func (a ByTime) Len() int { return len(a) }
func (a ByTime) Less(i, j int) bool {
	iTimeStr := a[i]["time"].(string)
	jTimeStr := a[j]["time"].(string)
	iTime := util.ParseFromStringToTime(iTimeStr)
	jTime := util.ParseFromStringToTime(jTimeStr)
	return iTime.After(jTime)
}
func (a ByTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type Exports []*model.ExportMap

func (ex Exports) Len() int { return len(ex) }
func (ex Exports) Less(i, j int) bool {
	iTimeStr := ex[i].ExportTime
	jTimeStr := ex[j].ExportTime
	iTime := util.ParseFromStringToTime(iTimeStr)
	jTime := util.ParseFromStringToTime(jTimeStr)
	return iTime.After(jTime)
}
func (ex Exports) Swap(i, j int) { ex[i], ex[j] = ex[j], ex[i] }
