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

type RoleGroup struct {
	roleGroupService service.IRoleGroup
}

func NewRoleGroup(r *gin.Engine, roleGroupService service.IRoleGroup) {
	handler := &RoleGroup{
		roleGroupService: roleGroupService,
	}

	Group := r.Group(constants.VERSION + "/crm/role-group")
	{
		Group.GET("", authMdw.AuthMiddleware(), handler.GetRoleGroups)
		Group.POST("", authMdw.AuthMiddleware(), handler.PostRoleGroup)
		Group.GET(":id", authMdw.AuthMiddleware(), handler.GetRoleGroupById)
		Group.PUT(":id", authMdw.AuthMiddleware(), handler.PutRoleGroup)
		Group.DELETE(":id", authMdw.AuthMiddleware(), handler.DeleteRoleGroupById)
		Group.POST("export", authMdw.AuthMiddleware(), handler.ExportRoleGroups)
	}
}

func (handler *RoleGroup) GetRoleGroups(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	limit, offset := util.GetLimitOffset(c.Query("limit"), c.Query("offset"))
	statusTmp := c.Query("status")
	var status sql.NullBool
	if len(statusTmp) > 0 {
		status.Valid = true
		status.Bool, _ = strconv.ParseBool(statusTmp)
	}
	startTimeTmp := c.Query("start_time")
	var startTime sql.NullTime
	if len(startTimeTmp) > 0 {
		startTime.Valid = true
		startTime.Time = util.ParseTime(startTimeTmp)
	}
	endTimeTmp := c.Query("end_time")
	var endTime sql.NullTime
	if len(startTimeTmp) > 0 {
		endTime.Valid = true
		endTime.Time = util.ParseTime(endTimeTmp)
	}
	filter := model.RoleGroupFilter{
		RoleGroupName: c.Query("role_group_name"),
		Status:        status,
		StartTime:     startTime,
		EndTime:       endTime,
	}

	code, result := handler.roleGroupService.GetRoleGroup(c, domainUuid, userUuid, limit, offset, filter)
	c.JSON(code, result)
}

func (handler *RoleGroup) PostRoleGroup(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	roleGroup := model.RoleGroup{}
	if err := c.BindJSON(&roleGroup); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, validSchema := validator.CheckSchema("role_group/post.json", roleGroup)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}

	code, result := handler.roleGroupService.PostRoleGroup(c, domainUuid, userUuid, roleGroup)
	c.JSON(code, result)
}

func (handler *RoleGroup) GetRoleGroupById(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("role group id is empty"))
		return
	}

	code, result := handler.roleGroupService.GetRoleGroupById(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (handler *RoleGroup) PutRoleGroup(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("role group id is empty"))
		return
	}
	roleGroup := model.RoleGroup{}
	if err := c.BindJSON(&roleGroup); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, result := handler.roleGroupService.PutRoleGroup(c, domainUuid, userUuid, id, roleGroup)
	c.JSON(code, result)
}

func (handler *RoleGroup) DeleteRoleGroupById(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("role group id is empty"))
		return
	}

	code, result := handler.roleGroupService.DeleteRoleGroupById(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (handler *RoleGroup) ExportRoleGroups(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	statusTmp := c.Query("status")
	var status sql.NullBool
	if len(statusTmp) > 0 {
		status.Valid = true
		status.Bool, _ = strconv.ParseBool(statusTmp)
	}
	startTimeTmp := c.Query("start_time")
	var startTime sql.NullTime
	if len(startTimeTmp) > 0 {
		startTime.Valid = true
		startTime.Time = util.ParseTime(startTimeTmp)
	}
	endTimeTmp := c.Query("end_time")
	var endTime sql.NullTime
	if len(startTimeTmp) > 0 {
		endTime.Valid = true
		endTime.Time = util.ParseTime(endTimeTmp)
	}
	filter := model.RoleGroupFilter{
		RoleGroupName: c.Query("role_group_name"),
		Status:        status,
		StartTime:     startTime,
		EndTime:       endTime,
	}

	filePath, err := handler.roleGroupService.ExportRoleGroups(c, domainUuid, userUuid, filter)
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
