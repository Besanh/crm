package util

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/araddon/dateparse"
	"github.com/golang-module/carbon"
	"github.com/google/uuid"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func ParseString(value any) string {
	str, ok := value.(string)
	if !ok {
		return str
	}
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Trim(str, "\r\n")
	str = strings.TrimSpace(str)
	return str
}

func ParseStringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

func ParseFloat64(value any) float64 {
	typeData := reflect.ValueOf(value)
	str := ""
	if typeData.Kind() == reflect.String {
		str = ParseString(value)
	} else {
		str = fmt.Sprintf("%v", value)
	}
	i, err := strconv.ParseFloat(str, 64)
	if err != nil {
		i = 0
	}
	return i
}

func ParseInt(str string) int {
	str = ParseString(str)
	i, err := strconv.Atoi(str)
	if err != nil {
		i = 0
	}
	return i
}

func ParseIsBool(value any) model.IsBool {
	val := ParseString(value)
	result := model.IsBool{
		IsNull: true,
		Value:  false,
	}
	if len(val) > 0 {
		result.IsNull = false
	}
	if val == "true" {
		result.Value = true
	}
	return result
}

func ParseOffset(offset string) int {
	offset = ParseString(offset)
	i, err := strconv.Atoi(offset)
	if err != nil {
		i = 0
	}
	return i
}

func ParseLimit(limit string) int {
	limit = ParseString(limit)
	i, err := strconv.Atoi(limit)
	if err != nil {
		i = 10
	}
	return i
}

func ParseTime(str string) time.Time {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Fatal(err)
		return time.Now()
	}
	time.Local = loc
	t, err := dateparse.ParseLocal(str)
	if err != nil {
		t = time.Now()
	}
	return t
}

func ParseMapToString(value any) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		return s
	}
	return output
}

func UrlEncode(s string) string {
	res := url.QueryEscape(s)
	return res
}

func UrlDecode(s string) string {
	res, err := url.QueryUnescape(s)
	if err != nil {
		return s
	}
	return res
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func GetPageSize(pageSize string) int64 {
	pageSizeInt, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		pageSizeInt = 50
	}
	return pageSizeInt
}

func CurrentTime() time.Time {
	return time.Now()
}

func CurrentTimeMicro() int64 {
	microTime := int64(time.Now().UnixNano() / 1000)
	return microTime
}

func TimeToString(valueTime time.Time) string {
	return TimeToStringLayout(valueTime, "2006-01-02 15:04:05")
}

func TimeToStringLayout(valueTime time.Time, layout string) string {
	return valueTime.Format(layout)
}

func ParseFromStringToTime(timeStr string) time.Time {
	return ParseFromStringToTimeLayout(timeStr, "2006-01-02 15:04:05")
}

func ParseFromStringToTimeLayout(timeStr string, layout string) time.Time {
	date, _ := time.Parse(layout, timeStr)
	return date
}

func CheckStartEndDate(startDate, endDate string) (time.Time, time.Time, error) {
	startTime := time.Now()
	endTime := time.Now()
	if startDate != "" {
		startTime = ParseFromStringToTime(startDate)
		if startTime.IsZero() {
			return time.Time{}, time.Time{}, errors.New("start_time is invalid")
		}
	} else {
		startDate = TimeToStringLayout(startTime, "2006-01-02") + " 00:00:00"
		startTime = ParseFromStringToTime(startDate)
	}
	if endDate != "" {
		endTime = ParseFromStringToTime(endDate)
		if endTime.IsZero() {
			return time.Time{}, time.Time{}, errors.New("end_time is invalid")
		}
	} else {
		endDate = TimeToStringLayout(startTime, "2006-01-02") + " 23:59:59"
		endTime = ParseFromStringToTime(endDate)
	}
	if startTime.After(endTime) {
		return time.Time{}, time.Time{}, errors.New("start_date must be after end_date")
	}
	return startTime, endTime, nil
}

func ParseQueryArray(slice []string) []string {
	result := make([]string, 0)
	for _, v := range slice {
		if len(v) > 0 {
			result = append(result, v)
		}
	}
	return result
}

func RemoveDuplicate(array []string) []string {
	m := make(map[string]string)
	for _, x := range array {
		m[x] = x
	}
	result := make([]string, 0)
	for x := range m {
		result = append(result, x)
	}
	return result
}

func InArray(item any, array any) bool {
	arr := reflect.ValueOf(array)
	if arr.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

func GetLocalTimeOfTime(val time.Time) time.Time {
	currentYear, currentMonth, currentDay := val.Date()
	loc, _ := time.LoadLocation("UTC")
	return time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, loc)
}

func ParseStartEndTime(startTimeStr, endTimeStr string, allowZero bool) (time.Time, time.Time, error) {
	today := time.Now()
	currentYear, currentMonth, currentDay := today.Date()
	loc, _ := time.LoadLocation("UTC")
	startTime := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, loc)
	endTime := time.Date(currentYear, currentMonth, currentDay, 23, 59, 59, 0, loc)
	if allowZero && len(startTimeStr) < 1 {
		startTime = time.Time{}
	} else if len(startTimeStr) > 1 {
		startTime = ParseFromStringToTime(startTimeStr)
		if startTime.IsZero() {
			return time.Time{}, time.Time{}, errors.New("start_time is invalid")
		}
	}
	if allowZero && len(endTimeStr) < 1 {
		endTime = time.Time{}
	} else if len(endTimeStr) > 1 {
		endTime = ParseFromStringToTime(endTimeStr)
		if endTime.IsZero() {
			return time.Time{}, time.Time{}, errors.New("end_time is invalid")
		}
	}
	if startTime.After(endTime) {
		return time.Time{}, time.Time{}, errors.New("start_date must be after end_date")
	}
	return startTime, endTime, nil
}

func GetStartEndCurrent() (time.Time, time.Time) {
	today := time.Now()
	currentYear, currentMonth, currentDay := today.Date()
	loc, _ := time.LoadLocation("UTC")
	startTime := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, loc)
	endTime := time.Date(currentYear, currentMonth, currentDay, 23, 59, 59, 0, loc)
	return startTime, endTime
}

func ParsesStringToStruct(value string, dest any) error {
	if err := json.Unmarshal([]byte(value), dest); err != nil {
		return err
	}
	return nil
}

func StringToBase64(str string) string {
	data := []byte(str)
	val := base64.StdEncoding.EncodeToString(data)
	return val
}

func ParseStruct(value any, dest any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, dest); err != nil {
		return err
	}
	return nil
}

func ToLower(value string) string {
	return strings.ToLower(value)
}

func InArr(array any, item any) bool {
	arr := reflect.ValueOf(array)
	if arr.Kind() != reflect.Slice {
		log.Error("invalid slice")
		return false
	}
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

func GetLimitOffset(limitQuery, offsetQuery string) (int, int) {
	limit := 0
	offset := 0
	limit = ParseLimit(limitQuery)
	offset = ParseOffset(offsetQuery)
	return limit, offset
}

func HandleResponseError(err error) (int, error) {
	j, err := json.Marshal(err)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	jsonMap := make(map[string]any)
	err = json.Unmarshal([]byte(j), &jsonMap)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	status := jsonMap["Status"].(int)
	return status, nil
}

func ParseToAnyArray(value []string) []any {
	result := make([]any, 0)
	for _, v := range value {
		result = append(result, v)
	}
	return result
}

func ParseAnyToAny(value any, dest any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, dest); err != nil {
		return err
	}
	return nil
}

func StripHtmlRegex(s string) string {
	const regex = `<.*?>`
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, "")
}

func ParseFromStringToTimeV2(timeStr string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func ParseFromToTime(startTimeStr, endTimeStr string, allowZero bool) (string, string, error) {
	today := time.Now()
	currentYear, currentMonth, currentDay := today.Date()
	loc, _ := time.LoadLocation("UTC")
	startTime := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, loc)
	endTime := time.Date(currentYear, currentMonth, currentDay, 23, 59, 59, 0, loc)

	if allowZero && len(startTimeStr) < 1 {
		startTime = time.Time{}
	} else if len(startTimeStr) > 1 {
		startTimeObj, err := ParseFromStringToTimeV2(startTimeStr)
		if err != nil {
			return "", "", err
		}
		startTime = startTimeObj
		if startTime.IsZero() {
			return "", "", fmt.Errorf("start_time is invalid")
		}
	}

	if allowZero && len(endTimeStr) < 1 {
		endTime = time.Time{}
	} else if len(endTimeStr) > 1 {
		endTimeObj, err := ParseFromStringToTimeV2(endTimeStr)
		if err != nil {
			return "", "", err
		}
		endTime = endTimeObj
		if endTime.IsZero() {
			return "", "", fmt.Errorf("end_time is invalid")
		}
	}

	if startTime.After(endTime) {
		return "", "", fmt.Errorf("start_date must be after end_date")
	}

	return startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"), nil
}

func TimestampToTime(timestamp int64) time.Time {
	return carbon.CreateFromTimestamp(timestamp).ToStdTime()
}
