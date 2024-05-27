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

type ClassifyTag struct {
	ClassifyTagService service.IClassifyTag
}

func NewClassifyTag(r *gin.Engine, classifyTagService service.IClassifyTag) {
	handler := &ClassifyTag{
		ClassifyTagService: classifyTagService,
	}

	Group := r.Group(constants.VERSION + "/crm/classify-tag")
	{
		Group.POST("", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), handler.PostClassifyTag)
		Group.GET("", authMdw.AuthMiddleware(), handler.GetClassifyTags)
		Group.GET(":id", authMdw.AuthMiddleware(), handler.GetClassifyTagById)
		Group.PUT(":id", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), handler.PutClassifyTagById)
		// Group.DELETE(":id", authMdw.AuthMiddleware(), authMdw.CheckLevelAdmin(), handler.DeleteClassifyTag)
		// Group.POST("export", authMdw.AuthMiddleware(), handler.ExportClassifyTags)
	}
}

func (handler *ClassifyTag) PostClassifyTag(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	ClassifyTag := model.ClassifyTag{}
	if err := c.ShouldBindJSON(&ClassifyTag); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, validSchema := validator.CheckSchema("classify_tag/post.json", ClassifyTag)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}

	code, result := handler.ClassifyTagService.PostClassifyTag(c, domainUuid, userUuid, &ClassifyTag)
	c.JSON(code, result)
}

func (handler *ClassifyTag) GetClassifyTags(c *gin.Context) {
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
	tagType := c.Query("tag_type")
	if len(tagType) > 0 && tagType == "all" {
		tagType = ""
	}
	filter := model.ClassifyTagFilter{
		TagName:   c.Query("tag_name"),
		TagType:   tagType,
		Status:    status,
		StartTime: c.Query("start_time"),
		EndTime:   c.Query("end_time"),
	}

	code, result := handler.ClassifyTagService.GetClassifyTags(c, domainUuid, limit, offset, filter)
	c.JSON(code, result)
}

func (handler *ClassifyTag) GetClassifyTagById(c *gin.Context) {
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

	code, result := handler.ClassifyTagService.GetClassifyTagById(c, domainUuid, id)
	c.JSON(code, result)
}

func (handler *ClassifyTag) PutClassifyTagById(c *gin.Context) {
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

	ClassifyTag := model.ClassifyTag{}
	if err := c.ShouldBindJSON(&ClassifyTag); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, validSchema := validator.CheckSchema("classify_tag/put.json", ClassifyTag)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}

	code, result := handler.ClassifyTagService.PutClassifyTagById(c, domainUuid, userUuid, id, ClassifyTag)
	c.JSON(code, result)
}

// func (handler *ClassifyTag) DeleteClassifyTag(c *gin.Context) {
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

// 	code, result := handler.ClassifyTagService.DeleteClassifyTag(c, domainUuid, id)
// 	c.JSON(code, result)
// }

// func (handler *ClassifyTag) ExportClassifyTags(c *gin.Context) {
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
// 	tagType := c.Query("tag_type")
// 	if len(tagType) > 0 && tagType == "all" {
// 		tagType = ""
// 	}
// 	filter := model.ClassifyTagFilter{
// 		TagName:   c.Query("tag_name"),
// 		TagType:   tagType,
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

// 	filePath, err := handler.ClassifyTagService.ExportClassifyTags(c, domainUuid, userUuid, fileType, filter)
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
