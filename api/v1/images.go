package api

import (
	"contactcenter-api/api"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/response"
	authMdw "contactcenter-api/middleware/auth"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func NewUploadLogo(r *gin.Engine, domainHost string) {

	Group := r.Group(constants.VERSION + "/crm/upload")
	{
		Group.POST("", authMdw.AuthMiddleware(), UploadImage(domainHost))
	}
}

func UploadImage(domainHost string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantUuid, _, err := api.GetInfoUser(c)
		if err != nil {
			c.JSON(response.BadRequestMsg(err))
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(response.ServiceUnavailableMsg(err.Error()))
			return
		}
		if file.Size > 50000000 { // 50MB
			c.JSON(response.BadRequestMsg("file size is too large . must be <= 50MB"))
			return
		}
		strArr := strings.Split(file.Filename, ".")
		fileExtension := strings.ToLower(strArr[1])
		if fileExtension != "png" && fileExtension != "jpg" && fileExtension != "jpeg" {
			c.JSON(response.BadRequestMsg("file must be image"))
			return
		}
		storageLogoPath := "public/logos/"
		if _, err := os.Stat(storageLogoPath); os.IsNotExist(err) {
			os.MkdirAll(storageLogoPath, 0755)
		}

		fileName := fmt.Sprintf("%s_%s", tenantUuid, strconv.Itoa(int(time.Now().UnixNano()))+"."+fileExtension)

		err = c.SaveUploadedFile(file, storageLogoPath+fileName)
		if err != nil {
			code, result := response.ServiceUnavailableMsg(err.Error())
			c.JSON(code, result)
			return
		}
		url := domainHost + "/" + constants.VERSION + "/crm/images/" + fileName

		c.JSON(response.Data(200, url))
	}
}
