package api

import (
	"contactcenter-api/common/cache"
	"contactcenter-api/common/constants"
	"contactcenter-api/common/model"
	"contactcenter-api/common/model/omni"
	authMdw "contactcenter-api/middleware/auth"
	"contactcenter-api/repository"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetInfoUser(c *gin.Context) (string, string, error) {
	domainUuid, ok := authMdw.GetUserDomainId(c)
	if !ok {
		err := errors.New("tenant is not exist")
		return "", "", err
	}
	userUuid, ok := authMdw.GetUserId(c)
	if !ok {
		err := errors.New("user is not exist")
		return "", "", err
	}
	return domainUuid, userUuid, nil
}

func ParseEventCalendar(c *gin.Context, userUuid, domainUuid string, eventCalendar model.EventCalendar) (model.EventCalendar, []model.EventCalendarTodo, error) {
	eventCalendarTodo := []model.EventCalendarTodo{}

	if title, ok := c.GetPostForm("title"); ok && len(title) > 0 {
		eventCalendar.Title = title
	}
	if description, ok := c.GetPostForm("description"); ok && len(description) > 0 {
		eventCalendar.Description = description
	}
	if categoryUuid, ok := c.GetPostForm("ecc_uuid"); ok && len(categoryUuid) > 0 {
		eventCalendar.EccUuid = categoryUuid
	} else {
		return eventCalendar, eventCalendarTodo, errors.New("event calendar category is not exist")
	}
	if remindTypeEvent, ok := c.GetPostForm("remind_type_event"); ok {
		eventCalendar.RemindTypeEvent = remindTypeEvent
	}
	if remindTime, ok := c.GetPostForm("remind_time"); ok {
		eventCalendar.RemindTime, _ = strconv.Atoi(remindTime)
	}
	if remindType, ok := c.GetPostForm("remind_type"); ok {
		eventCalendar.RemindType = remindType
	}
	if isWholeDay, ok := c.GetPostForm("is_whole_day"); ok {
		eventCalendar.IsWholeDay, _ = strconv.ParseBool(isWholeDay)
	}
	if isNotifyWeb, ok := c.GetPostForm("is_notify_web"); ok {
		eventCalendar.IsNotifyWeb, _ = strconv.ParseBool(isNotifyWeb)
	}
	if isNotifyEmail, ok := c.GetPostForm("is_notify_email"); ok {
		eventCalendar.IsNotifyEmail, _ = strconv.ParseBool(isNotifyEmail)
	}
	if isNotifySms, ok := c.GetPostForm("is_notify_sms"); ok {
		eventCalendar.IsNotifySms, _ = strconv.ParseBool(isNotifySms)
	}
	if isNotifyZns, ok := c.GetPostForm("is_notify_zns"); ok {
		eventCalendar.IsNotifyZns, _ = strconv.ParseBool(isNotifyZns)
	}
	if isNotifyCall, ok := c.GetPostForm("is_notify_call"); ok {
		eventCalendar.IsNotifyCall, _ = strconv.ParseBool(isNotifyCall)
	}
	if startTime, ok := c.GetPostForm("start_time"); ok {
		eventCalendar.StartTime, _ = time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	}
	if endTime, ok := c.GetPostForm("end_time"); ok {
		eventCalendar.EndTime, _ = time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
	}
	if repeat, ok := c.GetPostForm("repeat"); ok {
		eventCalendar.Repeat, _ = strconv.Atoi(repeat)
	}
	eventCalendar.Status = true
	eventCalendar.CreatedBy = userUuid
	eventCalendar.CreatedAt = time.Now()

	// todo
	if todos, ok := c.GetPostForm("todo"); ok {
		if err := json.Unmarshal([]byte(todos), &eventCalendarTodo); err != nil {
			return eventCalendar, eventCalendarTodo, err
		}
	}

	return eventCalendar, eventCalendarTodo, nil
}

func AuthOmni(c *gin.Context, id string) (omni.Omni, error) {
	omni := omni.Omni{}
	// Get cache
	_, err := cache.MCache.Get(constants.OMNI_INFO + "_" + id)
	if err != nil {
		return omni, err
	}
	domainUuid, _, err := GetInfoUser(c)
	if err != nil {
		return omni, err
	}

	omniExist, err := repository.OmniRepo.GetOmniById(c, domainUuid, id)
	if err != nil {
		return omni, err
	} else if len(omniExist.OmniUuid) < 1 {
		return omni, errors.New("omni is not exist")
	}

	// Set cache
	if err := cache.MCache.SetTTL(constants.OMNI_INFO+"_"+id, id, constants.TTL_OMNI); err != nil {
		return omni, err
	}
	omni = omniExist

	return omni, nil
}
