package api

import (
	"contactcenter-api/api"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	"contactcenter-api/service"

	authMdw "contactcenter-api/middleware/auth"

	"github.com/gin-gonic/gin"
)

type EventCalendarCategory struct {
	eventCalendarCategory service.IEventCalendarCategory
}

func NewEventCalendarCategory(server *gin.Engine, eventCalendarCategoryService service.IEventCalendarCategory) {
	handler := &EventCalendarCategory{
		eventCalendarCategory: eventCalendarCategoryService,
	}

	Group := server.Group(constants.VERSION + "/crm/event-calendar-category")
	{
		Group.GET("", authMdw.AuthMiddleware(), handler.GetEventCalendarCategories)
		Group.GET(":id", authMdw.AuthMiddleware(), handler.GetEventCalendarCategoryById)
		Group.POST("", authMdw.AuthMiddleware(), handler.InsertEventCalendarCategory)
		Group.PUT(":id", authMdw.AuthMiddleware(), handler.UpdateEventCalendarCategoryById)
		Group.DELETE(":id", authMdw.AuthMiddleware(), handler.DeleteEventCalendarCategoryById)
	}
}

func (h *EventCalendarCategory) GetEventCalendarCategories(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	filter := model.EventCalendarCategoryFilter{
		Title:     c.Query("title"),
		Color:     c.Query("color"),
		CreatedBy: userUuid,
	}

	limit, offset := util.GetLimitOffset(c.Query("limit"), c.Query("offset"))

	code, result := h.eventCalendarCategory.GetEventCalendarCategory(c, domainUuid, filter, limit, offset)
	c.JSON(code, result)
}

func (h *EventCalendarCategory) InsertEventCalendarCategory(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	eventCalendarCategory := model.EventCalendarCategory{}
	if err := c.BindJSON(&eventCalendarCategory); err != nil {
		c.JSON(response.BadRequest())
		return
	}

	eventCalendarCategory.CreatedBy = userUuid
	code, result := h.eventCalendarCategory.InsertEventCalendarCategory(c, domainUuid, eventCalendarCategory)
	c.JSON(code, result)
}

func (h *EventCalendarCategory) GetEventCalendarCategoryById(c *gin.Context) {
	domainUuid, _, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	code, result := h.eventCalendarCategory.GetEventCalendarCategoryById(c, domainUuid, id)
	c.JSON(code, result)
}

func (h *EventCalendarCategory) UpdateEventCalendarCategoryById(c *gin.Context) {
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

	eventCalendarCategory := model.EventCalendarCategory{}
	if err := c.BindJSON(&eventCalendarCategory); err != nil {
		c.JSON(response.BadRequest())
		return
	}

	eventCalendarCategory.UpdatedBy = userUuid

	code, result := h.eventCalendarCategory.UpdateEventCalendarCategoryById(c, domainUuid, eventCalendarCategory)
	c.JSON(code, result)
}

func (h *EventCalendarCategory) DeleteEventCalendarCategoryById(c *gin.Context) {
	domainUuid, _, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	code, result := h.eventCalendarCategory.DeleteEventCalendarCategoryById(c, domainUuid, id)
	c.JSON(code, result)
}
