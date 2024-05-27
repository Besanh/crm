package api

import (
	"contactcenter-api/api"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	authMdw "contactcenter-api/middleware/auth"
	"contactcenter-api/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type CallTransfer struct {
	callTransferService service.ICallTransfer
}

func NewCallTransfer(r *gin.Engine, callTransferService service.ICallTransfer) {
	handler := CallTransfer{
		callTransferService: callTransferService,
	}

	Group := r.Group(constants.VERSION + "/crm/call-transfer")
	{
		Group.GET(":phone_number", authMdw.AuthMiddleware(), handler.GetCallTransfer)
		Group.POST("", authMdw.AuthMiddleware(), handler.CreateCallTransfer)
	}
}

func (handler *CallTransfer) GetCallTransfer(c *gin.Context) {
	domainUuid, _, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	phoneNumber := c.Param("phone_number")
	if len(phoneNumber) < 1 {
		c.JSON(response.BadRequestMsg("phone number is empty"))
		return
	}
	code, result := handler.callTransferService.GetCallTransferByInfo(c, domainUuid, phoneNumber)
	c.JSON(code, result)
}

func (handler *CallTransfer) CreateCallTransfer(c *gin.Context) {
	domainUuid, _, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	callTransfer := model.CallTransfer{}
	if err := c.BindJSON(&callTransfer); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	callTransfer.Token = strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
	code, result := handler.callTransferService.CreateCallTransfer(c, domainUuid, callTransfer)
	c.JSON(code, result)
}
