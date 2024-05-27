package api

import (
	"contactcenter-api/common/constants"
	"contactcenter-api/common/model"
	"contactcenter-api/common/util"
	authMdw "contactcenter-api/middleware/auth"
	"contactcenter-api/service"
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Career struct {
	careerService service.ICareer
}

func NewCareer(r *gin.Engine, careerService service.ICareer) {
	handler := Career{
		careerService: careerService,
	}

	Group := r.Group(constants.VERSION + "/crm/career")
	{
		Group.GET("", authMdw.AuthMiddleware(), handler.GetCareers)
	}
}

func (handler *Career) GetCareers(c *gin.Context) {
	isSearchExactlyTmp := c.Query("is_search_exactly")
	isSearchExactly := sql.NullBool{}
	if len(isSearchExactlyTmp) > 0 {
		isSearchExactly.Bool, _ = strconv.ParseBool(isSearchExactlyTmp)
		isSearchExactly.Valid = true
	}
	filter := model.CareerFilter{
		CareerCode:      c.Query("career_code"),
		CareerName:      c.Query("career_name"),
		Source:          util.ParseQueryArray(c.QueryArray("source")),
		IsSearchExactly: isSearchExactly,
	}
	code, result := handler.careerService.GetCareers(c, filter)
	c.JSON(code, result)
}
