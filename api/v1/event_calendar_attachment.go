package api

import (
	"contactcenter-api/api"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	"contactcenter-api/service"
	"fmt"
	"net/http"
	"os"
	"time"

	authMdw "contactcenter-api/middleware/auth"

	"github.com/gin-gonic/gin"
)

type EventCalendarAttachment struct {
	eventCalendarAttachment service.IEventCalendarAttachment
}

func NewEventCalendarAttachment(r *gin.Engine, eventCalendarAttachmentService service.IEventCalendarAttachment) {
	handler := &EventCalendarAttachment{
		eventCalendarAttachment: eventCalendarAttachmentService,
	}

	Group := r.Group(constants.VERSION + "/crm/event-calendar-attachment")
	r.MaxMultipartMemory = 10 << 20
	{
		Group.POST(":id", authMdw.AuthMiddleware(), handler.PostEventCalendarAttachment)
		Group.GET(":id/download", handler.DownloadEventCalendarAttachment)
		Group.DELETE(":id", authMdw.AuthMiddleware(), handler.DeleteEventCalendarAttachmentById)
	}
}

func (handler *EventCalendarAttachment) PostEventCalendarAttachment(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	// attachment
	form, err := c.MultipartForm()
	if err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err))
		return
	}
	files := form.File["attachment[]"]
	dir := util.PUBLIC_DIR + domainUuid + "/event_calendar_attachments/" + userUuid
	code, result := handler.eventCalendarAttachment.PostEventCalendarAttachment(c, domainUuid, userUuid, id, dir, files)
	if code == http.StatusCreated {
		for _, file := range files {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				os.MkdirAll(dir, 0755)
			}
			filePath := dir + "/" + util.TimeToStringLayout(time.Now(), "2006_01_02_15_04") + "_" + file.Filename
			err := c.SaveUploadedFile(file, filePath)
			if err != nil {
				continue
			}
		}
	}
	c.JSON(code, result)
}

func (handler *EventCalendarAttachment) DownloadEventCalendarAttachment(c *gin.Context) {
	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	var errs error
	var result = model.EventCalendarAttachment{}
	result, errs = handler.eventCalendarAttachment.DownloadEventCalendarAttachment(c, id)
	if errs == nil {
		attachmentDir := result.PathFile

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
			"id":      id,
		})
	}
}

func (handler *EventCalendarAttachment) DeleteEventCalendarAttachmentById(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	code, result := handler.eventCalendarAttachment.DeleteEventCalendarAttachmentById(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}
