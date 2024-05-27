package api

import (
	"contactcenter-api/api"
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

type Unit struct {
	unitService service.IUnit
}

func NewUnit(r *gin.Engine, service service.IUnit) {
	handler := Unit{
		unitService: service,
	}
	Group := r.Group("v1/crm/unit")
	{
		Group.POST("", authMdw.AuthMiddleware(), handler.PostUnit)
		Group.GET("", authMdw.AuthMiddleware(), handler.GetUnits)
		Group.GET("view/:view_type", authMdw.AuthMiddleware(), handler.GetUnitsRecursive)
		Group.GET(":id", authMdw.AuthMiddleware(), handler.GetUnitById)
		Group.PUT(":id", authMdw.AuthMiddleware(), handler.PutUnitById)
		Group.DELETE(":id", authMdw.AuthMiddleware(), handler.DeleteUnitById)
		Group.GET("view-render", authMdw.AuthMiddleware(), handler.GetUnitTreeRender)
		Group.PATCH(":id/logo", authMdw.AuthMiddleware(), handler.PatchUnitLogo)
		Group.GET("logo_file/:file_name/download", handler.GetUnitLogo)
		Group.DELETE(":id/logo", authMdw.AuthMiddleware(), handler.DeleteUnitLogo)
		Group.POST("export", authMdw.AuthMiddleware(), handler.ExportUnits)
		Group.GET("tree-parent", authMdw.AuthMiddleware(), handler.GetTreeParent)
		Group.GET("tree-formular", authMdw.AuthMiddleware(), handler.GetTreeFormular)
		Group.GET("tree-exclude-main-current-unit", authMdw.AuthMiddleware(), handler.GetTreeExcludeMainCurrentUnit)
	}
}
func (handler *Unit) PostUnit(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	unit := model.Unit{}
	if err := c.BindJSON(&unit); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, validSchema := validator.CheckSchema("unit/post.json", unit)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}

	code, result := handler.unitService.PostUnit(c, domainUuid, userUuid, unit)
	c.JSON(code, result)
}

func (handler *Unit) GetUnits(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	limit, offset := util.GetLimitOffset(c.Query("limit"), c.Query("offset"))
	var status sql.NullBool
	if len(c.Query("status")) > 0 {
		status.Valid = true
		statusTmp, _ := strconv.ParseBool(c.Query("status"))
		status.Bool = statusTmp
	}

	filter := model.UnitFilter{
		ParentUnitUuid: c.Query("parent_unit_uuid"),
		UnitUuid:       c.Query("unit_uuid"),
		UnitName:       c.Query("unit_name"),
		UnitCode:       c.Query("unit_code"),
		UnitLeader:     c.Query("unit_leader"),
		Status:         status,
	}

	code, result := handler.unitService.GetUnits(c, domainUuid, userUuid, limit, offset, filter)
	c.JSON(code, result)
}

func (handler *Unit) GetUnitsRecursive(c *gin.Context) {
	domainUuid, _, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	viewType := c.Param("view_type")

	code, result := handler.unitService.GetUnitsRecursive(c, domainUuid, viewType)
	c.JSON(code, result)
}

func (handler *Unit) GetUnitTreeRender(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, result := handler.unitService.GetUnitTreeRender(c, domainUuid, userUuid)
	c.JSON(code, result)
}

func (handler *Unit) GetUnitById(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("unit id is empty"))
		return
	}

	code, result := handler.unitService.GetUnitById(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (handler *Unit) PutUnitById(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("unit id is empty"))
		return
	}

	unit := model.Unit{}
	if err := c.BindJSON(&unit); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, validSchema := validator.CheckSchema("unit/post.json", unit)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}

	code, result := handler.unitService.PutUnitById(c, domainUuid, userUuid, id, unit)
	c.JSON(code, result)
}

func (handler *Unit) DeleteUnitById(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("unit id is empty"))
		return
	}

	code, result := handler.unitService.DeleteUnitById(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (handler *Unit) PatchUnitLogo(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("unit is empty"))
		return
	}

	body := make(map[string]any)
	if err := c.BindJSON(&body); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	data, ok := body["data"].(string)
	if !ok || len(data) < 1 {
		c.JSON(response.BadRequestMsg("data is empty"))
		return
	}

	code, result := handler.unitService.PatchUnitLogo(c, domainUuid, userUuid, id, data)
	c.JSON(code, result)
}

func (handler *Unit) GetUnitLogo(c *gin.Context) {
	fileName := c.Param("file_name")
	if len(fileName) < 1 {
		c.JSON(response.BadRequestMsg("file name is empty"))
		return
	}

	var errs error
	var result string
	result, errs = handler.unitService.GetUnitLogo(c, fileName)
	if errs == nil {
		fileByte, err := os.ReadFile(result)
		if err != nil {
			c.JSON(response.ServiceUnavailableMsg(err.Error()))
			return
		}
		contentType := http.DetectContentType(fileByte)
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%q", result))
		c.Writer.Header().Add("Content-Type", contentType)
		c.Data(http.StatusOK, contentType, fileByte)
		return
	} else {
		c.JSON(http.StatusNotFound, map[string]any{
			"content":   "not found",
			"file_name": fileName,
		})
	}
}

func (handler *Unit) DeleteUnitLogo(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("unit uuid is empty"))
		return
	}

	code, result := handler.unitService.DeleteUnitLogo(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (handler *Unit) ExportUnits(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	var status sql.NullBool
	if len(c.Query("status")) > 0 {
		status.Valid = true
		statusTmp, _ := strconv.ParseBool(c.Query("status"))
		status.Bool = statusTmp
	}

	filter := model.UnitFilter{
		ParentUnitUuid: c.Query("parent_unit_uuid"),
		UnitUuid:       c.Query("unit_uuid"),
		UnitName:       c.Query("unit_name"),
		UnitCode:       c.Query("unit_code"),
		UnitLeader:     c.Query("unit_leader"),
		Status:         status,
	}

	filePath, err := handler.unitService.ExportUnits(c, domainUuid, userUuid, filter)
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

func (handler *Unit) GetTreeParent(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	filter := model.ParentUnitFilter{
		Level: c.Query("level"),
	}
	code, result := handler.unitService.GetUnitParentTree(c, domainUuid, userUuid, filter)
	c.JSON(code, result)
}

func (handler *Unit) GetTreeFormular(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	filter := model.UnitFilter{
		Level:         c.Query("level"),
		FromUnitLevel: c.Query("from_unit_level"),
		ToUnitLevel:   c.Query("to_unit_level"),
		Encompass:     util.ParseQueryArray(c.QueryArray("encompass")),
		Formular:      c.Query("formular"),
	}

	code, result := handler.unitService.GetUnitChildTree(c, domainUuid, userUuid, filter)
	c.JSON(code, result)
}

func (handler *Unit) GetTreeExcludeMainCurrentUnit(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	filter := model.UnitFilter{
		Level:         c.Query("level"),
		FromUnitLevel: c.Query("from_unit_level"),
		ToUnitLevel:   c.Query("to_unit_level"),
		Encompass:     util.ParseQueryArray(c.QueryArray("encompass")),
		Formular:      c.Query("formular"),
	}

	code, result := handler.unitService.GetTreeExcludeCurrentUnit(c, domainUuid, userUuid, filter)
	c.JSON(code, result)
}
