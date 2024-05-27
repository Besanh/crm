package api

import (
	"contactcenter-api/api"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	"contactcenter-api/common/validator"
	authMdw "contactcenter-api/middleware/auth"
	"contactcenter-api/service"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContactTag struct {
	contactTagService service.IContactTag
}

func NewContactTag(r *gin.Engine, contactTagService service.IContactTag) {
	handler := &ContactTag{
		contactTagService: contactTagService,
	}

	Group := r.Group(constants.VERSION + "/crm/contact-tag")
	{
		Group.POST("", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), handler.PostContactTag)
		Group.GET("", authMdw.AuthMiddleware(), handler.GetContactTags)
		Group.GET(":id", authMdw.AuthMiddleware(), handler.GetContactTagById)
		Group.PUT(":id", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), handler.PutContactTagById)
		Group.DELETE(":id", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), handler.DeleteContactTag)
		Group.POST("export", authMdw.AuthMiddleware(), handler.ExportContactTags)
	}
}

func (handler *ContactTag) PostContactTag(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	contactTag := model.ContactTag{}
	if err := c.ShouldBindJSON(&contactTag); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, validSchema := validator.CheckSchema("contact_tag/post.json", contactTag)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}

	code, result := handler.contactTagService.PostContactTag(c, domainUuid, userUuid, &contactTag)
	c.JSON(code, result)
}

func (handler *ContactTag) GetContactTags(c *gin.Context) {
	domainUuid, _, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	limit, offset := util.GetLimitOffset(c.Query("limit"), c.Query("offset"))

	statusTmp := c.Query("status")
	status := sql.NullBool{}
	if len(statusTmp) > 0 {
		status.Valid = true
		activeBool, _ := strconv.ParseBool(statusTmp)
		status.Bool = activeBool
	}
	tagType := c.Query("tag_type")
	if len(tagType) > 0 && tagType == "all" {
		tagType = ""
	}
	filter := model.ContactTagFilter{
		TagName:   c.Query("tag_name"),
		TagType:   tagType,
		Status:    status,
		StartTime: c.Query("start_time"),
		EndTime:   c.Query("end_time"),
	}

	code, result := handler.contactTagService.GetContactTags(c, domainUuid, limit, offset, filter)
	c.JSON(code, result)
}

func (handler *ContactTag) GetContactTagById(c *gin.Context) {
	domainUuid, _, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	code, result := handler.contactTagService.GetContactTagById(c, domainUuid, id)
	c.JSON(code, result)
}

func (handler *ContactTag) PutContactTagById(c *gin.Context) {
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

	contactTag := model.ContactTag{}
	if err := c.ShouldBindJSON(&contactTag); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, validSchema := validator.CheckSchema("contact_tag/put.json", contactTag)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}

	code, result := handler.contactTagService.PutContactTagById(c, domainUuid, userUuid, id, contactTag)
	c.JSON(code, result)
}

func (handler *ContactTag) DeleteContactTag(c *gin.Context) {
	domainUuid, _, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	code, result := handler.contactTagService.DeleteContactTag(c, domainUuid, id)
	c.JSON(code, result)
}

func (handler *ContactTag) ExportContactTags(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	statusTmp := c.Query("status")
	status := sql.NullBool{}
	if len(statusTmp) > 0 {
		status.Valid = true
		activeBool, _ := strconv.ParseBool(statusTmp)
		status.Bool = activeBool
	}
	tagType := c.Query("tag_type")
	if len(tagType) > 0 && tagType == "all" {
		tagType = ""
	}
	filter := model.ContactTagFilter{
		TagName:   c.Query("tag_name"),
		TagType:   tagType,
		Status:    status,
		StartTime: c.Query("start_time"),
		EndTime:   c.Query("end_time"),
		FileType:  c.Query("file_type"),
	}
	fileType := filter.FileType
	if len(fileType) < 1 {
		c.JSON(response.BadRequestMsg("file type is empty"))
		return
	}

	filePath, err := handler.contactTagService.ExportContactTags(c, domainUuid, userUuid, fileType, filter)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
	} else {
		fileByte, err := os.ReadFile(filePath)
		if err != nil {
			code, result := response.ServiceUnavailableMsg(err.Error())
			c.JSON(code, result)
			return
		}
		contentType := http.DetectContentType(fileByte)
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filePath))
		c.Writer.Header().Add("Content-Type", contentType)
		c.Data(http.StatusOK, contentType, fileByte)
		return
	}
}
