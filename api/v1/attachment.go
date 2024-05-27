package api

import (
	"contactcenter-api/common/constants"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	"contactcenter-api/service"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Attachment struct {
	Attachment service.IAttachment
}

func NewAttachment(r *gin.Engine, attachment service.IAttachment) {
	handler := &Attachment{
		Attachment: attachment,
	}
	r.MaxMultipartMemory = 10 << 20
	Group := r.Group(constants.VERSION + "/crm/attachment")
	{
		Group.GET(":domain_uuid/:folder/:file_name", handler.GetAttachment)
		Group.GET(":domain_uuid/:folder/:file_name/download", handler.DownloadAttachment)
		Group.DELETE(":domain_uuid/:folder/:file_name/:entity_uuid", handler.DeleteAttachment)
	}
}

func (handler *Attachment) DownloadAttachment(c *gin.Context) {
	domainUuid := c.Param("domain_uuid")
	if len(domainUuid) < 1 {
		c.JSON(response.BadRequestMsg("domainUuid is empty"))
		return
	}
	folder := c.Param("folder")
	if len(folder) < 1 {
		c.JSON(response.BadRequestMsg("folder is empty"))
		return
	}
	// userUuid := c.Param("user_uuid")
	// if len(folder) < 1 {
	// 	c.JSON(response.BadRequestMsg("userUuid is empty"))
	// 	return
	// }

	fileName := c.Param("file_name")
	if len(fileName) < 1 {
		c.JSON(response.BadRequestMsg("fileName is empty"))
		return
	}

	var errs error
	_, errs = handler.Attachment.DownloadAttachment(c, domainUuid, folder, fileName)
	if errs == nil {
		attachmentDir := util.PUBLIC_DIR + domainUuid + "/" + folder + "/" + fileName

		fileByte, err := os.ReadFile(attachmentDir)
		if err != nil {
			code, result := response.ServiceUnavailableMsg(err.Error())
			c.JSON(code, result)
			return
		}
		contentType := http.DetectContentType(fileByte)
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%q", attachmentDir))
		c.Writer.Header().Add("Content-Type", contentType)
		c.Data(http.StatusOK, contentType, fileByte)
		return
	} else {
		c.JSON(http.StatusNotFound, map[string]any{
			"message": "not found",
			"id":      fileName,
		})
	}
}

func (handler *Attachment) DeleteAttachment(c *gin.Context) {
	domainUuid := c.Param("domain_uuid")
	if len(domainUuid) < 1 {
		c.JSON(response.BadRequestMsg("domainUuid is empty"))
		return
	}

	folder := c.Param("folder")
	if len(folder) < 1 {
		c.JSON(response.BadRequestMsg("folder is empty"))
		return
	}
	// userUuid := c.Param("user_uuid")
	// if len(folder) < 1 {
	// 	c.JSON(response.BadRequestMsg("userUuid is empty"))
	// 	return
	// }

	fileName := c.Param("file_name")
	if len(fileName) < 1 {
		c.JSON(response.BadRequestMsg("fileName is empty"))
		return
	}

	entityUuid := c.Param("entity_uuid")
	if len(entityUuid) < 1 {
		c.JSON(response.BadRequestMsg("entity_uuid is empty"))
		return
	}

	code, result := handler.Attachment.DeleteAttachment(c, domainUuid, folder, fileName, entityUuid)
	c.JSON(code, result)
}

func (handler *Attachment) GetAttachment(c *gin.Context) {
	domainUuid := c.Param("domain_uuid")
	if len(domainUuid) < 1 {
		c.JSON(response.BadRequestMsg("domainUuid is empty"))
		return
	}
	folder := c.Param("folder")
	if len(folder) < 1 {
		c.JSON(response.BadRequestMsg("folder is empty"))
		return
	}
	// userUuid := c.Param("user_uuid")
	// if len(folder) < 1 {
	// 	c.JSON(response.BadRequestMsg("userUuid is empty"))
	// 	return
	// }

	fileName := c.Param("file_name")
	if len(fileName) < 1 {
		c.JSON(response.BadRequestMsg("fileName is empty"))
		return
	}

	var errs error
	_, errs = handler.Attachment.DownloadAttachment(c, domainUuid, folder, fileName)
	if errs == nil {
		attachmentDir := util.PUBLIC_DIR + domainUuid + "/" + folder + "/" + fileName
		c.File(attachmentDir)
	} else {
		c.JSON(http.StatusNotFound, map[string]any{
			"message": "not found",
			"id":      fileName,
		})
	}
}
