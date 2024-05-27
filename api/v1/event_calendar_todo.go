package api

import (
	"contactcenter-api/api"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/service"

	authMdw "contactcenter-api/middleware/auth"

	"github.com/gin-gonic/gin"
)

type EventCalendarTodo struct {
	eventCalendarTodo service.IEventCalendarTodo
}

func NewEventCalendarTodo(r *gin.Engine, eventCalendarTodoService service.IEventCalendarTodo) {
	handler := EventCalendarTodo{
		eventCalendarTodo: eventCalendarTodoService,
	}

	Group := r.Group(constants.VERSION + "/crm/event-calendar-todo")
	{
		Group.PUT(":typeUpdate/:id", authMdw.AuthMiddleware(), handler.PutEventCalendarTodo)
		Group.DELETE(":id", authMdw.AuthMiddleware(), handler.DeleteEventCalendarTodoById)
	}
}

func (handler *EventCalendarTodo) PutEventCalendarTodo(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	ecUuid := c.Param("id")
	if len(ecUuid) < 1 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	typeUpdate := c.Param("typeUpdate")
	if len(typeUpdate) < 1 {
		c.JSON(response.BadRequestMsg("typeUpdate is empty"))
		return
	}

	var code int
	var result any

	if typeUpdate == "bulk" {
		eventCalendarTodos := []model.EventCalendarTodo{}
		if err := c.ShouldBindJSON(&eventCalendarTodos); err != nil {
			c.JSON(response.BadRequestMsg(err))
			return
		}
		code, result = handler.eventCalendarTodo.PutEventCalendarTodos(c, domainUuid, userUuid, ecUuid, eventCalendarTodos)
		c.JSON(code, result)
	} else if typeUpdate == "single" {
		eventCalendarTodo := model.EventCalendarTodo{}
		if err := c.ShouldBindJSON(&eventCalendarTodo); err != nil {
			c.JSON(response.BadRequestMsg(err))
			return
		}
		code, result = handler.eventCalendarTodo.PutEventCalendarTodoById(c, domainUuid, userUuid, ecUuid, eventCalendarTodo)
		c.JSON(code, result)
	} else {
		c.JSON(response.BadRequestMsg("typeUpdate is invalid"))
	}

}

func (handle *EventCalendarTodo) DeleteEventCalendarTodoById(c *gin.Context) {
	domainUuid, _, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	ectUuid := c.Param("id")
	if len(ectUuid) < 1 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	code, result := handle.eventCalendarTodo.DeleteEventCalendarTodo(c, domainUuid, ectUuid)
	c.JSON(code, result)
}
