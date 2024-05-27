package calendar

import (
	"contactcenter-api/common/util"
	"os"
)

func GetAttachmentDir(fileName string) (string, error) {
	dirUserVoiceDirectory := util.PUBLIC_DIR + fileName
	if _, err := os.Stat(dirUserVoiceDirectory); os.IsNotExist(err) {
		return "", err
	}
	return dirUserVoiceDirectory, nil
}
