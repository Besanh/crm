package api

import (
	"contactcenter-api/api"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/log"
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

type Profile struct {
	service service.IProfile
}

func NewProfile(r *gin.Engine, profileService service.IProfile) {
	profile := &Profile{
		service: profileService,
	}
	Group := r.Group(constants.VERSION + "/crm/profile")
	{
		Group.GET("", authMdw.AuthMiddleware(), profile.GetProfiles)
		Group.GET(":id", authMdw.AuthMiddleware(), profile.GetProfileById)
		Group.POST("", authMdw.AuthMiddleware(), profile.PostProfile)
		Group.PUT(":id", authMdw.AuthMiddleware(), profile.PutProfile)
		Group.DELETE(":id", authMdw.AuthMiddleware(), profile.DeleteProfile)
		Group.PATCH(":id/avatar", authMdw.AuthMiddleware(), profile.PatchContactAvatar)
		Group.GET("avatar_file/:file_name", profile.GetAvatar)
		Group.PATCH(":id", authMdw.AuthMiddleware(), profile.PatchStatusProfile)
		Group.PATCH(":id/social-mapping-contact", authMdw.AuthMiddleware(), profile.PatchSocialMappingContact)
		Group.GET("manage-profile", authMdw.AuthMiddleware(), profile.GetManageProfiles)
		Group.DELETE("delete-transaction", authMdw.AuthMiddleware(), profile.DeleteProfileTransaction)
		Group.POST("convert-lead-to-profile", authMdw.AuthMiddleware(), profile.PostConvertLeadToProfile)
	}
}

func (p *Profile) GetProfileById(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}
	code, result := p.service.GetProfileInfoById(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (p *Profile) GetProfiles(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	statusTmp := c.Query("status")
	status := sql.NullBool{}
	if len(statusTmp) > 0 {
		status.Valid = true
		status.Bool, _ = strconv.ParseBool(statusTmp)
	}
	filter := model.ProfileFilter{
		Fullname:       c.Query("fullname"),
		PhoneNumber:    c.Query("phone_number"),
		Email:          c.Query("email"),
		CustomerId:     c.Query("customer_id"),
		NationalId:     c.Query("national_id"),
		Gender:         c.Query("gender"),
		Address:        c.Query("address"),
		Country:        c.Query("country"),
		Province:       c.Query("province"),
		District:       c.Query("district"),
		Ward:           c.Query("ward"),
		Birthday:       c.Query("birthday"),
		JobTitle:       c.Query("job_title"),
		Status:         status,
		UserOwerUuid:   c.Query("user_owner_uuid"),
		RefId:          c.Query("ref_id"),
		RefCode:        c.Query("ref_code"),
		ProfileCode:    c.Query("profile_code"),
		IdentityNumber: c.Query("identity_number"),
		Passport:       c.Query("passport"),
		Common:         c.Query("common"),
		ProfileType:    c.Query("profile_type"),
		FacebookUserId: c.Query("facebook_user_id"),
		ZaloUserId:     c.Query("zalo_user_id"),
	}
	var errs error
	filter.StartTime, filter.EndTime, errs = util.ParseStartEndTime(c.Query("start_time"), c.Query("end_time"), true)
	if errs != nil {
		c.JSON(response.BadRequestMsg(errs.Error()))
		return
	}
	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))
	code, result := p.service.GetProfilesInfo(c, domainUuid, userUuid, filter, limit, offset)
	c.JSON(code, result)
}

func (p *Profile) PostProfile(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	body := make(map[string]any)
	if err := c.BindJSON(&body); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	profile := model.ProfilePost{}
	if err := util.ParseStruct(body, &profile); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	code, validSchema := validator.CheckSchema("profile/post.json", profile)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}
	code, result := p.service.PostProfile(c, domainUuid, userUuid, profile)
	c.JSON(code, result)
}

func (p *Profile) PutProfile(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	body := make(map[string]any)
	if err := c.BindJSON(&body); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}
	profile := model.ProfilePost{}
	if err := util.ParseStruct(body, &profile); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	code, validSchema := validator.CheckSchema("profile/post.json", profile)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}
	code, result := p.service.PutProfile(c, domainUuid, userUuid, id, profile)
	c.JSON(code, result)
}

func (p *Profile) DeleteProfile(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	id := c.Param("id")
	code, result := p.service.DeleteProfile(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (p *Profile) PatchContactAvatar(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	body := make(map[string]any)
	if err := c.BindJSON(&body); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	data, ok := body["data"].(string)
	if !ok || len(data) < 1 {
		c.JSON(response.BadRequestMsg("data is empty"))
		return
	}
	id := c.Param("id")
	code, result := p.service.PatchProfileAvatar(c, domainUuid, userUuid, id, data)
	c.JSON(code, result)
}

func (p *Profile) GetAvatar(c *gin.Context) {
	fileName := c.Param("file_name")
	if len(fileName) < 1 {
		c.JSON(response.BadRequestMsg("file name is empty"))
		return
	}

	var errs error
	var result string
	result, errs = p.service.GetAvatar(c, fileName)
	if errs == nil {
		fileByte, err := os.ReadFile(result)
		if err != nil {
			c.JSON(response.ServiceUnavailableMsg(err.Error()))
			return
		}
		contentType := http.DetectContentType(fileByte)
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%q", result))
		c.Writer.Header().Add("Content-Type", contentType)
		c.Data(http.StatusOK, contentType, fileByte)
		return
	} else {
		c.JSON(http.StatusNotFound, map[string]any{
			"content":   "not found",
			"file_name": fileName,
		})
	}
}

func (p *Profile) PatchStatusProfile(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	profileId := c.Param("id")
	if len(profileId) < 1 {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	jsonBody := make(map[string]any, 0)
	if err := c.BindJSON(&jsonBody); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	status, _ := jsonBody["status"].(bool)
	code, result := p.service.PatchStatusProfile(c, domainUuid, userUuid, profileId, status)
	c.JSON(code, result)
}

func (p *Profile) PatchSocialMappingContact(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	profileId := c.Param("id")
	if len(profileId) < 1 {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	jsonBody := make(map[string]any, 0)
	if err := c.BindJSON(&jsonBody); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	facebook, _ := jsonBody["facebook"].(string)
	zalo, _ := jsonBody["zalo"].(string)
	socialMappingContact := model.SocialMappingContact{
		Facebook: facebook,
		Zalo:     zalo,
	}
	code, result := p.service.PatchSocialMappingContactProfile(c, domainUuid, userUuid, profileId, socialMappingContact)
	c.JSON(code, result)
}

func (p *Profile) GetManageProfiles(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	statusTmp := c.Query("status")
	status := sql.NullBool{}
	if len(statusTmp) > 0 {
		status.Valid = true
		status.Bool, _ = strconv.ParseBool(statusTmp)
	}
	filter := model.ProfileFilter{
		Fullname:       c.Query("fullname"),
		PhoneNumber:    c.Query("phone_number"),
		Email:          c.Query("email"),
		CustomerId:     c.Query("customer_id"),
		NationalId:     c.Query("national_id"),
		Gender:         c.Query("gender"),
		Address:        c.Query("address"),
		Country:        c.Query("country"),
		Province:       c.Query("province"),
		District:       c.Query("district"),
		Ward:           c.Query("ward"),
		Birthday:       c.Query("birthday"),
		JobTitle:       c.Query("job_title"),
		Status:         status,
		UserOwerUuid:   c.Query("user_owner_uuid"),
		RefId:          c.Query("ref_id"),
		RefCode:        c.Query("ref_code"),
		ProfileCode:    c.Query("profile_code"),
		IdentityNumber: c.Query("identity_number"),
		Passport:       c.Query("passport"),
		Common:         c.Query("common"),
		ProfileType:    c.Query("profile_type"),
		FacebookUserId: c.Query("facebook_user_id"),
		ZaloUserId:     c.Query("zalo_user_id"),
		CreatedBy:      c.Query("created_by"),
	}
	var errs error
	filter.StartTime, filter.EndTime, errs = util.ParseStartEndTime(c.Query("start_time"), c.Query("end_time"), true)
	if errs != nil {
		c.JSON(response.BadRequestMsg(errs.Error()))
		return
	}

	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))

	code, result := p.service.GetManageProfiles(c, domainUuid, userUuid, filter, limit, offset)
	c.JSON(code, result)
}

func (p *Profile) DeleteProfileTransaction(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	// TODO: check level admin, manager

	filter := model.ProfileFilter{}
	if err := c.ShouldBind(&filter); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	code, result := p.service.DeleteProfileTransaction(c, domainUuid, userUuid, filter)
	c.JSON(code, result)
}

func (p *Profile) PostConvertLeadToProfile(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	body := make(map[string]any)
	if err := c.BindJSON(&body); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err))
		return
	}
	profile := model.ProfilePost{}
	if err := util.ParseStruct(body, &profile); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err))
		return
	}

	code, result := p.service.PostConvertLeadToProfile(c, domainUuid, userUuid, profile)
	c.JSON(code, result)
}
