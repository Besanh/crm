package util

import (
	"bytes"
	"contactcenter-api/common/log"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

const (
	PUBLIC_DIR     = "/root/go/src/public/"
	UPLOAD_DIR     = PUBLIC_DIR + "upload/"
	UPLOAD_SMS_DIR = PUBLIC_DIR + "sms/"
	DOWNLOAD_DIR   = PUBLIC_DIR + "download/"
)

func DecodeAndSaveImageBase64(data string, dir string, fileName string) (string, error) {
	idx := strings.Index(data, ";base64,")
	if idx < 0 {
		return fileName, errors.New("invalid data")
	}
	imageType := data[11:idx]
	unbased, err := base64.StdEncoding.DecodeString(data[idx+8:])
	if err != nil {
		return fileName, err
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
	fileName = fileName + "." + imageType
	dir = dir + fileName
	f, err := os.OpenFile(dir, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return fileName, err
	}
	defer f.Close()
	r := bytes.NewReader(unbased)
	switch imageType {
	case "png":
		if im, err := png.Decode(r); err != nil {
			return fileName, err
		} else {
			png.Encode(f, im)
		}
	case "jpeg":
		if im, err := jpeg.Decode(r); err != nil {
			return fileName, err
		} else {
			jpeg.Encode(f, im, nil)
		}
	default:
		return fileName, errors.New("invalid image type")
	}
	return fileName, nil
}

func HandleExcelStreamWriter(fileName string, saveDir string, headers []string, rows [][]string) error {
	file := excelize.NewFile()
	SHEET1 := "Sheet1"
	index := 1
	streamWriter, err := file.NewStreamWriter(SHEET1)
	if err != nil {
		log.Error(err)
		return err
	}
	styleID, _ := file.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FCD5B4"}, Pattern: 1},
		Alignment: &excelize.Alignment{WrapText: true},
	})
	cell, _ := excelize.CoordinatesToCellName(1, index)
	values := []any{}
	for _, header := range headers {
		values = append(values, excelize.Cell{
			Value:   header,
			StyleID: styleID,
		})
	}
	if err := streamWriter.SetRow(cell, values); err != nil {
		log.Error(err)
		return err
	}
	index++
	styleID, _ = file.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#ffffff"}, Pattern: 1},
		Alignment: &excelize.Alignment{WrapText: false, Horizontal: "left"},
	})
	for _, row := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, index)
		values := []any{}
		for _, cellValue := range row {
			values = append(values, excelize.Cell{
				Value:   cellValue,
				StyleID: styleID,
			})
		}
		if err := streamWriter.SetRow(cell, values); err != nil {
			log.Error(err)
			break
		}
		index++
	}
	if err := streamWriter.Flush(); err != nil {
		log.Error(err)
		return err
	}
	_ = os.MkdirAll(filepath.Dir(saveDir), 0755)
	if err := file.SaveAs(saveDir + fileName); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func HandleCSVStreamWriter(fileName string, saveDir string, headers []string, rows [][]string, exportMap []string) error {
	_ = os.MkdirAll(filepath.Dir(saveDir), 0755)
	f, err := os.Create(saveDir + fileName)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	w.UseCRLF = true
	defer f.Close()
	bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
	f.Write(bomUtf8)
	var errExport error
	if errExport != nil {
		log.Error(errExport)
		return errExport
	}

	if err := w.Write(headers); err != nil {
		log.Error(err)
		return err
	}

	for _, row := range rows {
		records := make([]string, len(row))
		for k, v := range row {
			records[k] = fmt.Sprintf("%v", v)
		}
		if err = w.Write(records); err != nil {
			log.Error(err)
			return err
		}
	}
	w.Flush()
	if err := f.Close(); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
