package common

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"time"

	"github.com/google/uuid"
)

// entity, entityUuid, action, status, result string, oldValue, newValue any
func AddLogstash(domainUuid, userUuid, index, storageDialect string, logstash model.Logstash) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	oldDataStr := ""
	newDataStr := ""
	if logstash.OldData != nil {
		if bytes, err := json.Marshal(logstash.OldData); err != nil {
			log.Error(err)
		} else {
			oldDataStr = string(bytes)
		}
	}
	if logstash.NewData != nil {
		if bytes, err := json.Marshal(logstash.NewData); err != nil {
			log.Error(err)
		} else {
			newDataStr = string(bytes)
		}
	}
	logstash.OldData = oldDataStr
	logstash.NewData = newDataStr
	logstash.CreatedAt = time.Now()
	logstash.CreatedBy = userUuid

	if storageDialect == "SQL" {
		if err := repository.LogstashRepo.InsertLogstash(ctx, domainUuid, logstash); err != nil {
			return err
		}
	} else if storageDialect == "ES" {
		esDoc := map[string]any{}
		tmpByte, err := json.Marshal(logstash)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(tmpByte, &esDoc); err != nil {
			return err
		}
		if err := repository.ESRepo.InsertLog(ctx, domainUuid, index, logstash.LogstashUuid, esDoc); err != nil {
			return err
		}
	}

	return nil
}

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
