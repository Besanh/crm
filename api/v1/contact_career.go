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
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContactCareer struct {
	ContactCareerService service.IContactCareer
}

func NewContactCareer(r *gin.Engine, ContactCareerService service.IContactCareer) {
	handler := &ContactCareer{
		ContactCareerService: ContactCareerService,
	}

	Group := r.Group(constants.VERSION + "/crm/contact-career")
	{
		Group.POST("", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), handler.PostContactCareer)
		Group.GET("", authMdw.AuthMiddleware(), handler.GetContactCareers)
		Group.GET(":id", authMdw.AuthMiddleware(), handler.GetContactCareerById)
		Group.PUT(":id", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), handler.PutContactCareerById)
		Group.DELETE(":id", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), handler.DeleteContactCareer)
		Group.POST("/export", authMdw.AuthMiddleware(), handler.ExportContactCareers)
	}
}

func (handler *ContactCareer) PostContactCareer(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	contactCareer := model.ContactCareer{}
	if err := c.ShouldBindJSON(&contactCareer); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, validSchema := validator.CheckSchema("contact_career/post.json", contactCareer)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}

	code, result := handler.ContactCareerService.PostContactCareer(c, domainUuid, userUuid, &contactCareer)
	c.JSON(code, result)
}

func (handler *ContactCareer) GetContactCareers(c *gin.Context) {
	domainUuid, _, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	limit, offset := util.GetLimitOffset(c.Query("limit"), c.Query("offset"))

	statusTmp := c.Query("status")
	status := sql.NullBool{}
	if len(statusTmp) > 0 {
		status.Valid = true
		activeBool, _ := strconv.ParseBool(statusTmp)
		status.Bool = activeBool
	}
	filter := model.ContactCareerFilter{
		CareerName: c.Query("career_name"),
		CareerType: c.Query("career_type"),
		Status:     status,
		StartTime:  c.Query("start_time"),
		EndTime:    c.Query("end_time"),
	}

	code, result := handler.ContactCareerService.GetContactCareers(c, domainUuid, limit, offset, filter)
	c.JSON(code, result)
}

func (handler *ContactCareer) GetContactCareerById(c *gin.Context) {
	domainUuid, _, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	code, result := handler.ContactCareerService.GetContactCareerById(c, domainUuid, id)
	c.JSON(code, result)
}

func (handler *ContactCareer) PutContactCareerById(c *gin.Context) {
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

	contactCareer := model.ContactCareer{}
	if err := c.ShouldBindJSON(&contactCareer); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, validSchema := validator.CheckSchema("contact_career/put.json", contactCareer)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}

	code, result := handler.ContactCareerService.PutContactCareerById(c, domainUuid, userUuid, id, contactCareer)
	c.JSON(code, result)
}

func (handler *ContactCareer) DeleteContactCareer(c *gin.Context) {
	domainUuid, _, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	code, result := handler.ContactCareerService.DeleteContactCareer(c, domainUuid, id)
	c.JSON(code, result)
}

func (handler *ContactCareer) ExportContactCareers(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	statusTmp := c.Query("status")
	status := sql.NullBool{}
	if len(statusTmp) > 0 {
		status.Valid = true
		activeBool, _ := strconv.ParseBool(statusTmp)
		status.Bool = activeBool
	}
	filter := model.ContactCareerFilter{
		CareerName: c.Query("career_name"),
		CareerType: c.Query("career_type"),
		Status:     status,
		StartTime:  c.Query("start_time"),
		EndTime:    c.Query("end_time"),
		FileType:   c.Query("file_type"),
	}
	fileType := filter.FileType
	if len(fileType) < 1 {
		c.JSON(response.BadRequestMsg("file type is empty"))
		return
	}
	filePath, err := handler.ContactCareerService.ExportContactCareers(c, domainUuid, userUuid, fileType, filter)
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
