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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Contact struct {
	s service.IContact
}

func NewContact(r *gin.Engine, contactService service.IContact) {
	h := &Contact{
		s: contactService,
	}
	Group := r.Group(constants.VERSION + "/crm/contact")
	{
		Group.GET("", authMdw.AuthMiddleware(), h.GetContacts)
		Group.GET(":id", authMdw.AuthMiddleware(), h.GetContactById)
		Group.POST("", authMdw.AuthMiddleware(), h.PostContact)
		Group.PUT(":id", authMdw.AuthMiddleware(), h.PutContact)
		Group.DELETE(":id", authMdw.AuthMiddleware(), h.DeleteContact)
	}
}

func (h *Contact) GetContacts(c *gin.Context) {
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
	filter := model.ContactFilter{
		ContactType: c.Query("contact_type"),
		ContactName: c.Query("contact_name"),
	}
	var errs error
	filter.StartTime, filter.EndTime, errs = util.ParseStartEndTime(c.Query("start_time"), c.Query("end_time"), true)
	if errs != nil {
		c.JSON(response.BadRequestMsg(errs.Error()))
		return
	}
	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))
	code, result := h.s.GetContacts(c, domainUuid, userUuid, filter, limit, offset)
	c.JSON(code, result)
}

func (h *Contact) GetContactInfo(c *gin.Context) {
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
	filter := model.ContactFilter{
		ContactType: c.Query("contact_type"),
		ContactName: c.Query("contact_name"),
	}
	var errs error
	filter.StartTime, filter.EndTime, errs = util.ParseStartEndTime(c.Query("start_time"), c.Query("end_time"), true)
	if errs != nil {
		c.JSON(response.BadRequestMsg(errs.Error()))
		return
	}
	code, result := h.s.GetContactInfo(c, domainUuid, userUuid, filter)
	c.JSON(code, result)
}

func (h *Contact) PostContact(c *gin.Context) {
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
	contact := model.ContactPost{}
	if err := util.ParseStruct(body, &contact); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	code, validSchema := validator.CheckSchema("contact/post.json", contact)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}
	code, result := h.s.PostContact(c, domainUuid, userUuid, contact)
	c.JSON(code, result)
}

func (h *Contact) PutContact(c *gin.Context) {
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
	contact := model.ContactPost{}
	if err := util.ParseStruct(body, &contact); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	code, validSchema := validator.CheckSchema("contact/post.json", contact)
	if code != http.StatusOK {
		c.JSON(code, validSchema)
		return
	}
	code, result := h.s.PutContact(c, domainUuid, userUuid, id, contact)
	c.JSON(code, result)
}

func (h *Contact) DeleteContact(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	id := c.Param("id")
	code, result := h.s.DeleteContact(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (h *Contact) GetContactById(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}
	code, result := h.s.GetContactById(c, domainUuid, userUuid, id)
	c.JSON(code, result)
}

func (h *Contact) GetContactNotesOfId(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	id := c.Param("id")
	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))
	code, result := h.s.GetContactNotesOfId(c, domainUuid, userUuid, id, limit, offset)
	c.JSON(code, result)
}

func (h *Contact) PostContactNote(c *gin.Context) {
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
	contact := model.ContactNotePost{}
	if err := util.ParseStruct(body, &contact); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	code, result := h.s.PostContactNote(c, domainUuid, userUuid, contact)
	c.JSON(code, result)
}

func (h *Contact) PatchContactAvatar(c *gin.Context) {
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
	code, result := h.s.PatchContactAvatar(c, domainUuid, userUuid, id, data)
	c.JSON(code, result)
}

func (h *Contact) PostContactFile(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	isUpdateContact := false
	isSendMailStr, _ := c.GetPostForm("is_sendmail")
	if isSendMailStr == "true" {
		isUpdateContact = true
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	dir := util.PUBLIC_DIR + "contact/"
	dir += util.TimeToStringLayout(time.Now(), "2006_01_02")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
	filePath := dir + "/" + util.TimeToStringLayout(time.Now(), "2006_01_02_15_04_05") + "_" + file.Filename
	strArrs := strings.Split(file.Filename, ".") //  excel
	fileExtension := strArrs[len(strArrs)-1]
	if fileExtension != "xlsx" && fileExtension != "csv" {
		c.JSON(response.BadRequestMsg("File extension must be excel or csv"))
		return
	}
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	code, result := h.s.PostFileImportContacts(c, domainUuid, userUuid, filePath, fileExtension, isUpdateContact)
	c.JSON(code, result)
}

func (h *Contact) GetAvatar(c *gin.Context) {
	fileName := c.Param("file_name")
	if len(fileName) < 1 {
		c.JSON(response.BadRequestMsg("file name is empty"))
		return
	}

	var errs error
	var result string
	result, errs = h.s.GetAvatar(c, fileName)
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

func (h *Contact) PatchContactWithCallId(c *gin.Context) {
	domainUuid, userUuid, err := api.GetInfoUser(c)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	contactUuid := c.Param("contact_uuid")
	if len(contactUuid) < 1 {
		c.JSON(response.BadRequestMsg("contact_uuid is empty"))
		return
	}

	body := map[string]any{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	callId, ok := body["call_id"].(string)
	if !ok || len(callId) < 1 {
		c.JSON(response.BadRequestMsg("call_id is empty"))
		return
	}

	code, result := h.s.PatchContactWithCallId(c, domainUuid, userUuid, contactUuid, callId)
	c.JSON(code, result)
}
