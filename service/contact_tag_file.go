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

func (s *ContactTag) ExportContactTags(ctx context.Context, domainUuid, userUuid, fileType string, filter model.ContactTagFilter) (string, error) {
	fileType = "xlsx"
	total, _, err := repository.ContactTagRepo.GetContactTags(ctx, domainUuid, -1, 0, filter)
	if err != nil {
		log.Error(err)
		return "", err
	}

	timeStr := util.TimeToStringLayout(time.Now(), "2006_01_02_15_04_05")
	fileName := "Contact_tag_" + timeStr + "." + fileType
	exportMap := []string{fileName, util.TimeToString(time.Now()), "", fmt.Sprintf("%d", total), "In Progress", domainUuid, "role_group"}
	if err := common.SetExportValue(domainUuid, fileName, exportMap); err != nil {
		log.Error(err)
		return "", err
	}
	filePath, err := s.generateExportContactTags(ctx, domainUuid, fileName, fileType, exportMap, &filter)
	if err != nil {
		return "", err
	}
	return filePath, err
}

func (s *ContactTag) generateExportContactTags(ctx context.Context, domainUuid, fileName, fileType string, exportMap []string, param *model.ContactTagFilter) (string, error) {
	total, _, err := repository.ContactTagRepo.GetContactTags(ctx, domainUuid, -1, 0, *param)
	if err != nil {
		log.Error(err)
		return "", err
	}

	limit := 1000
	offset := 0
	headers := make([]string, 0)
	headers = append(headers, "Tag Name", "Desctiption", "Status", "Created At")
	rows := make([][]string, 0)

	for offset < total {
		_, contactTags, err := repository.ContactTagRepo.GetContactTags(ctx, domainUuid, limit, offset, *param)
		if err != nil {
			log.Error(err)
			return "", err
		}
		for _, contactTag := range *contactTags {
			row := make([]string, 0)
			row = append(row, contactTag.TagName, contactTag.Description, strconv.FormatBool(contactTag.Status), fmt.Sprintf("%s", contactTag.CreatedAt))
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
	dir := constants.EXPORT_DIR + "contact_tag/"
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
	filePath := dir + fileName
	return filePath, nil
}
