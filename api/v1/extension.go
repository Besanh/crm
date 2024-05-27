package api

import (
	"contactcenter-api/api"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	"contactcenter-api/common/validator"
	authMdw "contactcenter-api/middleware/auth"
	IService "contactcenter-api/service"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Extension struct {
	extensionService IService.IExtension
}

func NewExtension(r *gin.Engine, extensionService IService.IExtension) {
	h := &Extension{
		extensionService: extensionService,
	}
	Group := r.Group(constants.VERSION + "/crm/extension")
	{
		Group.GET("", authMdw.AuthMiddleware(), h.GetExtensions)
		Group.GET("unit", authMdw.AuthMiddleware(), h.GetExtensionsInUnit)
		Group.GET(":id", authMdw.AuthMiddleware(), h.GetExtensionByIdOrExten)
		Group.GET(":id/loggedin", authMdw.AuthMiddleware(), h.CheckExtensionIsLoggedIn)
		Group.GET(":id/qrcode", authMdw.AuthMiddleware(), h.GetExtensionQrCode)
		Group.POST("", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), h.PostExtension)
		Group.PUT(":id", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), h.PutExtension)
		Group.DELETE(":id", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), h.DeleteExtension)
		Group.PATCH(":id", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), h.PatchExtension)
	}
}

func (h *Extension) CheckExtensionIsLoggedIn(c *gin.Context) {
	tenantUuid, ok := authMdw.GetUserDomainId(c)
	if !ok {
		c.JSON(response.BadRequest())
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(response.BadRequestMsg("id is missing"))
		c.Abort()
		return
	}
	code, result := h.extensionService.CheckExtensionIsLoggedIn(c, tenantUuid, id)
	c.JSON(code, result)
}

func (h *Extension) GetExtensions(c *gin.Context) {
	tenantUuid, ok := authMdw.GetUserDomainId(c)
	if !ok {
		c.JSON(response.BadRequest())
		return
	}
	userUuid, _ := authMdw.GetUserId(c)
	enabled := sql.NullBool{}
	if len(c.Query("enabled")) > 0 {
		enabledTmp, _ := strconv.ParseBool(c.Query("enabled"))
		enabled.Valid = true
		enabled.Bool = enabledTmp
	}
	filter := model.ExtensionFilter{
		Extension: c.Query("extension"),
		Username:  c.Query("username"),
		Enabled:   enabled,
		UnitUuid:  c.Query("unit_uuid"),
		Common:    c.Query("common"),
		Fullname:  c.Query("fullname"),
	}
	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))
	code, result := h.extensionService.GetExtensions(c, tenantUuid, userUuid, filter, limit, offset)
	c.JSON(code, result)
}

func (h *Extension) GetExtensionQrCode(c *gin.Context) {
	tenantUuid, ok := authMdw.GetUserDomainId(c)
	if !ok {
		c.JSON(response.BadRequest())
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(response.BadRequestMsg("id is missing"))
		c.Abort()
		return
	}
	app := c.Query("app")
	code, result := h.extensionService.GetExtensionQrCode(c, tenantUuid, id, app)
	c.JSON(code, result)
}

func (h *Extension) PostExtension(c *gin.Context) {
	tenantUuid, ok := authMdw.GetUserDomainId(c)
	if !ok {
		c.JSON(response.BadRequest())
		return
	}
	userUuid, _ := authMdw.GetUserId(c)
	body := make(map[string]any)
	if err := c.BindJSON(&body); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	if code, validSchema := validator.CheckSchema("extension/post.json", body); code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}
	extension := model.ExtensionPost{}
	if err := util.ParseAnyToAny(body, &extension); err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	if extension.Password != extension.ConfirmPassword {
		c.JSON(response.BadRequestMsg("password and confirm password not match"))
		return
	}
	// if extension.IsFollowMe && len(extension.FollowMe.Destination) < 1 {
	// 	c.JSON(response.BadRequestMsg("follow_me.destination is missing"))
	// 	return
	// }
	if extension.IsRingGroup && len(extension.RingGroup.Main) < 1 {
		c.JSON(response.BadRequestMsg("ring_group.main is missing"))
		return
	}
	// extension.IsDeleteFollowMe = false
	extension.IsDeleteRingGroup = false
	code, result := h.extensionService.PostExtension(c, tenantUuid, userUuid, extension)
	c.JSON(code, result)
}

func (h *Extension) PutExtension(c *gin.Context) {
	tenantUuid, ok := authMdw.GetUserDomainId(c)
	if !ok {
		c.JSON(response.BadRequest())
		return
	}
	userUuid, _ := authMdw.GetUserId(c)
	id := c.Param("id")
	body := make(map[string]any)
	if err := c.BindJSON(&body); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	if code, validSchema := validator.CheckSchema("extension/post.json", body); code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}
	extension := model.ExtensionPost{}
	if err := util.ParseAnyToAny(body, &extension); err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	if extension.Password != extension.ConfirmPassword {
		c.JSON(response.BadRequestMsg("password and confirm password not match"))
		return
	}
	// isFollowMe, ok := body["is_follow_me"].(bool)
	// if ok && !isFollowMe {
	// 	extension.IsDeleteFollowMe = true
	// }
	// if extension.IsFollowMe && len(extension.FollowMe.Destination) < 1 {
	// 	c.JSON(response.BadRequestMsg("follow_me.destination is missing"))
	// 	return
	// }
	if extension.IsRingGroup && len(extension.RingGroup.Main) < 1 {
		c.JSON(response.BadRequestMsg("ring_group.main is missing"))
		return
	}
	code, result := h.extensionService.PutExtension(c, tenantUuid, userUuid, id, extension)
	c.JSON(code, result)
}

func (h *Extension) DeleteExtension(c *gin.Context) {
	tenantUuid, ok := authMdw.GetUserDomainId(c)
	if !ok {
		c.JSON(response.BadRequest())
		return
	}
	userUuid, _ := authMdw.GetUserId(c)
	id := c.Param("id")
	code, result := h.extensionService.DeleteExtension(c, tenantUuid, userUuid, id)
	c.JSON(code, result)
}

func (h *Extension) GetExtensionByIdOrExten(c *gin.Context) {
	tenantUuid, ok := authMdw.GetUserDomainId(c)
	if !ok {
		c.JSON(response.BadRequest())
		return
	}
	userUuid, _ := authMdw.GetUserId(c)
	id := c.Param("id")
	code, result := h.extensionService.GetExtensionByIdOrExten(c, tenantUuid, userUuid, id)
	c.JSON(code, result)
}

func (h *Extension) PatchExtension(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}
	body := make(map[string]any)
	if err := c.BindJSON(&body); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	status, ok := body["status"].(bool)
	if !ok {
		c.JSON(response.BadRequestMsg("status is invalid"))
		return
	}
	code, result := h.extensionService.PatchExtension(c, domainUuid, userUuid, id, status)
	c.JSON(code, result)
}

func (h *Extension) GetExtensionsInUnit(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	enabled := sql.NullBool{}
	if len(c.Query("enabled")) > 0 {
		enabledTmp, _ := strconv.ParseBool(c.Query("enabled"))
		enabled.Valid = true
		enabled.Bool = enabledTmp
	}
	filter := model.ExtensionFilter{
		Extension: c.Query("extension"),
		Username:  c.Query("username"),
		Enabled:   enabled,
		UnitUuid:  c.Query("unit_uuid"),
		Common:    c.Query("common"),
		Fullname:  c.Query("fullname"),
	}
	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))
	code, result := h.extensionService.GetExtensionsInUnit(c, domainUuid, userUuid, filter, limit, offset)
	c.JSON(code, result)
}
