package api

import (
	"contactcenter-api/api"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	authMdw "contactcenter-api/middleware/auth"
	"contactcenter-api/service"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Solution struct {
	solutionService service.ISolution
}

func NewSolution(r *gin.Engine, service service.ISolution) {
	handler := Solution{
		solutionService: service,
	}

	Group := r.Group(constants.VERSION + "/crm/solution")
	{
		Group.POST("", authMdw.AuthMiddleware(), handler.PostSolution)
		Group.GET("", authMdw.AuthMiddleware(), handler.GetSolutions)
		Group.GET(":id", authMdw.AuthMiddleware(), handler.GetSolutionById)
		Group.PUT(":id", authMdw.AuthMiddleware(), handler.PutSolutionById)
		Group.DELETE(":id", authMdw.AuthMiddleware(), handler.DeleteSolutionById)
		Group.POST("export", authMdw.AuthMiddleware(), handler.ExportSolutions)
	}
}
func (handler *Solution) PostSolution(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	solution := model.SolutionPost{}
	if err := c.Bind(&solution); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	code, result := handler.solutionService.PostSolution(c, domainUuid, userUuid, solution)
	c.JSON(code, result)
}

func (handler *Solution) GetSolutions(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	status := sql.NullBool{}
	if len(c.Query("status")) > 0 {
		status.Valid = true
		statusTmp, _ := strconv.ParseBool(c.Query("status"))
		status.Bool = statusTmp
	}
	filter := model.SolutionFilter{
		SolutionName: c.Query("solution_name"),
		SolutionCode: c.Query("solution_code"),
		Status:       status,
	}

	limit, offset := util.GetLimitOffset(c.Query("limit"), c.Query("offset"))

	code, result := handler.solutionService.GetSolutions(c, domainUuid, userUuid, filter, limit, offset)
	c.JSON(code, result)
}

func (handler *Solution) GetSolutionById(c *gin.Context) {
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

	code, result := handler.solutionService.GetSolutionById(c, domainUuid, id)
	c.JSON(code, result)
}

func (handler *Solution) PutSolutionById(c *gin.Context) {
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

	solution := model.Solution{}
	if err := c.Bind(&solution); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, result := handler.solutionService.PutSolutionById(c, domainUuid, userUuid, id, solution)
	c.JSON(code, result)
}

func (handler *Solution) DeleteSolutionById(c *gin.Context) {
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

	code, result := handler.solutionService.DeleteSolutionById(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (handler *Solution) ExportSolutions(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	var status sql.NullBool
	if len(c.Query("status")) > 0 {
		status.Valid = true
		statusTmp, _ := strconv.ParseBool(c.Query("status"))
		status.Bool = statusTmp
	}
	filter := model.SolutionFilter{
		SolutionName: c.Query("solution_name"),
		SolutionCode: c.Query("solution_code"),
		Status:       status,
		FileType:     c.Query("file_type"),
	}
	fileType := filter.FileType
	if len(fileType) < 1 {
		c.JSON(response.BadRequestMsg("file type is empty"))
		return
	}
	filePath, err := handler.solutionService.ExportSolutions(c, domainUuid, userUuid, fileType, filter)
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
