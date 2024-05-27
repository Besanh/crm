package service

import (
	"contactcenter-api/common/constants"
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	"contactcenter-api/repository"
	"contactcenter-api/service/common"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
	"github.com/xuri/excelize/v2"
)

func (s *UserCrm) PostExportUsers(ctx context.Context, domainUuid, userUuid string, limit, offset int, filter model.UserFilter) (int, any) {
	user, err := GetUserInfo(ctx, userUuid)
	if err != nil {
		return response.ServiceUnavailableMsg(err.Error())
	}
	if user.Level == constants.MANAGER {
		filter.ManageUserUuids = make([]string, 0)
		for _, u := range user.ManageUsers {
			filter.ManageUserUuids = append(filter.ManageUserUuids, u.UserUuid)
		}
	} else if user.Level == constants.LEADER {
		filter.ManageUserUuids = make([]string, 0)
		for _, u := range user.ManageUsers {
			filter.ManageUserUuids = append(filter.ManageUserUuids, u.UserUuid)
		}
	}

	timeStr := util.TimeToStringLayout(time.Now(), "2006_01_02_15_04_05")
	exportName := "Users_Export_" + timeStr + ".xlsx"
	exportMap := []string{exportName, util.TimeToString(time.Now()), "", "0", "In Progress", domainUuid, userUuid, "users"}
	if err := common.SetExportValue(userUuid, exportName, exportMap); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	go func(ctx context.Context, domainUuid, exportName string, exportMap []string, filter model.UserFilter) {
		if err := s.generateExportUsers(ctx, domainUuid, userUuid, exportName, exportMap, filter); err != nil {
			log.Error(err)
		}
	}(context.Background(), domainUuid, exportName, exportMap, filter)

	return response.Created(map[string]any{
		"export_name": exportName,
		"status":      "In Progress",
	})
}

func (s *UserCrm) generateExportUsers(ctx context.Context, domainUuid, userUuid, exportName string, exportMap []string, param model.UserFilter) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	total, _, err := repository.UserCrmRepo.GetUserCrms(ctx, domainUuid, 1, 0, param)
	if err != nil {
		log.Error(err)
		return err
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"

	if err := file.SetSheetProps(sheetName, nil); err != nil {
		log.Error(err)
	}

	if err := file.SetColWidth(sheetName, "A", "N", 20); err != nil {
		log.Error(err)
	}
	//set row height
	styleBorder, _ := file.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FCD5B4"}, Pattern: 1},
		Alignment: &excelize.Alignment{WrapText: true},
	})

	columnsHeader := []any{"Username", "Fullname", "Email", "Level", "QRCode", "Domain", "Outbound Proxy", "Extension", "Display Name"}

	for i := 0; i < len(columnsHeader); i++ {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		file.SetCellStyle(sheetName, cell, cell, styleBorder)
		str, _ := columnsHeader[i].(string)
		if err := file.SetCellRichText(sheetName, cell, []excelize.RichTextRun{
			{
				Text: str,
			},
		}); err != nil {
			log.Error(err)
		}
	}
	_, _ = file.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#ffffff"}, Pattern: 1},
		Alignment: &excelize.Alignment{WrapText: false, Horizontal: "center"},
	})

	const (
		INDEX_USERNAME int = 1
		INDEX_FULLNAME int = 2
		INDEX_EMAIL    int = 3
		INDEX_LEVEL    int = 4
		INDEX_QRCODE   int = 5
		//qrcode 3
		INDEX_DOMAIN_NAME    int = 6
		INDEX_OUTBOUND_PROXY int = 7
		INDEX_EXTENSION      int = 8
		INDEX_DISPLAY_NAME   int = 9
	)
	limitPerPart := 50
	indexRow := 2
	offset := 0

	for offset < total {
		_, users, err := repository.UserCrmRepo.GetUserCrms(ctx, domainUuid, limitPerPart, offset, param)
		if err != nil {
			log.Error(err)
			return err
		}
		for _, user := range users {
			extension, err := repository.ExtensionRepo.GetExtensionInfoByUserUuid(ctx, domainUuid, user.UserUuid)
			if err != nil {
				log.Error(err)
				return err
			}

			file.SetRowHeight(sheetName, indexRow, 250)
			cell, _ := excelize.CoordinatesToCellName(INDEX_USERNAME, indexRow)
			if err := file.SetCellValue(sheetName, cell, user.Username); err != nil {
				log.Error(err)
			}
			cell, _ = excelize.CoordinatesToCellName(INDEX_FULLNAME, indexRow)
			fullname := fmt.Sprintf("%s %s %s", user.LastName, user.MiddleName, user.FirstName)
			if err := file.SetCellValue(sheetName, cell, strings.TrimSpace(fullname)); err != nil {
				log.Error(err)
			}
			cell, _ = excelize.CoordinatesToCellName(INDEX_EMAIL, indexRow)
			if err := file.SetCellValue(sheetName, cell, user.Email); err != nil {
				log.Error(err)
			}
			cell, _ = excelize.CoordinatesToCellName(INDEX_LEVEL, indexRow)
			if err := file.SetCellValue(sheetName, cell, user.Level); err != nil {
				log.Error(err)
			}
			// qrcode
			if extension != nil {
				grandStream := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>`+
					`<AccountConfig version="1">`+
					`<Account>`+
					`<RegisterServer>%s</RegisterServer>`+
					`<OutboundServer>%s</OutboundServer>`+
					`<UserID>%s</UserID>`+
					`<AuthID>%s</AuthID>`+
					`<AuthPass>%s</AuthPass>`+
					`<AccountName>%s</AccountName>`+
					`<DisplayName>%s</DisplayName>`+
					`<Dialplan>{x+|\+x+|*x+|*xx*x+}</Dialplan>`+
					`<RandomPort>0</RandomPort>`+
					`<SecOutboundServer />`+
					`<Voicemail>*97</Voicemail>`+
					`</Account>`+
					`</AccountConfig>`,
					extension.DomainName,
					PBXInfo.PBXOutboundProxy,
					extension.Extension,
					extension.Extension,
					extension.Password,
					extension.Extension,
					extension.Extension+"@"+extension.DomainName,
				)

				byteArr, err := qrcode.Encode(grandStream, qrcode.Medium, 180)
				if err != nil {
					log.Error(err)
				}
				file.SetColWidth(sheetName, "E", "E", 30)
				pic := &excelize.Picture{
					Extension: "jpg",
					File:      byteArr,
				}
				if err := file.AddPictureFromBytes(sheetName, fmt.Sprintf("E%d", indexRow), pic); err != nil {
					log.Error(err)
				}

				cell, _ = excelize.CoordinatesToCellName(INDEX_DOMAIN_NAME, indexRow)
				if err := file.SetCellValue(sheetName, cell, extension.DomainName); err != nil {
					log.Error(err)
				}
				cell, _ = excelize.CoordinatesToCellName(INDEX_OUTBOUND_PROXY, indexRow)
				if err := file.SetCellValue(sheetName, cell, PBXInfo.PBXOutboundProxy); err != nil {
					log.Error(err)
				}

				cell, _ = excelize.CoordinatesToCellName(INDEX_EXTENSION, indexRow)
				if extension.Extension != "" {
					if err := file.SetCellValue(sheetName, cell, extension.Extension); err != nil {
						log.Error(err)
					}
				}
				cell, _ = excelize.CoordinatesToCellName(INDEX_DISPLAY_NAME, indexRow)
				if err := file.SetCellValue(sheetName, cell, extension.Extension+"@"+extension.DomainName); err != nil {

					log.Error(err)
				}
			}
			indexRow++
		}
		if offset < total {
			percentComplete := (float64(offset) / float64(total)) * 100
			exportMap[3] = fmt.Sprintf("%d", total)
			exportMap[4] = "In Progress (" + fmt.Sprintf("%.2f", percentComplete) + "%)"
			if err := common.SetExportValue(userUuid, exportName, exportMap); err != nil {
				log.Error(err)
				return err
			}
		}
		offset += limitPerPart
	}
	exportDir := constants.EXPORT_DIR + "users/"
	if err = os.MkdirAll(filepath.Dir(exportDir), 0755); err != nil {
		log.Error(err)
	}
	if err := file.SaveAs(exportDir + exportName); err != nil {
		log.Error(err)
		return err
	}
	exportMap[2] = util.TimeToString(time.Now())
	exportMap[3] = fmt.Sprintf("%d", total)
	exportMap[4] = "Done"
	if err := common.SetExportValue(userUuid, exportName, exportMap); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
