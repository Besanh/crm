package api

import (
	"contactcenter-api/api"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/response"
	authMdw "contactcenter-api/middleware/auth"
	"contactcenter-api/service"

	"github.com/gin-gonic/gin"
)

type UserLog struct {
	userLogService service.IUserLog
}

func NewUserLog(r *gin.Engine, service service.IUserLog) {
	handler := UserLog{
		userLogService: service,
	}
	Group := r.Group(constants.VERSION + "/crm/user-log")
	{
		Group.PATCH(":call_id", authMdw.AuthMiddleware(), handler.PatchUserLog)
	}
}

func (handler *UserLog) PatchUserLog(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	jsonBody := make(map[string]any)
	if err := c.BindJSON(&jsonBody); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	callId := c.Param("call_id")
	if len(callId) < 1 {
		c.JSON(response.BadRequestMsg("call_id is missing"))
		return
	}
	dispoSec, _ := jsonBody["dispo_sec"].(string)

	code, result := handler.userLogService.PatchUserLog(c, domainUuid, userUuid, callId, dispoSec)
	c.JSON(code, result)
}
