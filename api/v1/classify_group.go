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
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClassifyGroup struct {
	ClassifyGroupService service.IClassifyGroup
}

func NewClassifyGroup(r *gin.Engine, classifyGroupService service.IClassifyGroup) {
	handler := &ClassifyGroup{
		ClassifyGroupService: classifyGroupService,
	}

	Group := r.Group(constants.VERSION + "/crm/classify-group")
	{
		Group.POST("", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), handler.PostClassifyGroup)
		Group.GET("", authMdw.AuthMiddleware(), handler.GetClassifyGroups)
		Group.GET(":id", authMdw.AuthMiddleware(), handler.GetClassifyGroupById)
		Group.PUT(":id", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), handler.PutClassifyGroupById)
		// Group.DELETE(":id", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), handler.DeleteClassifyGroup)
		// Group.POST("export", authMdw.AuthMiddleware(), handler.ExportClassifyGroups)
	}
}

func (handler *ClassifyGroup) PostClassifyGroup(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	ClassifyGroup := model.ClassifyGroup{}
	if err := c.ShouldBindJSON(&ClassifyGroup); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, validSchema := validator.CheckSchema("classify_group/post.json", ClassifyGroup)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}

	code, result := handler.ClassifyGroupService.PostClassifyGroup(c, domainUuid, userUuid, &ClassifyGroup)
	c.JSON(code, result)
}

func (handler *ClassifyGroup) GetClassifyGroups(c *gin.Context) {
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
	GroupType := c.Query("Group_type")
	if len(GroupType) > 0 && GroupType == "all" {
		GroupType = ""
	}
	filter := model.ClassifyGroupFilter{
		GroupName: c.Query("Group_name"),
		GroupType: GroupType,
		Status:    status,
		StartTime: c.Query("start_time"),
		EndTime:   c.Query("end_time"),
	}

	code, result := handler.ClassifyGroupService.GetClassifyGroups(c, domainUuid, limit, offset, filter)
	c.JSON(code, result)
}

func (handler *ClassifyGroup) GetClassifyGroupById(c *gin.Context) {
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

	code, result := handler.ClassifyGroupService.GetClassifyGroupById(c, domainUuid, id)
	c.JSON(code, result)
}

func (handler *ClassifyGroup) PutClassifyGroupById(c *gin.Context) {
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

	ClassifyGroup := model.ClassifyGroup{}
	if err := c.ShouldBindJSON(&ClassifyGroup); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, validSchema := validator.CheckSchema("classify_group/put.json", ClassifyGroup)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}

	code, result := handler.ClassifyGroupService.PutClassifyGroupById(c, domainUuid, userUuid, id, ClassifyGroup)
	c.JSON(code, result)
}

// func (handler *ClassifyGroup) DeleteClassifyGroup(c *gin.Context) {
// 	domainUuid, _, err := api.GetInfoUser(c)
// 	if err != nil {
// 		c.JSON(response.BadRequestMsg(err))
// 		return
// 	}

// 	id := c.Param("id")
// 	if len(id) < 1 {
// 		c.JSON(response.BadRequestMsg("id is empty"))
// 		return
// 	}

// 	code, result := handler.ClassifyGroupService.DeleteClassifyGroup(c, domainUuid, id)
// 	c.JSON(code, result)
// }

// func (handler *ClassifyGroup) ExportClassifyGroups(c *gin.Context) {
// 	domainUuid, userUuid, err := api.GetInfoUser(c)
// 	if err != nil {
// 		c.JSON(response.BadRequestMsg(err))
// 		return
// 	}

// 	statusTmp := c.Query("status")
// 	status := sql.NullBool{}
// 	if len(statusTmp) > 0 {
// 		status.Valid = true
// 		activeBool, _ := strconv.ParseBool(statusTmp)
// 		status.Bool = activeBool
// 	}
// 	GroupType := c.Query("Group_type")
// 	if len(GroupType) > 0 && GroupType == "all" {
// 		GroupType = ""
// 	}
// 	filter := model.ClassifyGroupFilter{
// 		GroupName:   c.Query("Group_name"),
// 		GroupType:   GroupType,
// 		Status:    status,
// 		StartTime: c.Query("start_time"),
// 		EndTime:   c.Query("end_time"),
// 		FileType:  c.Query("file_type"),
// 	}
// 	fileType := filter.FileType
// 	if len(fileType) < 1 {
// 		c.JSON(response.BadRequestMsg("file type is empty"))
// 		return
// 	}

// 	filePath, err := handler.ClassifyGroupService.ExportClassifyGroups(c, domainUuid, userUuid, fileType, filter)
// 	if err != nil {
// 		c.JSON(response.ServiceUnavailableMsg(err.Error()))
// 	} else {
// 		fileByte, err := os.ReadFile(filePath)
// 		if err != nil {
// 			code, result := response.ServiceUnavailableMsg(err.Error())
// 			c.JSON(code, result)
// 			return
// 		}
// 		contentType := http.DetectContentType(fileByte)
// 		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filePath))
// 		c.Writer.Header().Add("Content-Type", contentType)
// 		c.Data(http.StatusOK, contentType, fileByte)
// 		return
// 	}
// }
