package api

import (
	"contactcenter-api/api"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	authMdw "contactcenter-api/middleware/auth"
	"contactcenter-api/service"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type EventCalendar struct {
	eventCalendar service.IEventCalendar
}

func NewEventCalendar(r *gin.Engine, eventCalendarService service.IEventCalendar) {
	handler := EventCalendar{
		eventCalendar: eventCalendarService,
	}

	Group := r.Group(constants.VERSION + "/crm/event-calendar")
	r.MaxMultipartMemory = 10 << 20
	{
		Group.POST("", authMdw.AuthMiddleware(), handler.PostEventCalendar)
		Group.GET("", authMdw.AuthMiddleware(), handler.GetEventCalendar)
		Group.GET(":id", authMdw.AuthMiddleware(), handler.GetEventCalendarById)
		Group.PUT(":id", authMdw.AuthMiddleware(), handler.PutEventCalendarById)
		Group.PUT("info-todo/:id", authMdw.AuthMiddleware(), handler.PutEventCalendarInfoAndTodoById)
		Group.PATCH(":id", authMdw.AuthMiddleware(), handler.PatchEventCalendar)
		Group.PUT(":id/expire", authMdw.AuthMiddleware(), handler.PutEventCalendarExpire)
		Group.PUT(":id/status", authMdw.AuthMiddleware(), handler.PutEventCalendarStatusById)
	}
}

func (handler *EventCalendar) PostEventCalendar(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
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
	eventCalendar := model.EventCalendar{}
	eventCalendar, eventCalendarTodo, err := api.ParseEventCalendar(c, userUuid, domainUuid, eventCalendar)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, result := handler.eventCalendar.PostEventCalendar(c, domainUuid, userUuid, eventCalendar, eventCalendarTodo, files)
	if code == http.StatusCreated {
		for _, file := range files {
			dir := util.PUBLIC_DIR + "event_calendar_attachments/" + util.TimeToStringLayout(time.Now(), "2006_01_02")
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

func (handler *EventCalendar) GetEventCalendar(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	filter := model.EventCalendarFilter{
		EccUuid:   c.Query("ecc_uuid"),
		CreatedBy: userUuid,
	}
	filter.StartTime, filter.EndTime, err = util.ParseStartEndTime(c.Query("start_time"), c.Query("end_time"), true)
	if err != nil {
		c.JSON(response.BadRequest())
		return
	}

	code, result := handler.eventCalendar.GetEventCalendar(c, domainUuid, userUuid, filter)
	c.JSON(code, result)
}

func (handler *EventCalendar) GetEventCalendarById(c *gin.Context) {
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

	code, result := handler.eventCalendar.GetEventCalendarById(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (handler *EventCalendar) PutEventCalendarById(c *gin.Context) {
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

	eventCalendar := model.EventCalendar{}
	if err := c.BindJSON(&eventCalendar); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, result := handler.eventCalendar.PutEventCalendarById(c, domainUuid, userUuid, id, eventCalendar)
	c.JSON(code, result)
}

func (handler *EventCalendar) PutEventCalendarInfoAndTodoById(c *gin.Context) {
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

	eventCalendar := model.EventCalendar{}
	if err := c.BindJSON(&eventCalendar); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, result := handler.eventCalendar.PutEventCalendarInfoAndTodoById(c, domainUuid, userUuid, id, eventCalendar)
	c.JSON(code, result)
}

func (handler *EventCalendar) PatchEventCalendar(c *gin.Context) {
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

	eventCalendar := model.EventCalendar{}
	if err := c.BindJSON(&eventCalendar); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, result := handler.eventCalendar.PatchEventCalendarById(c, domainUuid, userUuid, id, eventCalendar)
	c.JSON(code, result)
}

func (handler *EventCalendar) PutEventCalendarExpire(c *gin.Context) {
	// domainUuid, userUuid, err := api.GetInfoUser(c)
	// if err != nil {
	// 	c.JSON(response.BadRequestMsg(err))
	// 	return
	// }

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}
	// eventCalendar := model.EventCalendar{}
	// if err := c.BindJSON(&eventCalendar); err != nil {
	// 	c.JSON(response.BadRequestMsg(err))
	// 	return
	// }

	// code, result := handler.eventCalendar.PutEventCalendarExpire(c, domainUuid, userUuid, id, eventCalendar)
	// c.JSON(code, result)
}

func (handler *EventCalendar) PutEventCalendarStatusById(c *gin.Context) {
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

	jsonBody := make(map[string]any)
	if err := c.BindJSON(&jsonBody); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	var status bool
	status, _ = jsonBody["status"].(bool)

	code, result := handler.eventCalendar.PutEventCalendarStatusById(c, domainUuid, userUuid, id, status)
	c.JSON(code, result)
}
