package api

import (
	"contactcenter-api/api"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	authMdw "contactcenter-api/middleware/auth"
	"contactcenter-api/service"
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Workday struct {
	workdayService service.IWorkDay
}

func NewWorkday(r *gin.Engine, workdayService service.IWorkDay) {
	handler := &Workday{
		workdayService: workdayService,
	}
	Group := r.Group(constants.VERSION + "/crm/workday")
	{
		Group.GET("", authMdw.AuthMiddleware(), handler.GetWorkdays)
		Group.POST("", authMdw.AuthMiddleware(), handler.PostWorkdays)
		Group.GET(":id", authMdw.AuthMiddleware(), handler.GetWorkdayById)
		Group.GET("/workday-id/:id", authMdw.AuthMiddleware(), handler.GetWorkdayByWorkdayId)
		Group.DELETE(":id", authMdw.AuthMiddleware(), handler.DeleteWorkday)
		Group.DELETE("workday-id/:id", authMdw.AuthMiddleware(), handler.DeleteWorkdayByWorkdayId)
	}
	Group2 := r.Group("v1/crm/holiday")
	{
		Group2.POST("", authMdw.AuthMiddleware(), handler.PostHoliday)
	}
}

func (handler *Workday) GetWorkdays(c *gin.Context) {
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
	filter := model.WorkDayFilter{
		WorkDayName: c.Query("work_day_name"),
		StartTime:   c.Query("start_time"),
		EndTime:     c.Query("end_time"),
		Status:      status,
	}
	limit, offset := util.GetLimitOffset(c.Query("limit"), c.Query("offset"))

	code, result := handler.workdayService.GetWorkDays(c, domainUuid, userUuid, limit, offset, filter)
	c.JSON(code, result)
}

func (handler *Workday) PostWorkdays(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	workdays := make([]model.WorkDay, 0)
	if err := c.BindJSON(&workdays); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, result := handler.workdayService.PostWorkDays(c, domainUuid, userUuid, workdays)
	c.JSON(code, result)
}

func (handler *Workday) DeleteWorkday(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	code, result := handler.workdayService.DeleteWorkDay(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (handler *Workday) DeleteWorkdayByWorkdayId(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	code, result := handler.workdayService.DeleteWorkdayByWorkdayId(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (handler *Workday) PostHoliday(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	holiday := model.Holiday{}
	if err := c.ShouldBindJSON(&holiday); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, result := handler.workdayService.PostHoliday(c, domainUuid, userUuid, holiday)
	c.JSON(code, result)
}

func (handler *Workday) GetWorkdayById(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	holiday := model.Holiday{}
	if err := c.ShouldBindJSON(&holiday); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, result := handler.workdayService.PostHoliday(c, domainUuid, userUuid, holiday)
	c.JSON(code, result)
}

func (hanlder *Workday) GetWorkdayByWorkdayId(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	workdayId := c.Param("id")
	if len(workdayId) < 1 {
		c.JSON(response.BadRequestMsg("workdayId is empty"))
		return
	}

	code, result := hanlder.workdayService.GetWorkdayByWorkdayId(c, domainUuid, userUuid, workdayId)
	c.JSON(code, result)
}
