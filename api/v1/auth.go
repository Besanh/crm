package api

import (
	"contactcenter-api/api"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/response"
	_ "contactcenter-api/common/util"
	authMdw "contactcenter-api/middleware/auth"
	"contactcenter-api/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Auth struct {
	authService service.IAuth
}

func NewAuth(r *gin.Engine, authService service.IAuth) {
	handler := &Auth{
		authService: authService,
	}
	Group := r.Group(constants.VERSION + "/crm/auth")
	{
		Group.POST("", handler.Authen)
		Group.POST("token", handler.GenerateToken)
		Group.GET("auth-info", authMdw.AuthMiddleware(), handler.GetInfoByToken)
	}
}

func (handler *Auth) Authen(c *gin.Context) {
	var auth Login
	err := c.BindJSON(&auth)
	if err != nil {
		code, result := response.BadRequest()
		c.JSON(code, result)
		return
	}
	var username, domain string
	userDomain := strings.Split(auth.Username, "@")
	if len(userDomain) != 2 {
		code, result := response.BadRequestMsg("Username must be Username@Domain")
		c.JSON(code, result)
		return
	}
	username = userDomain[0]
	domain = userDomain[1]
	code, result := handler.authService.SigninContactCenter(c, domain, username, auth.Password)
	c.JSON(code, result)
}

func (handler *Auth) GenerateToken(c *gin.Context) {
	var body map[string]any
	err := c.BindJSON(&body)
	if err != nil {
		code, result := response.BadRequest()
		c.JSON(code, result)
		return
	}
	apiKey := body["api_key"].(string)
	if apiKey == "" {
		code, result := response.BadRequestMsg("api_key must not be null")
		c.JSON(code, result)
		return
	}
	code, result := handler.authService.GenerateTokenByApiKey(c, apiKey)
	c.JSON(code, result)
}

func (handler *Auth) GetInfoByToken(c *gin.Context) {
	_, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_uuid": userUuid,
	})
}
