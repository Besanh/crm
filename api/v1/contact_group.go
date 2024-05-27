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

type ContactGroup struct {
	contactGroupService service.IContactGroup
}

func NewContactGroup(r *gin.Engine, contactGroupService service.IContactGroup) {
	handler := &ContactGroup{
		contactGroupService: contactGroupService,
	}

	Group := r.Group(constants.VERSION + "/crm/contact-group")
	{
		Group.POST("", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), handler.PostContactGroup)
		Group.GET("", authMdw.AuthMiddleware(), handler.GetContactGroups)
		Group.GET(":id", authMdw.AuthMiddleware(), handler.GetContactGroupById)
		Group.PUT(":id", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), handler.PutContactGroupById)
		Group.DELETE(":id", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), handler.DeleteContactGroup)
		Group.POST("/export", authMdw.AuthMiddleware(), handler.ExportContactGroups)
	}
}

func (handler *ContactGroup) PostContactGroup(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	contactGroup := model.ContactGroup{}
	if err := c.ShouldBindJSON(&contactGroup); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, validSchema := validator.CheckSchema("contact_group/post.json", contactGroup)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}

	code, result := handler.contactGroupService.PostContactGroup(c, domainUuid, userUuid, &contactGroup)
	c.JSON(code, result)
}

func (handler *ContactGroup) GetContactGroups(c *gin.Context) {
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
	filter := model.ContactGroupFilter{
		GroupName: c.Query("group_name"),
		GroupType: c.Query("group_type"),
		Status:    status,
		StartTime: c.Query("start_time"),
		EndTime:   c.Query("end_time"),
	}

	code, result := handler.contactGroupService.GetContactGroups(c, domainUuid, limit, offset, filter)
	c.JSON(code, result)
}

func (handler *ContactGroup) GetContactGroupById(c *gin.Context) {
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

	code, result := handler.contactGroupService.GetContactGroupById(c, domainUuid, id)
	c.JSON(code, result)
}

func (handler *ContactGroup) PutContactGroupById(c *gin.Context) {
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

	contactGroup := model.ContactGroup{}
	if err := c.ShouldBindJSON(&contactGroup); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, validSchema := validator.CheckSchema("contact_group/put.json", contactGroup)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}

	code, result := handler.contactGroupService.PutContactGroupById(c, domainUuid, userUuid, id, contactGroup)
	c.JSON(code, result)
}

func (handler *ContactGroup) DeleteContactGroup(c *gin.Context) {
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

	code, result := handler.contactGroupService.DeleteContactGroup(c, domainUuid, id)
	c.JSON(code, result)
}

func (handler *ContactGroup) ExportContactGroups(c *gin.Context) {
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
	filter := model.ContactGroupFilter{
		GroupName: c.Query("group_name"),
		GroupType: c.Query("group_type"),
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

	filePath, err := handler.contactGroupService.ExportContactGroups(c, domainUuid, userUuid, fileType, filter)
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
