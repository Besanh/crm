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

type Domain struct {
	domainService service.IDomain
}

func NewDomain(r *gin.Engine, domainService service.IDomain) {
	handler := &Domain{
		domainService: domainService,
	}
	Group := r.Group(constants.VERSION + "/crm/domain")
	{
		Group.POST("", authMdw.AuthMiddleware(), handler.PostDomain)
		Group.GET("config", authMdw.AuthMiddleware(), handler.GetDomainConfigs)
		Group.GET(":id", authMdw.AuthMiddleware(), handler.GetDomainConfigById)
		Group.PUT("", authMdw.AuthMiddleware(), handler.PutDomainConfig)
		Group.PUT(":id", authMdw.AuthMiddleware(), handler.PutDomainConfigById)
		Group.POST(":id/config", authMdw.AuthMiddleware(), handler.PostDomainConfig)
	}
}

func (handler *Domain) GetDomainConfigs(c *gin.Context) {
	level, _ := authMdw.GetUserLevel(c)
	if level != authMdw.SUPERADMIN {
		c.JSON(response.ForbiddenLevel(level))
		return
	}
	limitQuery := c.Query("limit")
	offsetQuery := c.Query("offset")
	limit := util.ParseLimit(limitQuery)
	offset := util.ParseOffset(offsetQuery)
	code, result := handler.domainService.GetDomainConfigs(c, limit, offset)
	c.JSON(code, result)
}

func (handler *Domain) GetDomainConfigById(c *gin.Context) {
	level, _ := authMdw.GetUserLevel(c)
	if level != authMdw.SUPERADMIN {
		c.JSON(response.ForbiddenLevel(level))
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(response.BadRequestMsg("id is missing"))
		c.Abort()
		return
	}
	code, result := handler.domainService.GetDomainConfigById(c, id)
	c.JSON(code, result)
}

func (handler *Domain) PutDomainConfig(c *gin.Context) {
	level, _ := authMdw.GetUserLevel(c)
	if level != authMdw.SUPERADMIN {
		c.JSON(response.ForbiddenLevel(level))
		return
	}

	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	domainConfig := model.DomainConfig{}
	if err := c.BindJSON(&domainConfig); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	code, result := handler.domainService.PutDomainConfig(c, domainUuid, userUuid, domainConfig)
	c.JSON(code, result)
}

func (handler *Domain) PostDomainConfig(c *gin.Context) {
	level, _ := authMdw.GetUserLevel(c)
	if level != authMdw.SUPERADMIN {
		c.JSON(response.ForbiddenLevel(level))
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(response.BadRequestMsg("id is missing"))
		c.Abort()
		return
	}
	domainConfig := model.DomainConfigPut{}
	if err := c.BindJSON(&domainConfig); err != nil {
		code, result := response.BadRequest()
		c.JSON(code, result)
		return
	}
	code, validSchema := validator.CheckSchema("domain/domain_config.json", domainConfig)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}

	if err := domainConfig.ValidatePut(); err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	code, result := handler.domainService.PostDomainConfig(c, id, domainConfig)
	c.JSON(code, result)
}

func (handler *Domain) PutDomainConfigById(c *gin.Context) {
	level, _ := authMdw.GetUserLevel(c)
	if level != authMdw.SUPERADMIN {
		c.JSON(response.ForbiddenLevel(level))
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(response.BadRequestMsg("id is missing"))
		c.Abort()
		return
	}

	domainConfig := model.DomainConfigPut{}
	if err := c.BindJSON(&domainConfig); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	if err := domainConfig.ValidatePut(); err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	code, result := handler.domainService.PutDomainConfigById(c, id, domainConfig)
	c.JSON(code, result)
}

func (handler *Domain) PostDomain(c *gin.Context) {
	level, _ := authMdw.GetUserLevel(c)
	if level != authMdw.SUPERADMIN {
		c.JSON(response.ForbiddenLevel(level))
		return
	}
	domain := model.Domain{}
	if err := c.BindJSON(&domain); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	code, result := handler.domainService.PostDomain(c, domain)
	c.JSON(code, result)
}
