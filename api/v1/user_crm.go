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
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserCrm struct {
	userCrmService service.IUserCrm
}

func NewUserCrm(r *gin.Engine, userCrmService service.IUserCrm) {
	handler := &UserCrm{
		userCrmService: userCrmService,
	}

	Group := r.Group(constants.VERSION + "/crm/user-crm")
	{
		Group.POST("", authMdw.AuthMiddleware(), authMdw.CheckLevelManage(), handler.PostUserCrm)
		Group.GET("", authMdw.AuthMiddleware(), handler.GetUserCrms)
		Group.GET(":id", authMdw.AuthMiddleware(), handler.GetUserCrmById)
		Group.GET("view/:id", authMdw.AuthMiddleware(), handler.GetUserCrmViewById)
		Group.PUT(":id", authMdw.AuthMiddleware(), authMdw.CheckLevelManage(), handler.PutUserCrmById)
		Group.PATCH(":id", authMdw.AuthMiddleware(), authMdw.CheckLevelManage(), handler.PatchUserCrm)
		Group.DELETE(":id", authMdw.AuthMiddleware(), authMdw.CheckLevelManage(), handler.DeleteUserCrmById)
		Group.POST("file/export", authMdw.AuthMiddleware(), authMdw.CheckLevelManage(), handler.ExportExcelUsersCrm)
	}
}

func (handler *UserCrm) PostUserCrm(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	userCrm := model.UserPost{}
	if err := c.BindJSON(&userCrm); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, validSchema := validator.CheckSchema("user-crm/post.json", userCrm)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}

	code, result := handler.userCrmService.PostUserCrm(c, domainUuid, userUuid, userCrm)
	c.JSON(code, result)
}

func (handler *UserCrm) GetUserCrms(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	limit, offset := util.GetLimitOffset(c.Query("limit"), c.Query("offset"))
	filter := model.UserFilter{
		Fullname:  c.Query("fullname"),
		Name:      c.Query("username"),
		Email:     c.Query("email"),
		Level:     c.Query("level"),
		RoleUuid:  c.Query("role_uuid"),
		UnitUuid:  c.Query("unit_uuid"),
		Enabled:   c.Query("user_enabled"),
		Extension: c.Query("extension"),
		Common:    c.Query("common"),
	}

	code, result := handler.userCrmService.GetUserCrms(c, domainUuid, userUuid, limit, offset, filter)
	c.JSON(code, result)
}

func (handler *UserCrm) GetUserCrmById(c *gin.Context) {
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

	code, result := handler.userCrmService.GetUserCrmById(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (handler *UserCrm) GetUserCrmViewById(c *gin.Context) {
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

	code, result := handler.userCrmService.GetUserCrmViewById(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (handler *UserCrm) PutUserCrmById(c *gin.Context) {
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

	userCrm := model.UserPost{}
	if err := c.BindJSON(&userCrm); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	code, result := handler.userCrmService.PutUserCrmById(c, domainUuid, userUuid, id, userCrm)
	c.JSON(code, result)
}

func (handler *UserCrm) DeleteUserCrmById(c *gin.Context) {
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

	code, result := handler.userCrmService.DeleteUserCrmById(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (handler *UserCrm) ExportExcelUsersCrm(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	filter := model.UserFilter{
		Name:      c.Query("name"),
		Level:     c.Query("level"),
		Email:     c.Query("email"),
		UserUuid:  util.ParseQueryArray(c.QueryArray("user_uuid")),
		StartTime: c.Query("start_time"),
		EndTime:   c.Query("end_time"),
		Enabled:   c.Query("enable"),
		RoleUuid:  c.Query("role_uuid"),
	}
	limit, offset := util.GetLimitOffset(c.Query("limit"), c.Query("offset"))

	code, result := handler.userCrmService.PostExportUsers(c, domainUuid, userUuid, limit, offset, filter)
	c.JSON(code, result)
}

func (handler *UserCrm) PatchUserCrm(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("user is empty"))
		return
	}

	body := make(map[string]any)
	if err := c.BindJSON(&body); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	unitUuid, ok := body["unit_uuid"].(string)
	if !ok || len(unitUuid) < 1 {
		c.JSON(response.BadRequestMsg("unit_uuid is empty"))
		return
	}
	roleGroupUuid, ok := body["role_group_uuid"].(string)
	if !ok || len(roleGroupUuid) < 1 {
		c.JSON(response.BadRequestMsg("role_group_uuid is empty"))
		return
	}

	code, result := handler.userCrmService.PatchUserCrm(c, domainUuid, userUuid, id, unitUuid, roleGroupUuid)
	c.JSON(code, result)
}
