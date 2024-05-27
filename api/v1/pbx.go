package api

import (
	"contactcenter-api/api"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	authMdw "contactcenter-api/middleware/auth"
	"contactcenter-api/service"
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pbx struct {
	pbxService service.IPbx
}

func NewPbx(r *gin.Engine, pbxService service.IPbx) {
	handler := &Pbx{
		pbxService: pbxService,
	}

	Group := r.Group(constants.VERSION + "/crm/pbx")
	{
		Group.GET("", authMdw.AuthMiddleware(), handler.GetPbx)
		Group.GET("unit/:id", authMdw.AuthMiddleware(), handler.GetPbxByUnitId)
		Group.POST("", authMdw.AuthMiddleware(), handler.PostPbx)
		Group.GET(":id", authMdw.AuthMiddleware(), handler.GetPbxById)
		Group.PUT(":id", authMdw.AuthMiddleware(), handler.PutPbxById)
		Group.DELETE(":id", authMdw.AuthMiddleware(), handler.DeletePbxById)
	}
}

func (handler *Pbx) GetPbx(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err))
		return
	}

	var status sql.NullBool
	if len(c.Query("status")) > 0 {
		status.Valid = true
		statusTmp, _ := strconv.ParseBool(c.Query("status"))
		status.Bool = statusTmp
	}

	filter := model.PbxFilter{
		PbxName:   c.Query("pbx_name"),
		Status:    status,
		StartTime: c.Query("start_time"),
		EndTime:   c.Query("end_time"),
	}

	code, result := handler.pbxService.GetPbxs(c, domainUuid, userUuid, filter)
	c.JSON(code, result)
}

func (handler *Pbx) GetPbxByUnitId(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	code, result := handler.pbxService.GetPbxByUnitId(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (handler *Pbx) PostPbx(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err))
		return
	}

	body := model.Pbx{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, result := handler.pbxService.PostPbx(c, domainUuid, userUuid, body)
	c.JSON(code, result)
}

func (handler *Pbx) GetPbxById(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	code, result := handler.pbxService.GetPbxById(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}
func (handler *Pbx) PutPbxById(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	body := model.Pbx{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, result := handler.pbxService.PutPbxById(c, domainUuid, userUuid, id, body)
	c.JSON(code, result)
}

func (handler *Pbx) DeletePbxById(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	code, result := handler.pbxService.DeletePbxById(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}
