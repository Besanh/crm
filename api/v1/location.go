package api

import (
	"contactcenter-api/common/constants"
	"contactcenter-api/common/util"
	authMdw "contactcenter-api/middleware/auth"
	"contactcenter-api/service"

	"github.com/gin-gonic/gin"
)

type Location struct {
	location service.ILocation
}

func NewLocation(r *gin.Engine, location service.ILocation) {
	handler := &Location{
		location: location,
	}

	Group := r.Group(constants.VERSION + "/crm/location")
	{
		Group.GET("province", authMdw.AuthMiddleware(), handler.GetLocationProvince)
		Group.GET("district/:province_code", authMdw.AuthMiddleware(), handler.GetLocationDistrict)
		Group.GET("ward", authMdw.AuthMiddleware(), handler.GetWard)
		Group.GET("ward/:district_code", authMdw.AuthMiddleware(), handler.GetLocationWard)

	}
}

func (h *Location) GetLocationProvince(c *gin.Context) {
	limit, offset := util.GetLimitOffset(c.Query("limit"), c.Query("offset"))
	code, result := h.location.GetLocationProvince(c, limit, offset)
	c.JSON(code, result)
}

func (h *Location) GetLocationDistrict(c *gin.Context) {
	provinceCode := c.Param("province_code")
	limit, offset := util.GetLimitOffset(c.Query("limit"), c.Query("offset"))
	code, result := h.location.GetLocationDistrict(c, provinceCode, limit, offset)
	c.JSON(code, result)
}

func (h *Location) GetLocationWard(c *gin.Context) {
	districtCode := c.Param("district_code")
	limit, offset := util.GetLimitOffset(c.Query("limit"), c.Query("offset"))
	code, result := h.location.GetLocationWard(c, districtCode, limit, offset)
	c.JSON(code, result)
}

func (l *Location) GetWard(c *gin.Context) {
	limit, offset := util.GetLimitOffset(c.Query("limit"), c.Query("offset"))
	wardName := c.Query("ward_name")
	code, result := l.location.GetLocationWardByName(c, wardName, limit, offset)
	c.JSON(code, result)
}
