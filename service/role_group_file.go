package service

import (
	"contactcenter-api/common/constants"
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/util"
	"contactcenter-api/repository"
	"contactcenter-api/service/common"
	"context"
	"fmt"
	"strconv"
	"time"
)

func (s *RoleGroup) ExportRoleGroups(ctx context.Context, domainUuid, userUuid string, filter model.RoleGroupFilter) (string, error) {
	fileType := "xlsx"
	total, _, err := repository.RoleGroupRepo.GetRoleGroup(ctx, domainUuid, -1, 0, filter)
	if err != nil {
		log.Error(err)
		return "", err
	}

	timeStr := util.TimeToStringLayout(time.Now(), "2006_01_02_15_04_05")
	fileName := "Role_Group_" + timeStr + "." + fileType
	exportMap := []string{fileName, util.TimeToString(time.Now()), "", fmt.Sprintf("%d", total), "In Progress", domainUuid, "role_group"}
	if err := common.SetExportValue(domainUuid, fileName, exportMap); err != nil {
		log.Error(err)
		return "", err
	}
	filePath, err := s.generateExportRoleGroups(ctx, domainUuid, userUuid, fileName, fileName, fileType, exportMap, &filter)
	if err != nil {
		log.Error(err)

	}
	return filePath, err
}

func (service *RoleGroup) generateExportRoleGroups(ctx context.Context, domainUuid, userUuid, exportName, fileName, fileType string, exportMap []string, param *model.RoleGroupFilter) (string, error) {
	total, _, err := repository.RoleGroupRepo.GetRoleGroup(ctx, domainUuid, -1, 0, *param)
	if err != nil {
		log.Error(err)
		return "", err
	}

	limit := 1000
	offset := 0
	headers := make([]string, 0)
	headers = append(headers, "Role Group", "Status", "Description")
	rows := make([][]string, 0)

	for offset < total {
		_, roleGroups, err := repository.RoleGroupRepo.GetRoleGroup(ctx, domainUuid, limit, offset, *param)
		if err != nil {
			log.Error(err)
			return "", err
		}
		for _, roleGroup := range *roleGroups {
			row := make([]string, 0)
			row = append(row, roleGroup.RoleGroupName, strconv.FormatBool(roleGroup.Status), roleGroup.Description)
			rows = append(rows, row)
		}

		offset += limit
		percentComplete := (float64(offset) / float64(total)) * 100
		exportMap[3] = fmt.Sprintf("%d", total)
		exportMap[4] = "In Progress (" + fmt.Sprintf("%.2f", percentComplete) + "%)"
		if err := common.SetExportValue(domainUuid, fileName, exportMap); err != nil {
			log.Error(err)
			return "", err
		}
	}
	dir := constants.EXPORT_DIR + "role_group/"
	if fileType == "xlsx" {
		if err := util.HandleExcelStreamWriter(fileName, dir, headers, rows); err != nil {
			log.Error(err)
			return "", err
		}
	} else if fileType == "csv" {
		if err := util.HandleCSVStreamWriter(fileName, dir, headers, rows, exportMap); err != nil {
			log.Error(err)
			return "", err
		}
	}

	exportMap[2] = util.TimeToString(time.Now())
	exportMap[4] = "Done"
	if err := common.SetExportValue(domainUuid, fileName, exportMap); err != nil {
		log.Error(err)
		return "", err
	}
	fileName = dir + fileName
	return fileName, nil
}
