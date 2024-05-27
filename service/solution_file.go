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

func (s *Solution) ExportSolutions(ctx context.Context, domainUuid, userUuid, fileType string, filter model.SolutionFilter) (string, error) {
	fileType = "xlsx"
	total, _, err := repository.SolutionRepo.GetSolutions(ctx, domainUuid, -1, 0, filter)
	if err != nil {
		log.Error(err)
		return "", err
	}

	timeStr := util.TimeToStringLayout(time.Now(), "2006_01_02_15_04_05")
	fileName := "Solution_" + timeStr + "." + fileType
	exportMap := []string{fileName, util.TimeToString(time.Now()), "", fmt.Sprintf("%d", total), "In Progress", domainUuid, "role_group"}
	if err := common.SetExportValue(domainUuid, fileName, exportMap); err != nil {
		log.Error(err)
		return "", err
	}
	fileName, err = s.generateExportSolutions(ctx, domainUuid, userUuid, fileName, fileName, fileType, exportMap, &filter)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return fileName, nil

}

func (s *Solution) generateExportSolutions(ctx context.Context, domainUuid, userUuid, exportName, fileName, fileType string, exportMap []string, param *model.SolutionFilter) (string, error) {
	total, _, err := repository.SolutionRepo.GetSolutions(ctx, domainUuid, -1, 0, *param)
	if err != nil {
		log.Error(err)
		return "", err
	}

	limit := 1000
	offset := 0
	headers := make([]string, 0)
	headers = append(headers, "Solution Name", "Solution Code", "Status")
	rows := make([][]string, 0)

	for offset < total {
		_, solutions, err := repository.SolutionRepo.GetSolutions(ctx, domainUuid, -1, 0, *param)
		if err != nil {
			log.Error(err)
			return "", err
		}
		for _, solution := range *solutions {
			row := make([]string, 0)
			row = append(row, solution.SolutionName, solution.SolutionCode, strconv.FormatBool(solution.Status))
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
	dir := constants.EXPORT_DIR + "solution/"
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
	return fileName, nil
}
